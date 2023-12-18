package index

import (
	"reflect"
	"sort"
	"testing"

	"github.com/d4l-data4life/mex/mex/shared/solr"

	"github.com/d4l-data4life/mex/mex/services/index/endpoints/index/pb"
)

func Test_parseCLusterStatus(t *testing.T) {
	tests := []struct {
		name         string
		clusterStaus *solr.ClusterStatus
		want         *pb.SolrClusterStatus
		wantErr      bool
	}{
		{
			name: "Healthy status parsed OK",
			clusterStaus: &solr.ClusterStatus{
				Cluster: solr.ClusterInfo{
					Collections: map[string]interface{}{
						"mex_d4l": map[string]interface{}{
							"maxShardsPerNode": 1,
							"router": map[string]interface{}{
								"name": "compositeId",
							},
							"replicationFactor": 1,
							"znodeVersion":      11,
							"autoCreated":       true,
							"configName":        "my_config",
							"health":            "GREEN",
							"aliases":           []string{"someAlias"},
							"shards": map[string]interface{}{
								"shard1": map[string]interface{}{
									"range":  "80000000-ffffffff",
									"state":  "active",
									"health": "GREEN",
									"replicas": map[string]interface{}{
										"core_node1": map[string]interface{}{
											"state":     "active",
											"core":      "nex_d4l",
											"node_name": "node1",
											"base_url":  "http://127.0.1.1:8983/solr11",
										},
										"core_node2": map[string]interface{}{
											"state":     "active",
											"core":      "nex_d4l",
											"node_name": "node2",
											"base_url":  "http://127.0.1.1:8983/solr12",
											"leader":    "true",
										},
									},
								},
								"shard2": map[string]interface{}{
									"range":  "80000000-aaaaaa",
									"state":  "active",
									"health": "GREEN",
									"replicas": map[string]interface{}{
										"core_node1": map[string]interface{}{
											"state":     "active",
											"core":      "mex_d4l",
											"node_name": "node1",
											"base_url":  "http://127.0.1.1:8983/solr21",
											"leader":    "true",
										},
										"core_node2": map[string]interface{}{
											"state":     "active",
											"core":      "mex_d4l",
											"node_name": "node2",
											"base_url":  "http://127.0.1.1:8983/solr22",
										},
									},
								},
							},
						},
					},
					Aliases: map[string]interface{}{
						"someAlias": "mex_d4l",
					},
					Roles: map[string]interface{}{
						"overseer": []string{"7.8.9"},
					},
					LiveNodes: []string{"1.2.3", "4.5.6"},
				},
			},
			want: &pb.SolrClusterStatus{
				Collection: "mex_d4l",
				Health:     "GREEN",
				Shards: []*pb.ShardStatus{
					{
						Name:   "shard1",
						Health: "GREEN",
						State:  "active",
						Replicas: []*pb.ReplicaStatus{
							{
								Name:   "core_node1",
								State:  "active",
								Leader: false,
							},
							{
								Name:   "core_node2",
								State:  "active",
								Leader: true,
							},
						},
					},
					{
						Name:   "shard2",
						Health: "GREEN",
						State:  "active",
						Replicas: []*pb.ReplicaStatus{
							{
								Name:   "core_node1",
								State:  "active",
								Leader: true,
							},
							{
								Name:   "core_node2",
								State:  "active",
								Leader: false,
							},
						},
					},
				},
			},
		},
		{
			name: "Unhealthy status parsed OK",
			clusterStaus: &solr.ClusterStatus{
				Cluster: solr.ClusterInfo{
					Collections: map[string]interface{}{
						"mex_d4l": map[string]interface{}{
							"maxShardsPerNode": 1,
							"router": map[string]interface{}{
								"name": "compositeId",
							},
							"replicationFactor": 1,
							"znodeVersion":      11,
							"autoCreated":       true,
							"configName":        "my_config",
							"health":            "RED",
							"aliases":           []string{"someAlias"},
							"shards": map[string]interface{}{
								"shard1": map[string]interface{}{
									"range":  "80000000-ffffffff",
									"state":  "down",
									"health": "RED",
									"replicas": map[string]interface{}{
										"core_node1": map[string]interface{}{
											"state":     "down",
											"core":      "nex_d4l",
											"node_name": "node1",
											"base_url":  "http://127.0.1.1:8983/solr11",
										},
										"core_node2": map[string]interface{}{
											"state":     "down",
											"core":      "nex_d4l",
											"node_name": "node2",
											"base_url":  "http://127.0.1.1:8983/solr12",
										},
									},
								},
								"shard2": map[string]interface{}{
									"range":  "80000000-aaaaaa",
									"state":  "active",
									"health": "YELLOW",
									"replicas": map[string]interface{}{
										"core_node1": map[string]interface{}{
											"state":     "recovering",
											"core":      "mex_d4l",
											"node_name": "node1",
											"base_url":  "http://127.0.1.1:8983/solr21",
											"leader":    "true",
										},
										"core_node2": map[string]interface{}{
											"state":     "recovering",
											"core":      "mex_d4l",
											"node_name": "node2",
											"base_url":  "http://127.0.1.1:8983/solr22",
										},
									},
								},
							},
						},
					},
					Aliases: map[string]interface{}{
						"someAlias": "mex_d4l",
					},
					Roles: map[string]interface{}{
						"overseer": []string{"7.8.9"},
					},
					LiveNodes: []string{"1.2.3", "4.5.6"},
				},
			},
			want: &pb.SolrClusterStatus{
				Collection: "mex_d4l",
				Health:     "RED",
				Shards: []*pb.ShardStatus{
					{
						Name:   "shard1",
						Health: "RED",
						State:  "down",
						Replicas: []*pb.ReplicaStatus{
							{
								Name:   "core_node1",
								State:  "down",
								Leader: false,
							},
							{
								Name:   "core_node2",
								State:  "down",
								Leader: false,
							},
						},
					},
					{
						Name:   "shard2",
						Health: "YELLOW",
						State:  "active",
						Replicas: []*pb.ReplicaStatus{
							{
								Name:   "core_node1",
								State:  "recovering",
								Leader: true,
							},
							{
								Name:   "core_node2",
								State:  "recovering",
								Leader: false,
							},
						},
					},
				},
			},
		},
		{
			name: "Nil in shard field if no shards found",
			clusterStaus: &solr.ClusterStatus{
				Cluster: solr.ClusterInfo{
					Collections: map[string]interface{}{
						"mex_d4l": map[string]interface{}{
							"maxShardsPerNode": 1,
							"router": map[string]interface{}{
								"name": "compositeId",
							},
							"replicationFactor": 1,
							"znodeVersion":      11,
							"autoCreated":       true,
							"configName":        "my_config",
							"health":            "RED",
							"aliases":           []string{"someAlias"},
							"shards":            map[string]interface{}{},
						},
					},
					Aliases: map[string]interface{}{
						"someAlias": "mex_d4l",
					},
					Roles: map[string]interface{}{
						"overseer": []string{"7.8.9"},
					},
					LiveNodes: []string{"1.2.3", "4.5.6"},
				},
			},
			want: &pb.SolrClusterStatus{
				Collection: "mex_d4l",
				Health:     "RED",
				Shards:     nil,
			},
		},
		{
			name: "Error if replica is missing state",
			clusterStaus: &solr.ClusterStatus{
				Cluster: solr.ClusterInfo{
					Collections: map[string]interface{}{
						"mex_d4l": map[string]interface{}{
							"maxShardsPerNode": 1,
							"router": map[string]interface{}{
								"name": "compositeId",
							},
							"replicationFactor": 1,
							"znodeVersion":      11,
							"autoCreated":       true,
							"configName":        "my_config",
							"health":            "GREEN",
							"aliases":           []string{"someAlias"},
							"shards": map[string]interface{}{
								"shard1": map[string]interface{}{
									"range":  "80000000-ffffffff",
									"state":  "active",
									"health": "GREEN",
									"replicas": map[string]interface{}{
										"core_node1": map[string]interface{}{
											"state":     "active",
											"core":      "nex_d4l",
											"node_name": "node1",
											"base_url":  "http://127.0.1.1:8983/solr11",
										},
										"core_node2": map[string]interface{}{
											"state":     "active",
											"core":      "nex_d4l",
											"node_name": "node2",
											"base_url":  "http://127.0.1.1:8983/solr12",
											"leader":    "true",
										},
									},
								},
								"shard2": map[string]interface{}{
									"range":  "80000000-aaaaaa",
									"state":  "active",
									"health": "GREEN",
									"replicas": map[string]interface{}{
										"core_node1": map[string]interface{}{
											"state":     "active",
											"core":      "mex_d4l",
											"node_name": "node1",
											"base_url":  "http://127.0.1.1:8983/solr21",
											"leader":    "true",
										},
										"core_node2": map[string]interface{}{
											"core":      "mex_d4l",
											"node_name": "node2",
											"base_url":  "http://127.0.1.1:8983/solr22",
										},
									},
								},
							},
						},
					},
					Aliases: map[string]interface{}{
						"someAlias": "mex_d4l",
					},
					Roles: map[string]interface{}{
						"overseer": []string{"7.8.9"},
					},
					LiveNodes: []string{"1.2.3", "4.5.6"},
				},
			},
			wantErr: true,
		},
		{
			name: "Error if shard is missing state",
			clusterStaus: &solr.ClusterStatus{
				Cluster: solr.ClusterInfo{
					Collections: map[string]interface{}{
						"mex_d4l": map[string]interface{}{
							"maxShardsPerNode": 1,
							"router": map[string]interface{}{
								"name": "compositeId",
							},
							"replicationFactor": 1,
							"znodeVersion":      11,
							"autoCreated":       true,
							"configName":        "my_config",
							"health":            "GREEN",
							"aliases":           []string{"someAlias"},
							"shards": map[string]interface{}{
								"shard1": map[string]interface{}{
									"range":  "80000000-ffffffff",
									"state":  "active",
									"health": "GREEN",
									"replicas": map[string]interface{}{
										"core_node1": map[string]interface{}{
											"state":     "active",
											"core":      "nex_d4l",
											"node_name": "node1",
											"base_url":  "http://127.0.1.1:8983/solr11",
										},
										"core_node2": map[string]interface{}{
											"state":     "active",
											"core":      "nex_d4l",
											"node_name": "node2",
											"base_url":  "http://127.0.1.1:8983/solr12",
											"leader":    "true",
										},
									},
								},
								"shard2": map[string]interface{}{
									"range":  "80000000-aaaaaa",
									"health": "GREEN",
									"replicas": map[string]interface{}{
										"core_node1": map[string]interface{}{
											"state":     "active",
											"core":      "mex_d4l",
											"node_name": "node1",
											"base_url":  "http://127.0.1.1:8983/solr21",
											"leader":    "true",
										},
										"core_node2": map[string]interface{}{
											"core":      "mex_d4l",
											"node_name": "node2",
											"base_url":  "http://127.0.1.1:8983/solr22",
											"state":     "active",
										},
									},
								},
							},
						},
					},
					Aliases: map[string]interface{}{
						"someAlias": "mex_d4l",
					},
					Roles: map[string]interface{}{
						"overseer": []string{"7.8.9"},
					},
					LiveNodes: []string{"1.2.3", "4.5.6"},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseCLusterStatus(tt.clusterStaus)
			sortStatus(tt.want)
			sortStatus(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseClusterStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseCLusterStatus() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func sortStatus(st *pb.SolrClusterStatus) {
	if st == nil || len(st.Shards) == 0 {
		return
	}
	for _, wantShard := range st.Shards {
		sort.Slice(wantShard.Replicas, func(i, j int) bool {
			return wantShard.Replicas[i].Name > wantShard.Replicas[j].Name
		})
	}
	sort.Slice(st.Shards, func(i, j int) bool {
		return st.Shards[i].Name > st.Shards[j].Name
	})
}
