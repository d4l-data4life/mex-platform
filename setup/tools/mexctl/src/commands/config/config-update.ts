import { Command } from "commander";

import { request } from "../../utils/http";
import { Output } from "../../utils/out";
import { Config } from "../../utils/config";

export default function (): Command {
    return new Command("update").description("checkout new config and announce changes").option("--ref", "Git ref name", "main").action(handler);
}

async function handler(options, command: Command): Promise<void> {
    const out = new Output(command);
    const cfg = new Config(command);

    out.verbose(options);

    const updateConfigResponse = await request({
        method: "POST",
        uri: cfg.mexOrigin() + "/api/v0/config/update",
        verifyCerts: false,
        headers: {
            Authorization: `apikey ${cfg.effectiveInstance().configApiKey}`,
            "X-MEx-Trace": cfg.loggingTraceSecret(),
        },
        body: {
            refName: options.ref,
            changes: ["entity_types", "search_configs", "field_defs"],
        },
    });

    if (updateConfigResponse.statusCode !== 200) {
        out.error(updateConfigResponse.body.str);
        process.exit(1);
    }
}
