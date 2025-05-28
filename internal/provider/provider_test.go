package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"selectel-baremetal": providerserver.NewProtocol6WithError(New("test")()),
}

func TestProvider(t *testing.T) {
	provider := New("test")()

	if provider == nil {
		t.Fatal("Expected provider to be non-nil")
	}

	// Test that provider implements the expected interface
	if _, ok := provider.(*SelectelBaremetalProvider); !ok {
		t.Fatal("Expected provider to be of type *SelectelBaremetalProvider")
	}

	// Test that provider factories are configured
	if len(testAccProtoV6ProviderFactories) == 0 {
		t.Fatal("Expected testAccProtoV6ProviderFactories to be configured")
	}
}
