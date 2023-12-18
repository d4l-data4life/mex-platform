import { Command } from "commander";

import { request } from "../../utils/http";
import { getAccessToken } from "../../utils/jwt";
import { Output } from "../../utils/out";
import { Config } from "../../utils/config";

export default function (): Command {
    return new Command("list").description("list all entity types").option("--short", "short output").action(handler);
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

    const listEntityTypesResponse = await request({
        method: "GET",
        uri: cfg.effectiveInstance().mexOrigin + "/api/v0/config/files/entity_types",
        verifyCerts: false,
        headers: {
            Authorization: `Bearer ${jwt}`,
            "X-MEx-Trace": cfg.loggingTraceSecret(),
        },
    });

    if (listEntityTypesResponse.statusCode !== 200) {
        out.error(listEntityTypesResponse.body.str);
        process.exit(1);
    }

    const entityTypes = listEntityTypesResponse.body.obj.entityTypes;

    if (command.opts().short) {
        out.json(entityTypes.map((f) => f.name));
    } else {
        out.json(entityTypes);
    }
}
