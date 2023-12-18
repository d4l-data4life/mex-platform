package erepo

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/d4l-data4life/mex/mex/shared/entities"
)

const apiPath = "/api/v0/config/files"

type entitiesRepoDirectCMS struct {
	originCMS           string
	strictConfigParsing bool
}

func NewDirectCMSEntityTypesRepo(originCMS string, strictConfigParsing bool) entities.EntityRepo {
	return &entitiesRepoDirectCMS{
		originCMS:           originCMS,
		strictConfigParsing: strictConfigParsing,
	}
}

func (repo *entitiesRepoDirectCMS) GetEntityTypeNames(ctx context.Context, focalOnly bool) ([]string, error) {
	entityTypes, err := repo.getEntities(ctx, focalOnly)
	if err != nil {
		return nil, err
	}

	retTypeNames := make([]string, len(entityTypes))
	for i, entityType := range entityTypes {
		retTypeNames[i] = entityType.Name
	}

	return retTypeNames, nil
}

func (repo *entitiesRepoDirectCMS) GetEntityType(_ context.Context, entityTypeName string) (*entities.EntityType, error) {
	client := &http.Client{}
	resp, err := client.Get(fmt.Sprintf("%s%s/entity_types/%s", repo.originCMS, apiPath, entityTypeName))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("could not fetch single entity configuration from CMS - got response status code %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	protoEntityType := entities.EntityType{}
	discardUnknown := !repo.strictConfigParsing
	err = protojson.UnmarshalOptions{DiscardUnknown: discardUnknown}.Unmarshal(data, &protoEntityType)
	if err != nil {
		return nil, err
	}

	return &protoEntityType, nil
}

func (repo *entitiesRepoDirectCMS) ListEntityTypes(ctx context.Context, focalOnly bool) ([]*entities.EntityType, error) {
	retEntityTypes, err := repo.getEntities(ctx, focalOnly)
	if err != nil {
		return nil, err
	}

	return retEntityTypes, nil
}

// getEntities return all stored entities (or only the focal ones, depending on passed flag)
func (repo *entitiesRepoDirectCMS) getEntities(_ context.Context, focalOnly bool) ([]*entities.EntityType, error) {
	client := &http.Client{}
	resp, err := client.Get(fmt.Sprintf("%s%s/entity_types", repo.originCMS, apiPath))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("could not fetch entity configurations from CMS - got response status code %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	protoEntityTypes := entities.EntityTypeList{}
	discardUnknown := !repo.strictConfigParsing
	err = protojson.UnmarshalOptions{DiscardUnknown: discardUnknown}.Unmarshal(data, &protoEntityTypes)
	if err != nil {
		return nil, err
	}

	if !focalOnly {
		return protoEntityTypes.EntityTypes, nil
	}

	focalEntityTypes := []*entities.EntityType{}
	for _, et := range protoEntityTypes.EntityTypes {
		if et.Config.IsFocal {
			focalEntityTypes = append(focalEntityTypes, et)
		}
	}

	return focalEntityTypes, nil
}

func (*entitiesRepoDirectCMS) Purge(ctx context.Context) error { return nil }
