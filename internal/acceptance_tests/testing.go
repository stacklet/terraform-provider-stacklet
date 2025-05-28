// Copyright (c) 2025 - Stacklet, Inc.

package acceptance_tests

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"text/template"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/stacklet/terraform-provider-stacklet/internal/provider"
)

// importStateIDFuncFromAttrs returns an ImportStateIdFunc that creates an
// import ID from resource attributes. Each attribute is in the form
// `resource.name.attr`. If multiple attributes are provided, they're joined
// with ":".
func importStateIDFuncFromAttrs(attrs ...string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		values := make([]string, len(attrs))
		for i, attr := range attrs {
			parts := strings.Split(attr, ".")
			resource := strings.Join(parts[:len(parts)-1], ".")
			name := parts[len(parts)-1]

			res, ok := s.RootModule().Resources[resource]
			if !ok {
				return "", fmt.Errorf("resource '%s' not found in state", resource)
			}

			value, ok := res.Primary.Attributes[name]
			if !ok {
				return "", fmt.Errorf("resource '%s' doesn't have attribute '%s'", resource, attr)
			}
			values[i] = value
		}
		return strings.Join(values, ":"), nil
	}
}

// runRecordedAccTest runs an acceptance test, with the specified name and steps.
func runRecordedAccTest(t *testing.T, testName string, testSteps []resource.TestStep) {
	setupHTTPTransport(t, testName)
	renderConfigs(t, testSteps)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"stacklet": func() (tfprotov6.ProviderServer, error) {
				p := provider.New("test")()
				return providerserver.NewProtocol6WithError(p)()
			},
		},
		Steps: testSteps,
	})
}

func setupHTTPTransport(t *testing.T, testName string) {
	rt := newRecordedTransport(t, testName, http.DefaultTransport)
	if err := rt.loadRecording(); err != nil {
		t.Errorf("failed to load recording: %v", err)
	}

	origTransport := http.DefaultTransport
	http.DefaultTransport = rt

	t.Cleanup(func() {
		http.DefaultTransport = origTransport

		if err := rt.saveRecording(); err != nil {
			t.Errorf("failed to save recording: %v", err)
		}
	})
}

// configData holds the test configuration data.
type configData struct {
	Prefix string
}

func getConfigData() configData {
	return configData{
		Prefix: getenv("TF_ACC_PREFIX", "test"),
	}
}

func getenv(name, fallback string) string {
	value := os.Getenv(name)
	if value == "" {
		return fallback
	}
	return value
}

func renderConfigs(t *testing.T, testSteps []resource.TestStep) {
	data := getConfigData()
	for i, step := range testSteps {
		if step.Config != "" {
			config, err := renderConfig(step.Config, data)
			if err != nil {
				t.Errorf("failed to render config: %v", err)
			}
			step.Config = config
			testSteps[i] = step
		}
	}
}

func renderConfig(config string, data configData) (string, error) {
	t, err := template.New("config").Parse(config)
	if err != nil {
		return "", err
	}

	wr := bytes.Buffer{}
	if err := t.Execute(&wr, data); err != nil {
		return "", err
	}
	return wr.String(), nil
}

func prefixName(name string) string {
	return fmt.Sprintf("%s-%s", getConfigData().Prefix, name)
}
