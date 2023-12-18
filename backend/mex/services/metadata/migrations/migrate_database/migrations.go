//revive:disable:var-naming
package migrate_database

import (
	"fmt"
	"regexp"
)

type BindataMigrations struct{}

const assetNameInitScript = "init.sql"

func (BindataMigrations) InitScript() string {
	return string(MustAsset(assetNameInitScript))
}

func (BindataMigrations) ScriptForVersion(version int) string {
	pattern := regexp.MustCompile(fmt.Sprintf("^%02d_[a-zA-Z_0-9]+\\.sql$", version))

	for _, assetName := range AssetNames() {
		if pattern.MatchString(assetName) {
			return string(MustAsset(assetName))
		}
	}
	return ""
}
