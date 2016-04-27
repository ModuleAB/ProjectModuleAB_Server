package common

import (
	"fmt"
	"moduleab_server/oas"

	"github.com/astaxie/beego"
)

var DefaultOasClient *oas.OasClient

func InitOasClient(endpoint string) error {
	oasPort, err := beego.AppConfig.Int("aliapi::oasport")
	if err != nil {
		return fmt.Errorf("Bad config value type (expect int): apiapi::oasport")
	}
	oasUseSSL, err := beego.AppConfig.Bool("aliapi::oasusessl")
	if err != nil {
		return fmt.Errorf("Bad config value type (expect bool): apiapi::oasport")
	}

	DefaultClient = oas.NewOasClient(
		a.Endpoint,
		beego.AppConfig.String("aliapi::apikey"),
		beego.AppConfig.String("aliapi::secret"),
		oasPort,
		oasUseSSL,
	)
	return nil
}

func GetOasVaultId(name string) (string, error) {
	v := new(oas.VaultsList)
	for {
		id, v, err := DefaultClient.ListVaults(-1, v.Marker)
		beego.Debug("OAS request ID:", id)
		if err != nil {

			return "", err
		}
		for _, xv := range v.VaultList {
			if xv.VaultName == name {
				return xv.VaultID, nil
			}
		}
		if v.Marker == "" {
			return "", fmt.Errorf("Vault not found")
		}
	}
}
