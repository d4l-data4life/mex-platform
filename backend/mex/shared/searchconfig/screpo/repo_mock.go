package screpo

import (
	"context"
	"fmt"

	"github.com/d4l-data4life/mex/mex/shared/searchconfig"
	"github.com/d4l-data4life/mex/mex/shared/solr"
)

type mockSearchConfigRepo struct {
	content map[string]*searchconfig.SearchConfigObject
}

func NewMockSearchConfigRepo(content []*searchconfig.SearchConfigObject) searchconfig.SearchConfigRepo {
	contentByName := make(map[string]*searchconfig.SearchConfigObject)
	for _, sc := range content {
		contentByName[sc.Name] = sc
	}
	return &mockSearchConfigRepo{
		content: contentByName,
	}
}

func (repo *mockSearchConfigRepo) ListSearchConfigs(_ context.Context) (*searchconfig.SearchConfigList, error) {
	var scList *searchconfig.SearchConfigList
	for _, elem := range repo.content {
		scList.SearchConfigs = append(scList.SearchConfigs, elem)
	}
	return scList, nil
}

func (repo *mockSearchConfigRepo) ListSearchConfigsOfType(_ context.Context, objType string,
) (*searchconfig.SearchConfigList, error) {
	scList := searchconfig.SearchConfigList{
		SearchConfigs: []*searchconfig.SearchConfigObject{},
	}
	for _, elem := range repo.content {
		if elem.Type == objType {
			scList.SearchConfigs = append(scList.SearchConfigs, elem)
		}
	}
	return &scList, nil
}

func (repo *mockSearchConfigRepo) GetSearchConfigObject(_ context.Context, name string) (*searchconfig.SearchConfigObject, error) {
	elem, ok := repo.content[name]
	if !ok {
		return nil, fmt.Errorf("element not found")
	}
	return elem, nil
}

func (repo *mockSearchConfigRepo) GetFieldsForSearchFocus(_ context.Context, requestedSearchFocus string) ([]string, error) {
	for _, elem := range repo.content {
		if elem.Type == solr.MexSearchFocusType && elem.Name == requestedSearchFocus {
			return elem.Fields, nil
		}
	}
	return nil, fmt.Errorf("search focus not found")
}

func (repo *mockSearchConfigRepo) GetFieldsForAxis(_ context.Context, requestedOrdinalAxis string) ([]string, error) {
	for _, elem := range repo.content {
		if (elem.Type == solr.MexOrdinalAxisType || elem.Type == solr.MexHierarchyAxisType) && elem.
			Name == requestedOrdinalAxis {
			return elem.Fields, nil
		}
	}
	return nil, fmt.Errorf("ordinal axis not found")
}

func (*mockSearchConfigRepo) Purge(ctx context.Context) error { return nil }
