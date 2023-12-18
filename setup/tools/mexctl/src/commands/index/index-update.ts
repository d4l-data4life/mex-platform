import { Command } from "commander";

import { updateSolrIndex } from "../../api/mex-api";
import { getAccessToken } from "../../utils/jwt";
import { Output } from "../../utils/out";
import { Config } from "../../utils/config";

export default function (): Command {
    return new Command("update").description("delete index data, then re-index").action(handler);
}

async function handler(options, command: Command): Promise<void> {
    const out = new Output(command);
    const cfg = new Config(command);

    const jwt = await getAccessToken(cfg.effectiveInstance());
    out.verbose(jwt);

    if (cfg.dryRun()) {
        out.log("dry run: abort here");
        process.exit(0);
    }

    await updateSolrIndex({ mexOrigin: cfg.mexOrigin(), jwt, loggingTraceSecret: cfg.loggingTraceSecret() });
    out.log("index updated");
}
