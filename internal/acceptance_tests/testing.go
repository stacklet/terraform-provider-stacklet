// Copyright (c) 2025 - Stacklet, Inc.

package acceptance_tests

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/stacklet/terraform-provider-stacklet/internal/provider"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a new provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"stacklet": func() (tfprotov6.ProviderServer, error) {
		p := provider.New("test")()
		return providerserver.NewProtocol6WithError(p)()
	},
}

func runRecordedAccTest(t *testing.T, testName string, testSteps []resource.TestStep) {
	rt := newRecordedTransport(t, testName, http.DefaultTransport)
	if err := rt.loadRecordings(); err != nil {
		t.Errorf("failed to load recordings: %v", err)
	}

	origTransport := http.DefaultTransport
	http.DefaultTransport = rt

	t.Cleanup(func() {
		http.DefaultTransport = origTransport

		if err := rt.saveRecordings(); err != nil {
			t.Errorf("failed to save recordings: %v", err)
		}
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps:                    testSteps,
	})
}
