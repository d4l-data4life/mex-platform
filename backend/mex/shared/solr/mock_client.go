package solr

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// MockClient represents the API of a specific Solr instance & core
type MockClient struct {
	Origin               string
	Core                 string
	JdbcURL              string
	JdbcUser             string
	JdbcPassword         string
	CallQueue            []string
	ValuesToReturn       ReturnVals
	AlwaysFail           bool
	UniqueID             string
	FieldsAdded          []string
	CopyFieldsAdded      []string
	DynamicFieldsAdded   []string
	FieldsRemoved        []string
	CopyFieldsRemoved    []string
	DynamicFieldsRemoved []string
	DocsUploaded         int
}

type ReturnVals struct {
	Fields        []FieldDef
	CopyFields    []CopyFieldResponse
	DynamicFields []DynamicFieldDef
}

func NewMockClient(alwaysFails bool, uniqueID string, returnVals ReturnVals) MockClient {
	return MockClient{
		Origin:               "http://example.org/solr",
		Core:                 "test",
		JdbcURL:              "xys",
		JdbcUser:             "dummy",
		JdbcPassword:         "dummy",
		CallQueue:            []string{},
		ValuesToReturn:       returnVals,
		AlwaysFail:           alwaysFails,
		UniqueID:             uniqueID,
		FieldsAdded:          []string{},
		CopyFieldsAdded:      []string{},
		DynamicFieldsAdded:   []string{},
		FieldsRemoved:        []string{},
		CopyFieldsRemoved:    []string{},
		DynamicFieldsRemoved: []string{},
		DocsUploaded:         0,
	}
}

// DoJSONQuery directly uses the JSON API of a Solr instance to carry out a query
func (solrClient *MockClient) DoJSONQuery(context.Context, url.Values, *QueryBody) (*QueryResponse, int, error) {
	solrClient.CallQueue = append(solrClient.CallQueue, "DoJsonQuery")
	return &QueryResponse{}, http.StatusOK, nil
}

func (solrClient *MockClient) DoRequest(_ context.Context, _ string, _ string, _ []byte) (int, []byte, error) {
	return 0, nil, nil
}

func (solrClient *MockClient) GetSchemaUniqueKey(_ context.Context) (string, error) {
	solrClient.CallQueue = append(solrClient.CallQueue, "GetSchemaUniqueKey")
	return solrClient.UniqueID, nil
}

func (solrClient *MockClient) GetSchemaFields(_ context.Context) ([]FieldDef, error) {
	solrClient.CallQueue = append(solrClient.CallQueue, "GetSchemaFields")
	if solrClient.AlwaysFail {
		return nil, fmt.Errorf("provoked error")
	}
	returnFields := solrClient.ValuesToReturn.Fields
	return returnFields, nil
}

func (solrClient *MockClient) GetSchemaCopyFields(_ context.Context) ([]CopyFieldResponse, error) {
	solrClient.CallQueue = append(solrClient.CallQueue, "GetSchemaCopyFields")
	if solrClient.AlwaysFail {
		return nil, fmt.Errorf("provoked error")
	}
	returnFields := solrClient.ValuesToReturn.CopyFields
	return returnFields, nil
}

func (solrClient *MockClient) GetSchemaDynamicFields(_ context.Context) ([]DynamicFieldDef, error) {
	solrClient.CallQueue = append(solrClient.CallQueue, "GetSchemaDynamicFields")
	if solrClient.AlwaysFail {
		return nil, fmt.Errorf("provoked error")
	}
	returnFields := solrClient.ValuesToReturn.DynamicFields
	return returnFields, nil
}

func (solrClient *MockClient) AddSchemaFields(_ context.Context, fields []FieldDef) error {
	solrClient.CallQueue = append(solrClient.CallQueue, "AddSchemaFields")
	for _, f := range fields {
		solrClient.FieldsAdded = append(solrClient.FieldsAdded, f.Name)
	}
	return nil
}

