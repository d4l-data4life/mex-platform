package solr

import (
	"strings"
	"testing"
)

func Test_GetExactFieldName(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "Adds a post-fix to the input",
			input: "hello",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetExactFieldName(tt.input)
			if !strings.HasPrefix(got, tt.input) {
				t.Errorf("GetExactFieldName(): output '%s' did not start with the input '%s'", got, tt.input)
			}
		})
	}
}

func Test_GetLinkedFieldName(t *testing.T) {
	tests := []struct {
		name   string
		link   string
		target string
		want   string
	}{
		{
			name:   "Glues together link and target with the double underscore",
			link:   "hello",
			target: "there",
			want:   "hello__there",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetLinkedFieldName(tt.link, tt.target)
			if got != tt.want {
				t.Errorf("Wanted the linked field name '%s' but got '%s'", tt.want, got)
			}
		})
	}
}

func Test_GetTimestampName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "Post-fixes an underscore and the configured string to the input",
			input: "hello",
			want:  "hello" + "_" + RawValTimestampPostfix,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetRawValTimestampName(tt.input); got != tt.want {
				t.Errorf("GetRawValTimestampName(): output '%s' does not agree with expected output '%s'", got, tt.want)
			}
		})
	}
}

func Test_GetLangSpecificFieldName(t *testing.T) {
	tests := []struct {
		name      string
		fieldName string
		langCode  string
		want      string
		wantErr   bool
	}{
		{
			name:      "Accepts 'de' as language code and appends it to name with triple-underscore",
			fieldName: "hello",
			langCode:  GermanLangAbbrev,
			want:      "hello___de",
		},
		{
			name:      "Accepts 'en' as language code and appends it to name with triple-underscore",
			fieldName: "hello",
			langCode:  EnglishLangAbbrev,
			want:      "hello___en",
		},
		{
			name:      "Appends '___generic 'it to name if an empty string is passed as language code",
			fieldName: "hello",
			langCode:  GenericLangAbbrev,
			want:      "hello___" + GenericLanguageSuffix,
		},
		{
			name:      "Returns error when 'fr' is passed as language code",
			fieldName: "hello",
			langCode:  "fr",
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := GetLangSpecificFieldName(tt.fieldName, tt.langCode)
			if (gotErr != nil) != tt.wantErr {
				t.Errorf("GetLangSpecificFieldName(): wantErr '%v' but got '%s'", tt.wantErr, gotErr)
			} else if got != tt.want {
				t.Errorf("GetLangSpecificFieldName(): wanted '%s' but got '%s'", tt.want, got)
			}
		})
	}
}

func Test_GetLinkedDisplayID(t *testing.T) {
	tests := []struct {
		name      string
		fieldName string
		want      string
	}{
		{
			name:      "Camel-cased name (first letter lowercase)",
			fieldName: "aTestField",
			want:      "FIELD_A_TEST_FIELD",
		},
		{
			name:      "Camel-cased name (first letter uppercase)",
			fieldName: "AnotherTestField",
			want:      "FIELD_ANOTHER_TEST_FIELD",
		},
		{
			name:      "Mixed camel- and snake-cased name",
			fieldName: "aThird_testField",
			want:      "FIELD_A_THIRD_TEST_FIELD",
		},
		{
			name:      "Underscore followed by capital letter does not lead to double underscore",
			fieldName: "some_OtherField",
			want:      "FIELD_SOME_OTHER_FIELD",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetLinkedDisplayID(tt.fieldName); got != tt.want {
				t.Errorf("generateDisplayID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_CreateMergeId(t *testing.T) {
	tests := []struct {
		name  string
		rawId string
		want  string
	}{
		{
			name:  "Post-fixes the merge postfix, separated with the separator",
			rawId: "test",
			want:  "test#merged",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateMergeID(tt.rawId); got != tt.want {
				t.Errorf("CreateMergeID() = %v, want %v", got, tt.want)
			}
		})
	}
}
