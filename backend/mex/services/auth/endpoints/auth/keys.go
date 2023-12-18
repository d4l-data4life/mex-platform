package auth

import (
	"context"
	"encoding/json"

	pbAuth "github.com/d4l-data4life/mex/mex/services/auth/endpoints/auth/pb"
)

type rsaKeySpec struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Alg string `json:"alg"`
	Use string `json:"use"`
	E   string `json:"e"`
	N   string `json:"n"`
}

func (svc *Service) Keys(ctx context.Context, request *pbAuth.KeysRequest) (*pbAuth.KeysResponse, error) {
	publicKey, err := svc.signingKey.PublicKey()
	if err != nil {
		return nil, err
	}

	s, err := json.Marshal(publicKey)
	if err != nil {
		return nil, err
	}

	var keySpec rsaKeySpec
	err = json.Unmarshal(s, &keySpec)
	if err != nil {
		return nil, err
	}

	return &pbAuth.KeysResponse{
		Keys: []*pbAuth.KeysResponse_Key{
			{
				Kty: keySpec.Kty,
				Kid: keySpec.Kid,
				Use: keySpec.Use,
				Alg: keySpec.Alg,
				E:   keySpec.E,
				N:   keySpec.N,
			},
		},
	}, nil
}
