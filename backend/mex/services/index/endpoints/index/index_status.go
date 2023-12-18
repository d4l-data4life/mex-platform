package index

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/d4l-data4life/mex/mex/shared/errstat"
	L "github.com/d4l-data4life/mex/mex/shared/log"
	"github.com/d4l-data4life/mex/mex/shared/solr"

	"github.com/d4l-data4life/mex/mex/services/index/endpoints/index/pb"
)

func (svc *Service) IndexStatus(ctx context.Context, _ *pb.IndexStatusRequest) (*pb.IndexStatusResponse, error) {
	var solrStatus *pb.SolrClusterStatus
	var itemCount uint32
	var msg string
	if svc.Solr != nil {
		clusterStatus, respStatus, statusErr := svc.Solr.GetClusterStatus(ctx)
		if statusErr != nil || respStatus != 200 {
			svc.Log.Warn(ctx, L.Message("could not get Solr cluster status"))
		}
		var clusterErr error
		solrStatus, clusterErr = parseCLusterStatus(clusterStatus)
		if clusterErr != nil {
			svc.Log.Warn(ctx, L.Messagef("could not not parse cluster status (error: %s): %v", clusterErr.Error(), clusterStatus))
		}

		searchResult, statusCode, searchErr := svc.Solr.DoJSONQuery(ctx, nil, &solr.QueryBody{
			Query: "*",
			Limit: 0,
		},
		)
		if searchErr != nil {
			svc.Log.Error(ctx, L.Messagef("error executing main Solr query: %s", searchErr.Error()))
			return nil, errstat.MakeMexStatus(errstat.SolrQueryFailedInternal, fmt.Sprintf("solr query failed: %s", searchErr.Error())).Err()
		}
		if statusCode == http.StatusOK {
			itemCount = searchResult.Response.NumFound
		} else {
			msg = "Failed to retrieve no. of items in index"
		}
	} else {
		msg = "Could not find associated Solr index for service"
	}

	return &pb.IndexStatusResponse{
		ClusterStatus: solrStatus,
		ItemCount:     itemCount,
		Message:       msg,
	}, nil
}

func parseCLusterStatus(clusterStaus *solr.ClusterStatus) (*pb.SolrClusterStatus, error) {
	var solrStatus pb.SolrClusterStatus

	if len(clusterStaus.Cluster.Collections) == 0 {
		return nil, fmt.Errorf("no collection information in solrStatus response")
	}
	if len(clusterStaus.Cluster.Collections) > 1 {
		var colNames []string
		for collectionName := range clusterStaus.Cluster.Collections {
			colNames = append(colNames, collectionName)
		}
		return nil, fmt.Errorf("multiple collection returned in solrStatus response: %s", strings.Join(colNames, ","))
	}

	for collectionName, collectionStatus := range clusterStaus.Cluster.Collections {
		solrStatus.Collection = collectionName
		coercedCollectionStatus, collectionOk := collectionStatus.(map[string]interface{})
		if !collectionOk {
			return nil, fmt.Errorf("could not parse the solrStatus of the collection '%s'", collectionName)
		}

		var healthErr error
		solrStatus.Health, healthErr = getStringProperty("health", coercedCollectionStatus, false)
		if healthErr != nil {
			return nil, fmt.Errorf("collection '%s': %s", collectionName, healthErr.Error())
		}
		parsedShards, shardsErr := getObjectProperty("shards", coercedCollectionStatus, false)
		if shardsErr != nil {
			return nil, fmt.Errorf("collection '%s': %s", collectionName, shardsErr.Error())
		}
		for shardName, shardStatus := range parsedShards {
			shard, shardErr := parseShard(shardName, shardStatus)
			if shardErr != nil {
				return nil, fmt.Errorf("collection '%s', %s", collectionName, shardErr.Error())
			}
			solrStatus.Shards = append(solrStatus.Shards, shard)
		}
	}
	return &solrStatus, nil
}

func parseShard(shardName string, shardStatus interface{}) (*pb.ShardStatus, error) {
	var shard pb.ShardStatus
	shard.Name = shardName
	coercedShardStatus, shardOk := shardStatus.(map[string]interface{})
	if !shardOk {
		return nil, fmt.Errorf("could not parse the status of the shard '%s'", shardName)
	}

	var stateErr error
	shard.State, stateErr = getStringProperty("state", coercedShardStatus, false)
	if stateErr != nil {
		return nil, fmt.Errorf("shard '%s': %s", shardName, stateErr.Error())
	}
	var healthErr error
	shard.Health, healthErr = getStringProperty("health", coercedShardStatus, false)
	if healthErr != nil {
		return nil, fmt.Errorf("shard '%s': %s", shardName, healthErr.Error())
	}
	parsedReplicas, replicaErr := getObjectProperty("replicas", coercedShardStatus, false)
	if replicaErr != nil {
		return nil, fmt.Errorf("shard '%s', %s", shardName, replicaErr.Error())
	}
	for replicaName, replicaStatus := range parsedReplicas {
		replica, repErr := parseReplica(replicaName, replicaStatus)
		if repErr != nil {
			return nil, fmt.Errorf("shard '%s', %s", shardName, repErr.Error())
		}
		shard.Replicas = append(shard.Replicas, replica)
	}

	return &shard, nil
}

func parseReplica(replicaName string, replicaStatus interface{}) (*pb.ReplicaStatus, error) {
	var replica pb.ReplicaStatus
	replica.Name = replicaName
	coercedReplicaStatus, replicaOk := replicaStatus.(map[string]interface{})
	if !replicaOk {
		return nil, fmt.Errorf("could not parse the status of the replica '%s'", replicaName)
	}

	var repErr error
	replica.State, repErr = getStringProperty("state", coercedReplicaStatus, false)
	if repErr != nil {
		return nil, fmt.Errorf("replica '%s': %s", replicaName, repErr.Error())
	}
	var leaderErr error
	replica.Leader, leaderErr = getBoolProperty("leader", coercedReplicaStatus, true)
	if leaderErr != nil {
		return nil, fmt.Errorf("replica '%s': %s", replicaName, leaderErr.Error())
	}

	return &replica, nil
}

func getStringProperty(propName string, obj map[string]interface{}, setDefault bool) (string, error) {
	rawVal, valOk := obj[propName]
	if !valOk {
		if setDefault {
			return "", nil
		}
		return "", fmt.Errorf("could not property '%s'", propName)
	}
	val, castOk := rawVal.(string)
	if !castOk {
		return "", fmt.Errorf("could not cast string property '%s' (raw val: %v)", propName, rawVal)
	}
	return val, nil
}

func getBoolProperty(propName string, obj map[string]interface{}, setDefault bool) (bool, error) {
	stringVal, err := getStringProperty(propName, obj, setDefault)
	if err != nil {
		return false, fmt.Errorf("could not get bool property '%s': %s", propName, err.Error())
	}
	// If setDefault is true, stringVal will be "" if the property was missing and hence return as false
	if strings.TrimSpace(strings.ToLower(stringVal)) == "true" {
		return true, nil
	}
	return false, nil
}

func getObjectProperty(propName string, obj map[string]interface{}, setDefault bool) (map[string]interface{}, error) {
	rawVal, valOk := obj[propName]
	if !valOk {
		if setDefault {
			return nil, nil
		}
		return nil, fmt.Errorf("could not get property '%s'", propName)
	}
	val, castOk := rawVal.(map[string]interface{})
	if !castOk {
		return nil, fmt.Errorf("could not cast property '%s'", propName)
	}
	return val, nil
}
