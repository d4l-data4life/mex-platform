package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

// Returns a valid pgtype.TextFromString if the string is non-empty; else an invalid pgtype.TextFromString.
func TextFromString(s string) pgtype.Text {
	if s == "" {
		return pgtype.Text{}
	}
	return pgtype.Text{String: s, Valid: true}
}

func StringArray(arr []string) pgtype.Array[string] {
	return pgtype.Array[string]{
		Elements: arr,
		Dims:     []pgtype.ArrayDimension{{LowerBound: 0, Length: int32(len(arr))}},
		Valid:    true,
	}
}

func TextFromStringPtr(sp *string) pgtype.Text {
	if sp == nil {
		return pgtype.Text{}
	}
	return pgtype.Text{Valid: true, String: *sp}
}

func StringOrNil(s pgtype.Text) *string {
	if s.Valid {
		return &s.String
	}
	return nil
}

func Int32OrNil(i pgtype.Int4) *int32 {
	if i.Valid {
		return &i.Int32
	}
	return nil
}
