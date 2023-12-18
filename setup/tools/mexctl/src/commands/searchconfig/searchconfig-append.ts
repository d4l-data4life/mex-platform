import * as E from "fp-ts/lib/Either";
import { Command } from "commander";
import notifier from "node-notifier";

import { getAccessToken } from "../../utils/jwt";
import { Output } from "../../utils/out";
import { Config } from "../../utils/config";
import { loadData } from "../../utils/data-loader";
import { SearchConfig_T } from "../../api/defs";
import { addSearchConfigs, rebuildSolrSchema, updateSolrIndex } from "../../api/mex-api";

export default function (): Command {
    return new Command("append")
        .description("append one or more new ordinal axes or search foci")
        .requiredOption("-s, --search-configs <configs...>", "search configs file; can be specified multiple times")
        .action(handler);
}

type Options = {
    searchConfigs: string[];
};
async function handler(options: Options, command: Command): Promise<void> {
    const out = new Output(command);
    const cfg = new Config(command);

    const jwt = await getAccessToken(cfg.effectiveInstance());
    out.verbose(jwt);

    out.log(`Loading search configs from the following files:\n${options.searchConfigs.join("\n")}`);
    const searchConfigs = loadData(options.searchConfigs, SearchConfig_T);
    if (E.isLeft(searchConfigs)) {
        out.error(searchConfigs.left);
        process.exit(1);
    }

    out.log("MEx Core service origin:", cfg.mexOrigin());

    out.log(`Preparing to append ${searchConfigs.right.length} search config elements`);

    if (cfg.dryRun()) {
        out.log("dry run: exiting here");
        process.exit(0);
    }

    const config = { mexOrigin: cfg.mexOrigin(), jwt, loggingTraceSecret: cfg.loggingTraceSecret() };
    try {
        // eslint-disable-next-line no-console
        out.log("--- Starting search config append ---");

        // Add new configs
        out.log("Adding search configs");
        await addSearchConfigs(config, searchConfigs.right);

        // Propagate config to Solr
        out.log("Rebuilding Solr schema");
        await rebuildSolrSchema(config);

        // Update Solr index
        out.log("Updating search index");
        await updateSolrIndex(config);

        out.log("--- Search config elements appended ---");
    } catch (e) {
        out.log(`Failed with error ${e}`);
    }

    notifier.notify({
        title: "MEx search config append",
        message: `${searchConfigs.right.length} search config elements appended`,
    });
}
