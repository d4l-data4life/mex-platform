import * as E from "fp-ts/lib/Either";
import { Command } from "commander";
import notifier from "node-notifier";

import { getAccessToken } from "../../utils/jwt";
import { Output } from "../../utils/out";
import { Config } from "../../utils/config";
import { loadData } from "../../utils/data-loader";
import { EntityType_T, FieldDef_T, SearchConfig_T } from "../../api/defs";
import {
    addEntities,
    addFields,
    addSearchConfigs,
    rebuildSolrSchema,
    removeAllEntities,
    removeAllFields,
    removeAllSearchConfigs,
    updateSolrIndex,
} from "../../api/mex-api";

export default function (): Command {
    return new Command("configure")
        .description("set configuration")
        .requiredOption("-f, --fields <fields...>", "field defs file; can be specified multiple times")
        .requiredOption("-e, --entities <entities...>", "entity types file; can be specified multiple times")
        .requiredOption("-s, --search-configs <configs...>", "search configs file; can be specified multiple times")
        .action(handler);
}

type Options = {
    fields: string[];
    entities: string[];
    searchConfigs: string[];
};

async function handler(options: Options, command: Command): Promise<void> {
    const out = new Output(command);
    const cfg = new Config(command);

    const jwt = await getAccessToken(cfg.effectiveInstance());
    out.verbose(jwt);

    out.log(`Loading field defs from the following files:\n${options.fields.join("\n")}`);
    const fieldDefs = loadData(options.fields, FieldDef_T);
    if (E.isLeft(fieldDefs)) {
        out.error(fieldDefs.left);
        process.exit(1);
    }

    out.log(`Loading entity types from the following files:\n${options.entities.join("\n")}`);
    const entityTypes = loadData(options.entities, EntityType_T);
    if (E.isLeft(entityTypes)) {
        console.error(entityTypes.left);
        process.exit(1);
    }

    out.log(`Loading search configs from the following files:\n${options.searchConfigs.join("\n")}`);
    const searchConfigs = loadData(options.searchConfigs, SearchConfig_T);
    if (E.isLeft(searchConfigs)) {
        out.error(searchConfigs.left);
        process.exit(1);
    }

    out.log("MEx Core service origin:", cfg.mexOrigin());

    out.log(`Preparing to configure system with:
        - ${entityTypes.right.length} (non-relation) entities
        - ${searchConfigs.right.length} search configs
        - ${fieldDefs.right.length} custom fields`);

    if (cfg.dryRun()) {
        out.log("dry run: exiting here");
        process.exit(0);
    }

    const config = { mexOrigin: cfg.mexOrigin(), jwt, loggingTraceSecret: cfg.loggingTraceSecret() };
    try {
        // eslint-disable-next-line no-console
        out.log("--- Starting configuration rebuild ---");

        out.log("Removing search configs");
        await removeAllSearchConfigs(config);

        out.log("Removing existing entities");
        await removeAllEntities(config);

        out.log("Removing existing fields");
        await removeAllFields(config);

        out.log("Adding fields");
        await addFields(config, fieldDefs.right);

        out.log("Adding entities");
        await addEntities(config, entityTypes.right);

        out.log("Adding search configs");
        await addSearchConfigs(config, searchConfigs.right);

        // Propagate config to Solr
        out.log("Rebuilding Solr schema");
        await rebuildSolrSchema(config);

        // Update Solr index
        out.log("Updating search index");
        await updateSolrIndex(config);

        out.log("--- Configuration rebuild complete ---");
    } catch (e) {
        out.log(`Failed with error ${e}`);
    }

    notifier.notify({
        title: "MEx config set-up",
    });
}
