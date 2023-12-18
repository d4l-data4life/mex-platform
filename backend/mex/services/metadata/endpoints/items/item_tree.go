package items

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"

	"github.com/d4l-data4life/mex/mex/shared/coll/forest"
	"github.com/d4l-data4life/mex/mex/shared/errstat"
	L "github.com/d4l-data4life/mex/mex/shared/log"

	"github.com/d4l-data4life/mex/mex/services/metadata/business/datamodel"
	pbItems "github.com/d4l-data4life/mex/mex/services/metadata/endpoints/items/pb"
)

func (svc *Service) ComputeItemsTree(ctx context.Context, request *pbItems.ComputeItemsTreeRequest) (*pbItems.ComputeItemsTreeResponse, error) {
	if request.NodeEntityType == "" || request.LinkFieldName == "" {
		return nil, errstat.MakeGRPCStatus(codes.InvalidArgument, "specify node entity type and link field name", request).Err()
	}

	queries := datamodel.New(svc.DB)

	begin := time.Now()
	nodes, err := queries.DbReadTreeNodes(ctx, request.NodeEntityType, request.LinkFieldName, request.DisplayFieldName)
	if err != nil {
		return nil, err
	}

	svc.Log.Info(ctx, L.Messagef("# of nodes: %d, query duration: %v", nodes.Size(), time.Since(begin)))

	// Determine depth values (formerly done in SQL query, but that was too slow).
	f := forest.NewForestWriter[*pbItems.ComputeItemsTreeResponse_TreeNode]()
	for _, k := range nodes.Keys() {
		n := nodes.GetByKeyOrNil(k)
		f.Add(k, emptyIfNil((*n).ParentNodeId), *n)
	}

	r := f.Seal()
	if r == nil {
		return nil, errstat.MakeGRPCStatus(codes.Aborted, "invalid forest", request).Err()
	}

	response := pbItems.ComputeItemsTreeResponse{
		Nodes: make([]*pbItems.ComputeItemsTreeResponse_TreeNode, nodes.Size()),
	}

	for i := 0; i < r.Size(); i++ {
		payload, _ := r.GetByIndex(i)

		response.Nodes[i] = *payload
		response.Nodes[i].Depth = int32(r.MustDepth((*payload).NodeId))
	}

	svc.Log.Info(ctx, L.Messagef("# of items: %d", len(response.Nodes)))
	return &response, nil
}

func emptyIfNil(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
