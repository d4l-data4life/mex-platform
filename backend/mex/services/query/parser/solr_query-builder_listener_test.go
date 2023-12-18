package parser

import (
	"reflect"
	"testing"

	"github.com/d4l-data4life/mex/mex/shared/solr"
)

func TestAddEditDistance(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		distance uint32
		want     string
	}{
		{
			name:     "If edit distance is zero, the input is returned unchanged",
			input:    "text",
			distance: 0,
			want:     "text",
		},
		{
			name:     "If passed string is shorter than the lower cut-off, no fuzzy operator is added",
			input:    "in",
			distance: 1,
			want:     "in",
		},
		{
			name:     "If the passed string length is equal to the lower cut-off, no fuzzy operator is added",
			input:    "ten",
			distance: 2,
			want:     "ten",
		},
		{
			name:     "If the passed string length is 1 above the lower cut-off, a distance-1 fuzzing operator is added (assuming the upper cut-off is not too close)",
			input:    "texts",
			distance: 2,
			want:     "texts~1",
		},
		{
			name:     "If the passed string length is equal to the upper cut-off, the full edit-distance is applied",
			input:    "textuality",
			distance: 2,
			want:     "textuality~2",
		},
		{
			name:     "If the passed string length is above the upper cut-off, the full edit-distance is applied",
			input:    "incredulously",
			distance: 2,
			want:     "incredulously~2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := addEditDistance(tt.input, tt.distance); got != tt.want {
				t.Errorf("AddEditDistance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name             string
		matchingOpConfig MatchingOpsConfig
		want             *SolrQueryBuilderListener
		wantErr          bool
	}{
		{
			name: "Edit distance larger than allowed max causes error",
			matchingOpConfig: MatchingOpsConfig{
				Term: []MatchingFieldConfig{
					{
						FieldName:       "A",
						MaxEditDistance: solr.MaxEditDistance + 1,
					},
				},
				Phrase: []MatchingFieldConfig{
					{
						FieldName:       "B",
						MaxEditDistance: 0,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Setting an edit distance > 0 for a phrase field causes error",
			matchingOpConfig: MatchingOpsConfig{
				Term: []MatchingFieldConfig{
					{
						FieldName:       "A",
						MaxEditDistance: 0,
					},
				},
				Phrase: []MatchingFieldConfig{
					{
						FieldName:       "B",
						MaxEditDistance: 1,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Adding a term match field with no name causes an error",
			matchingOpConfig: MatchingOpsConfig{
				Term: []MatchingFieldConfig{
					{
						FieldName:       "",
						MaxEditDistance: 0,
					},
				},
				Phrase: []MatchingFieldConfig{
					{
						FieldName:       "B",
						MaxEditDistance: 0,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Adding a phrase match field with no name causes an error",
			matchingOpConfig: MatchingOpsConfig{
				Term: []MatchingFieldConfig{
					{
						FieldName:       "A",
						MaxEditDistance: 0,
					},
				},
				Phrase: []MatchingFieldConfig{
					{
						FieldName:       "",
						MaxEditDistance: 0,
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := NewListener(tt.matchingOpConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewListener() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewListener() got = %v, want %v", got, tt.want)
			}
		})
	}
}
