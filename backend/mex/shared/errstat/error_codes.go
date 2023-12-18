package errstat

import "google.golang.org/grpc/codes"

type MexErrorCode uint32

type MexGrpcStatusCoupling struct {
	MexErrorString string
	GrpcCode       codes.Code
}

const (
	OtherError MexErrorCode = iota
	QueryEngineCreationFailedInternal
	QueryConstructionErrorInternal
	QueryCreationFailedInternal
	InvalidClientQuery
	InvalidConfigurationClient
	SolrQueryFailedInternal
	SolrResponseProcessingInternal
	ElementNotFound
)

var codeToStatus = map[MexErrorCode]MexGrpcStatusCoupling{
	OtherError:                        {"E_OTHER", codes.Internal},
	QueryEngineCreationFailedInternal: {"E_QUERY_ENGINE_SETUP", codes.Internal},
	QueryConstructionErrorInternal:    {"E_QUERY_CONSTRUCTION", codes.Internal},
	QueryCreationFailedInternal:       {"E_QUERY_CREATION_FAILED", codes.Internal},
	InvalidClientQuery:                {"E_INVALID_CLIENT_QUERY", codes.InvalidArgument},
	InvalidConfigurationClient:        {"E_CONFIGURATION_MISSING", codes.InvalidArgument},
	SolrQueryFailedInternal:           {"E_SOLR_QUERY", codes.Internal},
	SolrResponseProcessingInternal:    {"E_SOLR_RESPONSE_PROCESSING", codes.Internal},
	ElementNotFound:                   {"E_REQUESTED_ELEMENT_NOT_FOUND", codes.NotFound},
}
