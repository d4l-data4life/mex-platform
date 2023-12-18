import { Command } from "commander";

import { getAccessToken } from "../../utils/jwt";
import { Output } from "../../utils/out";
import { Config } from "../../utils/config";
import { request } from "../../utils/http";

export default function (): Command {
    return new Command("list").description("list search config elements").option("--short", "short output").action(handler).action(handler);
}

type Options = {
    short: boolean;
};

async function handler(options: Options, command: Command): Promise<void> {
    const out = new Output(command);
    const cfg = new Config(command);

    const jwt = await getAccessToken(cfg.effectiveInstance());
    out.verbose(jwt);

    if (cfg.dryRun()) {
        out.log("dry run: exiting here");
        process.exit(0);
    }

    // First check that no element with that name & type exists
    const listSearchConfigResponse = await request({
        method: "GET",
        uri: cfg.mexOrigin() + "/api/v0/config/files/search_configs",
        verifyCerts: false,
        headers: {
            Authorization: `Bearer ${jwt}`,
            "X-MEx-Trace": cfg.loggingTraceSecret(),
        },
    });

    if (listSearchConfigResponse.statusCode !== 200) {
        out.error(listSearchConfigResponse.body);
        process.exit(1);
    }

    if (command.opts().short) {
        out.log(listSearchConfigResponse.body.obj.searchConfigs.map((f) => f.name));
    } else {
        out.log(listSearchConfigResponse.body.obj.str);
    }
}
