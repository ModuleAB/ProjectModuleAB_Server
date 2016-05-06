package common

import (
	"fmt"
	"moduleab_server/oas"

	"github.com/astaxie/beego"
)

type OasClient struct {
	*oas.OasClient
}

func (o *OasClient) GetOasVaultId(name string) (string, error) {
	v := new(oas.VaultsList)
	for {
		id, v, err := o.ListVaults(-1, v.Marker)
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

var DefaultOasClient *OasClient

func NewOasClient(endpoint string) (*OasClient, error) {
	oasPort, err := beego.AppConfig.Int("aliapi::oasport")
	if err != nil {
		return nil, fmt.Errorf("Bad config value type (expect int): apiapi::oasport")
	}
	oasUseSSL, err := beego.AppConfig.Bool("aliapi::oasusessl")
	if err != nil {
		return nil, fmt.Errorf("Bad config value type (expect bool): apiapi::oasport")
	}
	o := new(OasClient)
	o.OasClient = oas.NewOasClient(
		endpoint,
		beego.AppConfig.String("aliapi::apikey"),
		beego.AppConfig.String("aliapi::secret"),
		oasPort,
		oasUseSSL,
	)
	return o, nil
}
