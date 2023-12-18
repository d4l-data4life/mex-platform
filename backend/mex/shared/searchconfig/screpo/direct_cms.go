package screpo

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/d4l-data4life/mex/mex/shared/searchconfig"
	"github.com/d4l-data4life/mex/mex/shared/solr"
)

const apiPath = "/api/v0/config/files"

type searchConfigRepoDirectCMS struct {
	originCMS           string
	strictConfigParsing bool
}

func NewDirectCMSSearchConfigRepo(originCMS string, strictConfigParsing bool) searchconfig.SearchConfigRepo {
	return &searchConfigRepoDirectCMS{
		originCMS:           originCMS,
		strictConfigParsing: strictConfigParsing,
	}
}

// ListSearchConfigs lists all config objects
func (repo *searchConfigRepoDirectCMS) ListSearchConfigs(_ context.Context) (*searchconfig.SearchConfigList, error) {
	client := &http.Client{}
	resp, err := client.Get(fmt.Sprintf("%s%s/search_configs", repo.originCMS, apiPath))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("could not fetch search configurations (of specific type) from CMS - got response status code %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	protoSearchConfigList := searchconfig.SearchConfigList{}
	discardUnknown := !repo.strictConfigParsing
	err = protojson.UnmarshalOptions{DiscardUnknown: discardUnknown}.Unmarshal(data, &protoSearchConfigList)
	if err != nil {
		return nil, err
	}

	return &protoSearchConfigList, nil
}

// ListSearchConfigsOfType lists all config objects of a certain type
func (repo *searchConfigRepoDirectCMS) ListSearchConfigsOfType(_ context.Context, objType string,
) (*searchconfig.SearchConfigList, error) {
	client := &http.Client{}
	resp, err := client.Get(fmt.Sprintf("%s%s/search_configs?type=%s", repo.originCMS, apiPath, objType))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("could not fetch search configurations (of specific type) from CMS - got response status code %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	protoSearchConfigList := searchconfig.SearchConfigList{}
	discardUnknown := !repo.strictConfigParsing
	err = protojson.UnmarshalOptions{DiscardUnknown: discardUnknown}.Unmarshal(data, &protoSearchConfigList)
	if err != nil {
		return nil, err
	}

	return &protoSearchConfigList, nil
}

// GetSearchConfigObject returns one search config object by id
func (repo *searchConfigRepoDirectCMS) GetSearchConfigObject(_ context.Context, name string) (*searchconfig.SearchConfigObject, error) {
	client := &http.Client{}
	resp, err := client.Get(fmt.Sprintf("%s%s/search_config/%s", repo.originCMS, apiPath, name))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("could not fetch single search configuration from CMS - got response status code %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	protoSearchConfig := searchconfig.SearchConfigObject{}
	discardUnknown := !repo.strictConfigParsing
	err = protojson.UnmarshalOptions{DiscardUnknown: discardUnknown}.Unmarshal(data, &protoSearchConfig)
	if err != nil {
		return nil, err
	}

	return &protoSearchConfig, nil
}

// GetFieldsForSearchFocus return the requested search focus (error if not found)
func (repo *searchConfigRepoDirectCMS) GetFieldsForSearchFocus(ctx context.Context, requestedSearchFocus string) ([]string, error) {
	searchFociList, err := repo.ListSearchConfigsOfType(ctx, solr.MexSearchFocusType)
	if err != nil {
		return nil, err
	}
	var baseHighlightFields []string
	for _, sf := range searchFociList.GetSearchConfigs() {
		if sf.Name == requestedSearchFocus {
			baseHighlightFields = sf.Fields
		}
	}
	if len(baseHighlightFields) == 0 {
		return nil, fmt.Errorf("could not find requested search focus")
	}
	return baseHighlightFields, nil
}

// GetFieldsForAxis returns fields for an axis (ordinal or hierarchy)
func (repo *searchConfigRepoDirectCMS) GetFieldsForAxis(ctx context.Context, requestedAxis string) ([]string, error) {
	searchFociList, err := repo.ListSearchConfigsOfType(ctx, solr.MexOrdinalAxisType)
	if err != nil {
		return nil, err
	}
	var baseHighlightFields []string
	for _, sf := range searchFociList.GetSearchConfigs() {
		if sf.Name == requestedAxis {
			baseHighlightFields = sf.Fields
			break
		}
	}
	// Didn't find an ordinal axis - check the hierarchy axes
	if len(baseHighlightFields) == 0 {
		hierarchyAxisList, err := repo.ListSearchConfigsOfType(ctx, solr.MexHierarchyAxisType)
		if err != nil {
			return nil, err
		}
		for _, sf := range hierarchyAxisList.GetSearchConfigs() {
			if sf.Name == requestedAxis {
				baseHighlightFields = sf.Fields
				break
			}
		}
	}
	if len(baseHighlightFields) == 0 {
		return nil, fmt.Errorf("could not find requested ordinal axis")
	}
	return baseHighlightFields, nil
}

func (repo *searchConfigRepoDirectCMS) Purge(_ context.Context) error {
	return nil
}
