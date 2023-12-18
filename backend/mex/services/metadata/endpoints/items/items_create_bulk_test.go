package items

import (
	"reflect"
	"testing"

	"github.com/d4l-data4life/mex/mex/shared/items"

	"github.com/d4l-data4life/mex/mex/services/metadata/business/canonical"
)

func Test_getEffectiveHashList(t *testing.T) {

	testItems := []*items.Item{
		{
			EntityType: "Source",
			BusinessId: "id1",
			Values:     nil,
		},
		{
			EntityType: "Resource",
			BusinessId: "id2",
			Values:     nil,
		},
	}
	hashes := make([]string, 2)
	for i, item := range testItems {
		hashes[i] = canonical.Fingerprint(item)
	}

	tests := []struct {
		name              string
		items             []*items.Item
		precomputedHashes []string
		want              []string
		wantErr           bool
	}{
		{
			name:              "If an empty list of hashes is passed, hashes are computed",
			items:             testItems,
			precomputedHashes: []string{},
			want:              hashes,
		},
		{
			name: "If as many hashes as items are passed but they are all empty strings, the hashes are computed",
			items: []*items.Item{
				{
					EntityType: "Source",
					BusinessId: "id1",
					Values:     nil,
				},
				{
					EntityType: "Resource",
					BusinessId: "id2",
					Values:     nil,
				},
			},
			precomputedHashes: []string{"", ""},
			want:              hashes,
		},
		{
			name:              "If as many hashes as items are passed and they are all non-empty strings, the passed hashes are returned",
			items:             testItems,
			precomputedHashes: []string{"a", "b"},
			want:              []string{"a", "b"},
		},
		{
			name:              "If the passed hashes are a mixture of empty and non-empty string, an error is returned",
			items:             testItems,
			precomputedHashes: []string{"a", ""},
			wantErr:           true,
		},
		{
			name:              "If the no. of passed hashes does not mix the no. of items, an error is returned",
			items:             testItems,
			precomputedHashes: []string{"a", "b", "c"},
			wantErr:           true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getEffectiveHashList(tt.items, tt.precomputedHashes)
			if (err != nil) != tt.wantErr {
				t.Errorf("getEffectiveHashList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getEffectiveHashList() got = %v, want %v", got, tt.want)
			}
		})
	}
}
