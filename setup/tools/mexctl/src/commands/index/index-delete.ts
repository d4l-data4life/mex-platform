import { Command } from "commander";

import { request } from "../../utils/http";
import { getAccessToken } from "../../utils/jwt";
import { Output } from "../../utils/out";
import { Config } from "../../utils/config";

export default function (): Command {
    return new Command("delete").description("delete the indexed data but keep the schema").action(handler);
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

    const deleteResponse = await request({
        method: "DELETE",
        uri: cfg.mexOrigin() + "/api/v0/metadata/index",
        verifyCerts: false,
        headers: {
            Authorization: `Bearer ${jwt}`,
            "X-MEx-Trace": cfg.loggingTraceSecret(),
        },
    });

    if (deleteResponse.statusCode !== 204) {
        out.error(deleteResponse.body);
    } else {
        out.log("index deleted");
    }
}
