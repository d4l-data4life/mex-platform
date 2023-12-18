package index

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/d4l-data4life/mex/mex/shared/entities"
	sharedEntities "github.com/d4l-data4life/mex/mex/shared/entities"
	"github.com/d4l-data4life/mex/mex/shared/entities/erepo"

	"github.com/d4l-data4life/mex/mex/services/metadata/business/datamodel"
)

func Test_isOfFocalType(t *testing.T) {

	mockedEntityRepo := erepo.NewMockedEntityTypesRepo([]*sharedEntities.EntityType{
		{
			Name: "focalType1",
			Config: &sharedEntities.EntityTypeConfig{
				BusinessIdFieldName:   "identifier",
				AggregationEntityType: "Resource",
				AggregationAlgorithm:  "simple",
				IsFocal:               true,
			},
		},
		{
			Name: "nonFocalType`",
			Config: &sharedEntities.EntityTypeConfig{
				BusinessIdFieldName:   "identifier",
				AggregationEntityType: "Resource",
				AggregationAlgorithm:  "simple",
				IsFocal:               false,
			},
		},
		{
			Name: "focalType2",
			Config: &sharedEntities.EntityTypeConfig{
				BusinessIdFieldName:   "identifier",
				AggregationEntityType: "Resource",
				AggregationAlgorithm:  "simple",
				IsFocal:               true,
			},
		},
		{
			Name: "nonFocalType2",
			Config: &sharedEntities.EntityTypeConfig{
				BusinessIdFieldName:   "identifier",
				AggregationEntityType: "Resource",
				AggregationAlgorithm:  "simple",
				IsFocal:               false,
			},
		}})

	tests := []struct {
		name       string
		entityRepo entities.EntityRepo
		item       datamodel.ItemsWithBusinessID
		want       bool
		wantErr    bool
	}{
		{
			name:       "Returns true for focal entity type",
			entityRepo: mockedEntityRepo,
			item: datamodel.ItemsWithBusinessID{
				ItemID:              "abc123",
				CreatedAt:           pgtype.Timestamptz{Time: time.Time{}, Valid: true},
				Owner:               "blipp",
				EntityName:          "focalType2",
				BusinessID:          "123abc",
				BusinessIDFieldName: pgtype.Text{},
			},
			want: true,
		},
		{
			name:       "Returns false for non-focal entity type",
			entityRepo: mockedEntityRepo,
			item: datamodel.ItemsWithBusinessID{
				ItemID:              "abc123",
				CreatedAt:           pgtype.Timestamptz{Time: time.Time{}, Valid: true},
				Owner:               "blipp",
				EntityName:          "nonFocalType2",
				BusinessID:          "123abc",
				BusinessIDFieldName: pgtype.Text{},
			},
			want: false,
		},
		{
			name:       "Returns false for unknown entity type",
			entityRepo: mockedEntityRepo,
			item: datamodel.ItemsWithBusinessID{
				ItemID:              "abc123",
				CreatedAt:           pgtype.Timestamptz{Time: time.Time{}, Valid: true},
				Owner:               "blipp",
				EntityName:          "notAType",
				BusinessID:          "123abc",
				BusinessIDFieldName: pgtype.Text{},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := isOfFocalType(context.TODO(), tt.item, tt.entityRepo)
			if (err != nil) != tt.wantErr {
				t.Errorf("isOfFocalType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("isOfFocalType() got = %v, want %v", got, tt.want)
			}
		})
	}
}
