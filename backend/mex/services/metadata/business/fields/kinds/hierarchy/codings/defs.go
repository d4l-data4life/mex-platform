package codings

import "context"

type Coding struct {
	Code       string
	ParentCode string
	Display    string
	Language   string
	Depth      int
}

type Codings interface {
	Reset()
	TransitiveClosure(ctx context.Context, codeSystemNameOrEntityType string, linktType string, displayFieldName string, code string) (map[string][]Coding, error)
	Resolve(ctx context.Context, codeSystemNameOrEntityType string, linktType string, displayFieldName string, code string) ([]Coding, error)
	GetCodeSystemNames(ctx context.Context) ([]string, error)
	Close() error
}
