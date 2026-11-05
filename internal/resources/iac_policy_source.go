// Copyright Stacklet, Inc. 2025, 2026

package resources

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/stacklet/terraform-provider-stacklet/internal/api/iac"
	"github.com/stacklet/terraform-provider-stacklet/internal/errors"
	"github.com/stacklet/terraform-provider-stacklet/internal/models"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &iacPolicySourceResource{}
	_ resource.ResourceWithConfigure   = &iacPolicySourceResource{}
	_ resource.ResourceWithImportState = &iacPolicySourceResource{}
)

// iacPolicySourceResource defines the resource implementation.
type iacPolicySourceResource struct {
	iacAPIResource
}

func (r *iacPolicySourceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iac_policy_source"
}

func (r *iacPolicySourceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Stacklet IaC policy source, a repository containing IaC policies. Deleting a policy source also deletes its policies, removing them from all policy collections.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The unique name of the policy source.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"url": schema.StringAttribute{
				Description: "The URL of the remote repository.",
				Required:    true,
			},
			"branch": schema.StringAttribute{
				Description: "The branch to scan for policies. If not set, the default branch is used.",
				Optional:    true,
			},
			"paths": schema.SetAttribute{
				Description: "The directories to scan for policies. If not set, the whole repository is scanned.",
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
			"extensions": schema.SetAttribute{
				Description: "The extensions of policy files, including the leading dot.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Default: setdefault.StaticValue(
					types.SetValueMust(types.StringType, []attr.Value{
						types.StringValue(".yml"),
						types.StringValue(".yaml"),
					}),
				),
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(regexp.MustCompile(`^\.`), "must start with a dot"),
					),
				},
			},
			"auto_scan": schema.BoolAttribute{
				Description: "Whether the repository is periodically scanned for policy changes.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"webhook_url": schema.StringAttribute{
				Description: "A URL which triggers scans of the remote repository.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"token_auth": schema.SingleNestedAttribute{
				Description: "Authentication to the remote repository with a user and password or token.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"user": schema.StringAttribute{
						Description: "The user with access to the remote repository.",
						Required:    true,
					},
					"token_wo": schema.StringAttribute{
						Description: "The password or token to access the remote repository.",
						Required:    true,
						Sensitive:   true,
						WriteOnly:   true,
					},
					"token_wo_version": schema.StringAttribute{
						Description: "Change value to update token_wo.",
						Required:    true,
					},
				},
				Validators: []validator.Object{
					objectvalidator.ConflictsWith(path.MatchRoot("ssh_auth")),
				},
			},
			"ssh_auth": schema.SingleNestedAttribute{
				Description: "Authentication to the remote repository with an SSH key.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"user": schema.StringAttribute{
						Description: "The user with access to the remote repository.",
						Required:    true,
					},
					"private_key_wo": schema.StringAttribute{
						Description: "The SSH private key to access the remote repository.",
						Required:    true,
						Sensitive:   true,
						WriteOnly:   true,
					},
					"private_key_wo_version": schema.StringAttribute{
						Description: "Change value to update private_key_wo.",
						Required:    true,
					},
					"passphrase_wo": schema.StringAttribute{
						Description: "The passphrase for the SSH private key.",
						Optional:    true,
						Sensitive:   true,
						WriteOnly:   true,
						Validators: []validator.String{
							stringvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("passphrase_wo_version")),
						},
					},
					"passphrase_wo_version": schema.StringAttribute{
						Description: "Change value to update passphrase_wo.",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.AlsoRequires(path.MatchRelative().AtParent().AtName("passphrase_wo")),
						},
					},
				},
			},
		},
	}
}

func (r *iacPolicySourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, config models.IaCPolicySourceResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	paths, extensions, diags := r.pathsAndExtensions(ctx, plan)
	resp.Diagnostics.Append(diags...)
	auth, diags := r.auth(ctx, config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := iac.PolicySourceCreateInput{
		Name:       plan.Name.ValueString(),
		URL:        plan.URL.ValueString(),
		Branch:     plan.Branch.ValueStringPointer(),
		Paths:      paths,
		Extensions: extensions,
		AutoScan:   plan.AutoScan.ValueBool(),
		Auth:       auth,
	}
	source, err := r.api.PolicySource.Create(ctx, input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(plan.Update(ctx, source)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *iacPolicySourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.IaCPolicySourceResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	source, err := r.api.PolicySource.Read(ctx, state.Name.ValueString())
	if err != nil {
		handleAPIError(ctx, &resp.State, &resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(state.Update(ctx, source)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *iacPolicySourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state, config models.IaCPolicySourceResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	paths, extensions, diags := r.pathsAndExtensions(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := iac.PolicySourceUpdateInput{
		URL:        plan.URL.ValueString(),
		Branch:     plan.Branch.ValueStringPointer(),
		Paths:      paths,
		Extensions: extensions,
		AutoScan:   plan.AutoScan.ValueBool(),
	}
	// The API replaces authentication as a whole, so if any part of it
	// changes, it's sent again in full from the config. Write-only values
	// are always null in plan and state, so only the user and versions are
	// compared.
	if !plan.TokenAuth.Equal(state.TokenAuth) || !plan.SSHAuth.Equal(state.SSHAuth) {
		auth, diags := r.auth(ctx, config)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		input.Auth = iac.PolicySourceAuthUpdate{Set: true, Auth: auth}
	}
	source, err := r.api.PolicySource.Update(ctx, plan.Name.ValueString(), input)
	if err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(plan.Update(ctx, source)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *iacPolicySourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.IaCPolicySourceResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.api.PolicySource.Delete(ctx, state.Name.ValueString()); err != nil {
		errors.AddDiagError(&resp.Diagnostics, err)
		return
	}
}

func (r *iacPolicySourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importState(ctx, req, resp, []string{"name"})
}

// auth returns the authentication for the policy source from the config, or
// nil if not configured.
func (r *iacPolicySourceResource) auth(ctx context.Context, config models.IaCPolicySourceResource) (*iac.PolicySourceAuth, diag.Diagnostics) {
	var diags diag.Diagnostics

	switch {
	case !config.TokenAuth.IsNull():
		var tokenAuth models.IaCPolicySourceTokenAuth
		diags.Append(config.TokenAuth.As(ctx, &tokenAuth, basetypes.ObjectAsOptions{})...)
		return &iac.PolicySourceAuth{
			Username: tokenAuth.User.ValueString(),
			Token:    tokenAuth.TokenWO.ValueStringPointer(),
		}, diags
	case !config.SSHAuth.IsNull():
		var sshAuth models.IaCPolicySourceSSHAuth
		diags.Append(config.SSHAuth.As(ctx, &sshAuth, basetypes.ObjectAsOptions{})...)
		return &iac.PolicySourceAuth{
			Username:      sshAuth.User.ValueString(),
			SSHKey:        sshAuth.PrivateKeyWO.ValueStringPointer(),
			SSHPassphrase: sshAuth.PassphraseWO.ValueStringPointer(),
		}, diags
	default:
		return nil, diags
	}
}

func (r *iacPolicySourceResource) pathsAndExtensions(ctx context.Context, plan models.IaCPolicySourceResource) ([]string, []string, diag.Diagnostics) {
	var diags diag.Diagnostics
	var paths, extensions []string
	if !plan.Paths.IsNull() {
		diags.Append(plan.Paths.ElementsAs(ctx, &paths, false)...)
	}
	diags.Append(plan.Extensions.ElementsAs(ctx, &extensions, false)...)
	return paths, extensions, diags
}
