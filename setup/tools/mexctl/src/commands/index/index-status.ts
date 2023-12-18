import { Command } from "commander";

import { request } from "../../utils/http";
import { getAccessToken } from "../../utils/jwt";
import { Output } from "../../utils/out";
import { Config } from "../../utils/config";

interface IndexStatus {
    clusterStatus?: ClusterStatus;
    itemCount?: number;
    message?: string;
}

interface ClusterStatus {
    collection?: string;
    health?: string;
    shards?: ShardStatus[];
}

interface ShardStatus {
    name?: string;
    health?: string;
    state?: string;
    replicas?: ReplicaStatus[];
}

interface ReplicaStatus {
    name?: string;
    state?: string;
    leader?: boolean;
}

export default function (): Command {
    return new Command("status").description("status commands").action(handler);
}

async function handler(options, command: Command): Promise<void> {
    const out = new Output(command);
    const cfg = new Config(command);

    const jwt = await getAccessToken(cfg.effectiveInstance());
    out.verbose(jwt);

    if (cfg.dryRun()) {
        out.log("dry run: exiting here");
        process.exit(0);
    }

    // First check that no element with that name & type exists
    const solrClusterStatusResponse = await request({
        method: "GET",
        uri: cfg.mexOrigin() + "/api/v0/metadata/index",
        verifyCerts: false,
        headers: {
            Authorization: `Bearer ${jwt}`,
            "X-MEx-Trace": cfg.loggingTraceSecret(),
        },
    });

    if (solrClusterStatusResponse.statusCode !== 200) {
        out.error(solrClusterStatusResponse.body);
        process.exit(1);
    }

    const indexStatus: IndexStatus = solrClusterStatusResponse.body.obj;
    out.json(indexStatus);
    prettyPrintStatus(indexStatus, out);
}

function prettyPrintStatus(indexStatus: IndexStatus, out: Output) {
    if (typeof indexStatus.message !== "undefined" && indexStatus.message.trim().length > 0) {
        out.log(`Message generated by status check: ${indexStatus.message}`);
    }
    out.log(`Index contains ${indexStatus.itemCount} item(s)`);

    if (typeof indexStatus.clusterStatus === "undefined") {
        out.log("No Solr cluster status returned!");
    } else {
        const clusterStatus = indexStatus.clusterStatus;
        out.log(`Solr collection '${clusterStatus.collection}' - overall health ${clusterStatus.health}`);
        if (!Array.isArray(clusterStatus.shards) || clusterStatus.shards.length === 0) {
            out.log("No shards found!");
        } else {
            const shards = clusterStatus.shards;
            out.log("   " + `Shards (${shards.length} in total):`);
            shards.forEach((shardStatus, i) => {
                if (!Array.isArray(shardStatus.replicas) || shardStatus.replicas.length === 0) {
                    out.log("   " + "No replicas found for shard!!");
                } else {
                    const replicas = shardStatus.replicas;
                    out.log(
                        "      " +
                            `${i + 1}) ${shardStatus.name}: overall health ${shardStatus.health}, overall state ${shardStatus.state} - ${
                                replicas.length
                            } replica(s):`,
                    );
                    const leaderNo = replicas.map((r) => r.leader).filter((l) => l).length;
                    if (leaderNo == 0) {
                        out.log("      " + "No leader replica!");
                    } else if (leaderNo > 2) {
                        out.log("      " + "Multiple leader replicas!");
                    }
                    replicas.forEach((rep, j) => {
                        const headline = rep.leader
                            ? `${j + 1}) ${rep.name} (leader) - state: ${rep.state}`
                            : `${j + 1}) ${rep.name} - state: ${rep.state}`;
                        out.log("          " + headline);
                    });
                }
            });
        }
    }
}