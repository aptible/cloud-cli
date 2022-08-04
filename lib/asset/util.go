package libasset

import (
	"strings"

	cac "github.com/aptible/cloud-api-clients/clients/go"
)

func GetName(asset cac.AssetOutput) string {
	assetName := asset.CurrentAssetParameters.Data["name"]
	if assetName == "" {
		assetName = "N/A"
	}

	return assetName.(string)
}

func FilterByType(assets []cac.AssetOutput, types []string) []cac.AssetOutput {
	filteredResults := make([]cac.AssetOutput, 0)
	for _, result := range assets {
		for _, _type := range types {
			if strings.Contains(result.Asset, _type) {
				filteredResults = append(filteredResults, result)
			}
		}
	}
	return filteredResults
}
