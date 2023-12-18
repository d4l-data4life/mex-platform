import { Command } from "commander";

import getIndexDeleteCommand from "./index/index-delete";
import getIndexUpdateCommand from "./index/index-update";
import getIndexRebuildCommand from "./index/index-rebuild";
import getIndexStatusCommand from "./index/index-status";

export default function (): Command {
    return new Command("index")
        .description("index commands")
        .addCommand(getIndexDeleteCommand())
        .addCommand(getIndexUpdateCommand())
        .addCommand(getIndexRebuildCommand())
        .addCommand(getIndexStatusCommand());
}
