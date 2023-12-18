import { Command } from "commander";

import { getAccessToken } from "../../utils/jwt";
import { Output } from "../../utils/out";
import { Config } from "../../utils/config";

export default function (): Command {
    return new Command("token").description("retrieve access token").action(handler);
}

async function handler(options, command: Command): Promise<void> {
    const out = new Output(command);
    const cfg = new Config(command);

    const token = await getAccessToken(cfg.effectiveInstance());
    out.log(token);
}