func (solrClient *MockClient) RemoveSchemaFields(_ context.Context, fieldNames []string) error {
	solrClient.CallQueue = append(solrClient.CallQueue, "RemoveSchemaFields")
	if solrClient.AlwaysFail {
		return fmt.Errorf("provoked error")
	}
	solrClient.FieldsRemoved = append(solrClient.FieldsRemoved, fieldNames...)
	return nil
}

func (solrClient *MockClient) AddSchemaCopyFields(_ context.Context, fields []CopyFieldDef) error {
	solrClient.CallQueue = append(solrClient.CallQueue, "AddSchemaCopyFields")
	for _, f := range fields {
		solrClient.CopyFieldsAdded = append(solrClient.CopyFieldsAdded, f.Source+"-->"+strings.Join(f.Destination, "+"))
	}
	return nil
}

func (solrClient *MockClient) RemoveSchemaCopyFields(_ context.Context, fields []RemoveCopyFieldSubBody) error {
	solrClient.CallQueue = append(solrClient.CallQueue, "RemoveSchemaCopyFields")
	if solrClient.AlwaysFail {
		return fmt.Errorf("provoked error")
	}
	for _, f := range fields {
		solrClient.CopyFieldsRemoved = append(solrClient.CopyFieldsRemoved, f.Source+"-->"+f.Dest)
	}
	return nil
}

func (solrClient *MockClient) AddSchemaDynamicFields(_ context.Context, dynamicField []DynamicFieldDef) error {
	solrClient.CallQueue = append(solrClient.CallQueue, "AddSchemaDynamicFields")
	for _, f := range dynamicField {
		solrClient.DynamicFieldsAdded = append(solrClient.DynamicFieldsAdded, f.Name)
	}
	return nil
}

func (solrClient *MockClient) RemoveSchemaDynamicFields(_ context.Context, fieldNamePatterns []string) error {
	solrClient.CallQueue = append(solrClient.CallQueue, "RemoveSchemaDynamicFields")
	if solrClient.AlwaysFail {
		return fmt.Errorf("provoked error")
	}
	solrClient.DynamicFieldsRemoved = append(solrClient.DynamicFieldsRemoved, fieldNamePatterns...)
	return nil
}

func (solrClient *MockClient) GetCollections(_ context.Context) ([]string, error) {
	if solrClient.AlwaysFail {
		return []string{}, fmt.Errorf("provoked error")
	}
	return []string{}, nil
}

func (solrClient *MockClient) CreateCollection(_ context.Context, _ string, _ string, _ uint32) error {
	if solrClient.AlwaysFail {
		return fmt.Errorf("provoked error")
	}
	return nil
}

func (solrClient *MockClient) DeleteCollection(_ context.Context, _ string) error {
	if solrClient.AlwaysFail {
		return fmt.Errorf("provoked error")
	}
	return nil
}

func (solrClient *MockClient) Ping(_ context.Context) error {
	if solrClient.AlwaysFail {
		return fmt.Errorf("provoked error")
	}
	return nil
}

func (solrClient *MockClient) GetClusterStatus(_ context.Context) (*ClusterStatus, int, error) {
	if solrClient.AlwaysFail {
		return nil, 0, fmt.Errorf("provoked error")
	}
	return &ClusterStatus{}, http.StatusOK, nil
}

func (solrClient *MockClient) DropIndex(_ context.Context) error {
	if solrClient.AlwaysFail {
		return fmt.Errorf("provoked error")
	}
	return nil
}

func (solrClient *MockClient) AddDocuments(_ context.Context, _ []string, uploadCount int) error {
	if solrClient.AlwaysFail {
		return fmt.Errorf("provoked error")
	}
	solrClient.DocsUploaded += uploadCount
	return nil
}

func (solrClient *MockClient) RemoveDocuments(_ context.Context, _ []string) error {
	if solrClient.AlwaysFail {
		return fmt.Errorf("provoked error")
	}
	return nil
}
