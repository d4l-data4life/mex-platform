import { Command } from "commander";
import notifier from "node-notifier";
import * as E from "fp-ts/lib/Either";

import { request } from "../../utils/http";
import { getAccessToken } from "../../utils/jwt";
import { Output } from "../../utils/out";
import { Config } from "../../utils/config";
import { loadData } from "../../utils/data-loader";
import { EntityTypesList, EntityType_T } from "../../api/defs";
import { addEntities, rebuildSolrSchema, updateSolrIndex } from "../../api/mex-api";

export default function (): Command {
    return new Command("append")
        .description("append one or more new entity types")
        .requiredOption("-e, --entities <entities...>", "entity types file; can be specified multiple times")
        .action(handler);
}

type Options = {
    entities: string[];
};

async function handler(options: Options, command: Command): Promise<void> {
    const out = new Output(command);
    const cfg = new Config(command);

    const jwt = await getAccessToken(cfg.effectiveInstance());
    out.verbose(jwt);

    out.log(`Loading entity types from following files:\n${options.entities.join("\n")}`);
    const entityTypes = loadData(options.entities, EntityType_T);
    if (E.isLeft(entityTypes)) {
        out.error(entityTypes.left);
        process.exit(1);
    }

    out.log("MEx Core service origin:", cfg.mexOrigin());

    out.log(`Preparing to append ${entityTypes.right.length} entity types elements`);

    if (cfg.dryRun()) {
        out.log("dry run: exiting here");
        process.exit(0);
    }

    const config = { mexOrigin: cfg.mexOrigin(), jwt, loggingTraceSecret: cfg.loggingTraceSecret() };
    try {
        out.log("--- Starting entity type append ---");

        // First check that no element with that name & type exists
        const listEntityTypesResponse = await request({
            method: "GET",
            uri: cfg.mexOrigin() + "/cms/entity_types",
            verifyCerts: false,
            headers: {
                Authorization: `Bearer ${jwt}`,
            },
        });

        if (listEntityTypesResponse.statusCode !== 200) {
            out.error(listEntityTypesResponse.body);
            process.exit(1);
        }

        for (const existingEntityType of (listEntityTypesResponse.body.obj as EntityTypesList).entityTypes) {
            for (const newEntityType of entityTypes.right) {
                if (existingEntityType.name === newEntityType.name) {
                    out.error(`An entity type of name '${existingEntityType.name}' already exists`);
                    process.exit(1);
                }
            }
        }

        // Add new configs
        await addEntities(config, entityTypes.right);

        // Propagate config to Solr
        out.log("Rebuilding Solr schema");
        await rebuildSolrSchema(config);

        // Update Solr index
        out.log("Updating search index");
        await updateSolrIndex(config);

        out.log("--- Entity types appended ---");
    } catch (e) {
        out.log(`Failed with error ${e}`);
    }

    notifier.notify({
        title: "MEx entity types append",
        message: `${entityTypes.right.length} entity types appended`,
    });
}
