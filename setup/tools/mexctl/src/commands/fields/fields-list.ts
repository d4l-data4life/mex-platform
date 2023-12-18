import { Command } from "commander";

import { request } from "../../utils/http";
import { getAccessToken } from "../../utils/jwt";
import { Output } from "../../utils/out";
import { Config } from "../../utils/config";

export default function (): Command {
    return new Command("list").description("list all field definitions").option("--short", "short output").action(handler);
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
        out.log("dry run: abort here");
        process.exit(0);
    }

    const listFieldsResponse = await request({
        method: "GET",
        uri: cfg.mexOrigin() + "/api/v0/config/files/field_defs",
        verifyCerts: false,
        headers: {
            Authorization: `Bearer ${jwt}`,
            "X-MEx-Trace": cfg.loggingTraceSecret(),
        },
    });

    if (listFieldsResponse.statusCode !== 200) {
        out.error(listFieldsResponse.body.str);
        process.exit(1);
    }

    const fieldDefs = listFieldsResponse.body.obj.fieldDefs;
    if (command.opts().short) {
        out.json(fieldDefs.map((f) => f.name));
    } else {
        out.json(fieldDefs);
    }
}
