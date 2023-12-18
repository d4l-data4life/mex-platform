import { Command } from "commander";

import getDataLoadCommand from "./data/data-load";
import getDataCleanCommand from "./data/data-clean";
import getDataAppendCommand from "./data/data-append";
import getDataConfigureCommand from "./data/data-configure";

export default function (): Command {
    return new Command("data")
        .description("data commands")
        .addCommand(getDataLoadCommand())
        .addCommand(getDataCleanCommand())
        .addCommand(getDataAppendCommand())
        .addCommand(getDataConfigureCommand());
}
