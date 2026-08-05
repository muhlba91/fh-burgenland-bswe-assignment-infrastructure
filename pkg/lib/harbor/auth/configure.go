package auth

import (
	"os"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumiverse/pulumi-harbor/sdk/v3/go/harbor"
)

// Configure configures Harbor authentication with OIDC using Dex.
// ctx: pulumi.Context.
func Configure(ctx *pulumi.Context) error {
	_, err := harbor.NewConfigAuth(ctx, "harbor-auth-dex", &harbor.ConfigAuthArgs{
		AuthMode:         pulumi.String("oidc_auth"),
		PrimaryAuthMode:  pulumi.Bool(true),
		OidcAutoOnboard:  pulumi.Bool(true),
		OidcClientId:     pulumi.String("registry"),
		OidcClientSecret: pulumi.String(os.Getenv("DEX_REGISTRY_CLIENT_SECRET")),
		OidcEndpoint:     pulumi.String("https://auth.hochschule-burgenland.muehlbachler.xyz/api/dex"),
		OidcName:         pulumi.String("GITHUB"),
		OidcUserClaim:    pulumi.String("preferred_username"),
		OidcGroupsClaim:  pulumi.String("groups"),
		OidcAdminGroup:   pulumi.String("fhburgenland-bswe:dmuehlbachler"),
		OidcScope:        pulumi.String("openid,email,profile,groups,offline_access"),
	})
	return err
}
