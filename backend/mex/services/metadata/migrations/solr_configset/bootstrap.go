//revive:disable:var-naming
package solr_configset

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	L "github.com/d4l-data4life/mex/mex/shared/log"
	"github.com/d4l-data4life/mex/mex/shared/solr"
	"github.com/d4l-data4life/mex/mex/shared/utils"
)

type collection struct {
	ConfigName string `json:"configName"`
}

type cluster struct {
	Collections map[string]collection `json:"collections"`
}

type bodyGetCollectionsResponse struct {
	Cluster cluster `json:"cluster"`
}

type bodyGetConfigsetsResponse struct {
	Configsets []string `json:"configSets"`
}

func Bootstrap(ctx context.Context, log L.Logger, solrClient solr.ClientAPI, configsetName string) error {
	var err error
	var body []byte

	zipData, err := zipConfigset(ctx, log, configsetName)
	if err != nil {
		return err
	}

	// We only can reliably update/overwrite the existing configset in Solr if it is not used by an existing collection/core.
	// Just trying to delete a configset is not possible because Solr does not allow for a proper programmatic distinguitation
	// of the conditions 'configset not found' and 'configset used by a collection'. (In the first case we would upload the
	// configset, in the second case this would likely (*) lead to an error.)
	//
	// This is why we check the existing configsets and their possible collections and only upload the configset if it is safe.
	//
	// (*) Not all overwriting leads to an error. Only if there are fields created based on field types from the configset.

	statusCode, body, err := solrClient.DoRequest(ctx, "GET", "/api/cluster/configs?omitHeader=true", nil)
	if err != nil {
		return err
	}
	if statusCode != http.StatusOK {
		return fmt.Errorf("error getting configsets: %d: %s", statusCode, string(body))
	}
	var bodyGetConfigsetsResponse bodyGetConfigsetsResponse
	if err := json.Unmarshal(body, &bodyGetConfigsetsResponse); err != nil {
		return err
	}

	if utils.Contains(bodyGetConfigsetsResponse.Configsets, configsetName) {
		if _, body, err = solrClient.DoRequest(ctx, "GET", "/solr/admin/collections?action=CLUSTERSTATUS&wt=json", nil); err != nil {
			return err
		}

		var bodyGetCollectionsResponse bodyGetCollectionsResponse
		if err := json.Unmarshal(body, &bodyGetCollectionsResponse); err != nil {
			return err
		}

		collections := utils.KeysOfMap(bodyGetCollectionsResponse.Cluster.Collections)
		log.Info(ctx, L.Messagef("Solr configsets: %v", collections))

		for _, collection := range collections {
			if c, ok := bodyGetCollectionsResponse.Cluster.Collections[collection]; ok {
				if c.ConfigName == configsetName {
					log.Info(ctx, L.Messagef("configset '%s' used by collection '%s'; not attempting deletion/update", configsetName, collection))
					return nil
				}
			}
		}

		if statusCode, _, err = solrClient.DoRequest(ctx, "DELETE", fmt.Sprintf("/api/cluster/configs/%s?omitHeader=true", configsetName), nil); err != nil {
			return err
		}
		if statusCode != http.StatusOK {
			return fmt.Errorf("cannot delete/overwrite Solr configset: %d: %s", statusCode, configsetName)
		}
	}

	if statusCode, body, err = solrClient.DoRequest(ctx, "PUT", fmt.Sprintf("/api/cluster/configs/%s", configsetName), zipData); err != nil {
		return err
	}
	if statusCode != http.StatusOK {
		return fmt.Errorf("configset upload to Solr failed: %d: %s", statusCode, string(body))
	}

	log.Info(ctx, L.Messagef("configset uploaded: %s", configsetName))
	return nil
}

func zipConfigset(ctx context.Context, log L.Logger, configsetName string) ([]byte, error) {
	var b bytes.Buffer
	z := zip.NewWriter(&b)

	err := zipAssets(ctx, log, z, configsetName+"/", configsetName)
	if err != nil {
		return nil, err
	}

	err = z.Close()
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func zipAssets(ctx context.Context, log L.Logger, z *zip.Writer, prefix string, path string) error {
	dirContents, err := AssetDir(path)
	if err != nil {
		return err
	}

	for _, asset := range dirContents {
		// This call is needed to make sure there is a zip entry for the folder.
		// Without it the Solr zip reader is unable to read/detect subfolders in the zip.
		_, _ = z.Create(strings.TrimPrefix(path, prefix))

		data, err := Asset(path + "/" + asset)
		if err == nil {
			// If no error: asset is a file: zip it
			log.Info(ctx, L.Messagef("add file to configset zip file: %s", path+"/"+asset))
			w, err := z.Create(strings.TrimPrefix(path+"/"+asset, prefix))
			if err != nil {
				return err
			}

			_, err = w.Write(data)
			if err != nil {
				return err
			}

			continue
		}

		// Arriving here means the asset name refers to a folder: recurse into it
		log.Info(ctx, L.Messagef("recursing into configset folder: %s", path+"/"+asset))
		err = zipAssets(ctx, log, z, prefix, path+"/"+asset)
		if err != nil {
			return err
		}
	}

	return nil
}
