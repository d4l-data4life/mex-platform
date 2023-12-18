import { Command } from "commander";

import { getAccessToken } from "../../utils/jwt";
import { Output } from "../../utils/out";
import { Config } from "../../utils/config";
import { request } from "../../utils/http";

export default function (): Command {
    return new Command("show")
        .description("show a tree")
        .requiredOption("-n, --node-entity-type <entity type>", "entity type of the items representing tree nodes")
        .requiredOption("-l, --link-field-name <field name>", "field name establishing edges between the nodes")
        .requiredOption("-d, --display-field-name <field name>", "field name of field to render in nodes")
        .action(handler);
}

type Options = {
    nodeEntityType: string;
    linkFieldName: string;
    displayFieldName: string;

    short: boolean;
};

async function handler(options: Options, command: Command): Promise<void> {
    const out = new Output(command);
    const cfg = new Config(command);

    const jwt = await getAccessToken(cfg.effectiveInstance());
    out.verbose(jwt);
    out.verbose(options);

    if (cfg.dryRun()) {
        out.log("dry run: abort here");
        process.exit(0);
    }

    const treeResponse = await request({
        method: "POST",
        uri: `${cfg.mexOrigin()}/api/v0/metadata/tree`,
        headers: {
            Authorization: `Bearer ${jwt}`,
            "X-MEx-Trace": cfg.loggingTraceSecret(),
        },
        body: {
            nodeEntityType: options.nodeEntityType,
            linkFieldName: options.linkFieldName,
            displayFieldName: options.displayFieldName,
        },
    });

    if (treeResponse.statusCode !== 200) {
        out.error(treeResponse.body.str);
        process.exit(1);
    }

    out.json(treeResponse.body.obj);
}
