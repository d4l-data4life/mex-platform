import { Command } from "commander";

import { getAccessToken } from "../../utils/jwt";
import { Output } from "../../utils/out";
import { Config } from "../../utils/config";
import { removeAllEntities, removeAllFields, removeAllItems, removeAllSearchConfigs, updateSolrIndex } from "../../api/mex-api";

export default function (): Command {
    return new Command("clean").description("remove all data and configurations").option("--drop-config", "drop Solr config", false).action(handler);
}

type Options = {
    dropConfig: boolean;
};

async function handler(options: Options, command: Command): Promise<void> {
    const out = new Output(command);
    const cfg = new Config(command);

    const jwt = await getAccessToken(cfg.effectiveInstance());
    out.verbose(jwt);

    out.log("MEx Core service origin:", cfg.mexOrigin());

    if (cfg.dryRun()) {
        out.log("dry run: exiting here");
        process.exit(0);
    }

    const config = { mexOrigin: cfg.mexOrigin(), jwt, loggingTraceSecret: cfg.loggingTraceSecret() };
    try {
        // eslint-disable-next-line no-console
        out.log("--- Starting system clean-up ---");

        // Remove old data and configs elements from MEx DB
        out.log("Removing existing items");
        await removeAllItems(config);

        if (options.dropConfig) {
            out.log("Removing search configs");
            await removeAllSearchConfigs(config);

            out.log("Removing existing entities");
            await removeAllEntities(config);

            out.log("Removing existing fields");
            await removeAllFields(config);
        }

        // Update Solr index
        out.log("Updating search index");
        await updateSolrIndex(config);

        out.log("--- System clean-up complete ---");
    } catch (e) {
        out.log(`Failed with error ${e}`);
    }
}
