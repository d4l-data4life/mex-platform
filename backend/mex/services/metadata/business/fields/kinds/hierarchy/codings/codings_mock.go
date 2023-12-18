package codings

import (
	"context"
	"fmt"

	"github.com/d4l-data4life/mex/mex/shared/utils"
)

type MockedCodings struct {
	Codings map[string][]Coding
}

func (codings *MockedCodings) Reset() {
}

func (codings *MockedCodings) TransitiveClosure(_ context.Context, codeSystemName string, _ string, _ string, code string) (map[string][]Coding, error) {
	return nil, fmt.Errorf("not impl in mock")
}

func (codings *MockedCodings) Resolve(ctx context.Context, codeSystemName string, _ string, _string, code string) ([]Coding, error) {
	codeSystem, ok := codings.Codings[codeSystemName]
	if !ok {
		return nil, fmt.Errorf("mock: unknown code system: %s", codeSystemName)
	}

	codes := []Coding{}
	for _, c := range codeSystem {
		if c.Code == code {
			codes = append(codes, c)
		}
	}

	return codes, nil
}

func (codings *MockedCodings) GetCodeSystemNames(ctx context.Context) ([]string, error) {
	return utils.KeysOfMap(codings.Codings), nil
}

func (codings *MockedCodings) Close() error {
	return nil
}
