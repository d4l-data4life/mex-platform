package frepo

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"google.golang.org/protobuf/encoding/protojson"

	sharedFields "github.com/d4l-data4life/mex/mex/shared/fields"

	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields"
	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields/hooks"
	"github.com/d4l-data4life/mex/mex/services/metadata/business/fields/linked"
)

const apiPath = "/api/v0/config/files"

type fieldDefsRepoDirectCMS struct {
	originCMS           string
	parsers             hooks.FieldDefinitionHooks
	strictConfigParsing bool
}

func NewDirectCMSFieldDefsRepo(originCMS string, parsers hooks.FieldDefinitionHooks, strictConfigParsing bool) fields.FieldRepo {
	return &fieldDefsRepoDirectCMS{
		originCMS:           originCMS,
		parsers:             parsers,
		strictConfigParsing: strictConfigParsing,
	}
}

// GetFieldDefNames returns the names of all defined fields, including linked fields
func (repo *fieldDefsRepoDirectCMS) GetFieldDefNames(ctx context.Context) ([]string, error) {
	fieldDefs, err := repo.ListFieldDefs(ctx)
	if err != nil {
		return nil, err
	}
	fieldNames := make([]string, len(fieldDefs))
	for i, fd := range fieldDefs {
		fieldNames[i] = fd.Name()
	}
	return fieldNames, nil
}

func (repo *fieldDefsRepoDirectCMS) GetFieldDefByName(ctx context.Context, fieldName string) (fields.BaseFieldDef, error) {
	/*
		Pulling all fields to return just one is relatively wasteful. However, the alternative is building potentially
		brittle logic for detecting linked fields & extracting the base field name. Also, since we typically use a
		cached repo on deployed system, this is unlikely to have a significant practical impact.
	*/
	fieldDefs, err := repo.ListFieldDefs(ctx)
	if err != nil {
		return nil, err
	}
	for _, fd := range fieldDefs {
		if fd.Name() == fieldName {
			return fd, nil
		}
	}
	return nil, fmt.Errorf("field not found")
}

// GetFieldDefsByKind returns all fields (including linked fields) of a given kind.
func (repo *fieldDefsRepoDirectCMS) GetFieldDefsByKind(ctx context.Context, fieldKind string) ([]fields.BaseFieldDef, error) {
	unfilteredFieldDefs, err := repo.ListFieldDefs(ctx)
	if err != nil {
		return nil, err
	}
	var fieldDefs []fields.BaseFieldDef
	for _, fd := range unfilteredFieldDefs {
		if fd.Kind() == fieldKind {
			fieldDefs = append(fieldDefs, fd)
		}
	}
	return fieldDefs, nil
}

// ListFieldDefs returns the configurations of all configured fields, including linked fields
func (repo *fieldDefsRepoDirectCMS) ListFieldDefs(ctx context.Context) ([]fields.BaseFieldDef, error) {
	client := &http.Client{}
	resp, err := client.Get(fmt.Sprintf("%s%s/field_defs", repo.originCMS, apiPath))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("could not fetch fields configurations from CMS - got response status code %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	protoFieldDefs := sharedFields.FieldDefList{}
	discardUnknown := !repo.strictConfigParsing
	err = protojson.UnmarshalOptions{DiscardUnknown: discardUnknown}.Unmarshal(data, &protoFieldDefs)
	if err != nil {
		return nil, err
	}

	// List all user-configured fields (that are not linked fields)
	fieldDefs := make([]fields.BaseFieldDef, len(protoFieldDefs.FieldDefs))
	for i, protoFieldDef := range protoFieldDefs.FieldDefs {
		fieldDef, parseErr := repo.parseField(ctx, protoFieldDef)
		if parseErr != nil {
			return nil, parseErr
		}
		fieldDefs[i] = fieldDef
	}

	// Add linked field definitions
	linkedFieldDefs, err := linked.GetLinkedFieldDefs(fieldDefs)
	if err != nil {
		return nil, fmt.Errorf("could not generate linked field configs: %w", err)
	}
	fieldDefs = append(fieldDefs, linkedFieldDefs...)
	return fieldDefs, nil
}

// parseField converts the protobuf object for a field def into a field def object with the standard interface
func (repo *fieldDefsRepoDirectCMS) parseField(ctx context.Context, protoFieldDef *sharedFields.FieldDef) (fields.BaseFieldDef, error) {
	parser := repo.parsers.GetHook(protoFieldDef.Kind)
	if parser == nil {
		return nil, fmt.Errorf("could not find fields def: %s", protoFieldDef.Kind)
	}

	fieldDef, err := parser.ValidateDefinition(ctx, protoFieldDef)
	if err != nil {
		return nil, fmt.Errorf("could not validate field definition: %s/%s (%w)", protoFieldDef.Name, protoFieldDef.Kind, err)
	}

	if fieldDef == nil {
		panic("parsed result is nil: " + protoFieldDef.Name)
	}
	return fieldDef, nil
}

func (repo *fieldDefsRepoDirectCMS) Purge(_ context.Context) error { return nil }
