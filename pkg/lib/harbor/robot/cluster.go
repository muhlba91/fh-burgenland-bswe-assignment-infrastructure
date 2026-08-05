package robot

import (
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/muhlba91/fh-burgenland-bswe-assignment-infrastructure/pkg/lib/config"
	"github.com/muhlba91/pulumi-shared-library/pkg/lib/vault/secret"
	"github.com/muhlba91/pulumi-shared-library/pkg/util/encoding"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumiverse/pulumi-harbor/sdk/v3/go/harbor"
)

// CreateCluster configures the cluster-level Harbor robot account.
// ctx: The Pulumi context.
func CreateCluster(
	ctx *pulumi.Context,
) error {
	if config.Classroom.Vault == nil || *config.Classroom.Vault == "" {
		log.Info().Msg("[harbor] vault is disabled, skipping harbor cluster robot account creation")
		return nil
	}

	ra, _ := harbor.NewRobotAccount(ctx, fmt.Sprintf("harbor-robot-%s", config.Environment), &harbor.RobotAccountArgs{
		Level:       pulumi.String("system"),
		Description: pulumi.String(fmt.Sprintf("Robot account for %s", config.Environment)),
		Permissions: &harbor.RobotAccountPermissionArray{
			&harbor.RobotAccountPermissionArgs{
				Kind:      pulumi.String("project"),
				Namespace: pulumi.String("*"),
				Accesses: &harbor.RobotAccountPermissionAccessArray{
					&harbor.RobotAccountPermissionAccessArgs{
						Action:   pulumi.String("pull"),
						Resource: pulumi.String("repository"),
						Effect:   pulumi.String("allow"),
					},
				},
			},
		},
	})

	vaultValue, _ := (pulumi.All(ra.FullName, ra.Secret).ApplyT(func(args []any) string {
		name, ok := args[0].(string)
		if !ok {
			log.Error().Msgf("[harbor] failed to cast robot name for %s", name)
		}
		secretKey, ok := args[1].(string)
		if !ok {
			log.Error().Msgf("[harbor] failed to cast robot secret for %s", name)
		}
		dockerconfigjson, errDCMarshal := json.Marshal(map[string]map[string]map[string]string{
			"auths": {
				harborURL(): {
					"username": name,
					"password": secretKey,
					"auth":     encoding.B64Encode(fmt.Sprintf("%s:%s", name, secretKey)),
				},
			},
		})
		if errDCMarshal != nil {
			log.Error().Err(errDCMarshal).Msgf("[harbor] failed to marshal dockerconfigjson for %s", name)
		}
		data, errMarshal := json.Marshal(map[string]string{
			"name":             name,
			"secret_key":       secretKey,
			"dockerconfigjson": encoding.B64Encode(string(dockerconfigjson)),
		})
		if errMarshal != nil {
			log.Error().Err(errMarshal).Msgf("[harbor] failed to marshal credentials for %s", name)
		}
		return string(data)
	})).(pulumi.StringOutput)

	_, errVault := secret.Create(ctx, &secret.CreateOptions{
		Key:   "registry-credentials",
		Value: vaultValue,
		Path:  *config.Classroom.Vault,
	})

	return errVault
}
