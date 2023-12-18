import { Command } from "commander";

import { getAccessToken } from "../../utils/jwt";
import { Output } from "../../utils/out";
import { Config } from "../../utils/config";
import { deleteSolrIndex, removeAllEntities, removeAllFields, removeAllItems, removeAllSearchConfigs, updateSolrIndex } from "../../api/mex-api";

export default function (): Command {
    return new Command("delete").description("remove all items").action(handler);
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

    const config = { mexOrigin: cfg.mexOrigin(), jwt, loggingTraceSecret: cfg.loggingTraceSecret() };
    try {
        out.log("removing existing items");
        await removeAllItems(config);
        await deleteSolrIndex(config);
    } catch (e) {
        out.log(`failed with error: ${e}`);
    }
}
