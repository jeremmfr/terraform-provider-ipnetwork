package provider_test

import (
	"github.com/jeremmfr/terraform-provider-ipnetwork/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){ //nolint:gochecknoglobals
	"ipnetwork": providerserver.NewProtocol6WithError(provider.New()),
}
