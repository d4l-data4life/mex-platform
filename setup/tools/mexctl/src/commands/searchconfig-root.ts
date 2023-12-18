import { Command } from "commander";

import getSearchConfigListCommand from "./searchconfig/searchconfig-list";
import getSearchConfigAppendCommand from "./searchconfig/searchconfig-append";

export default function (): Command {
    return new Command("searchconfig")
        .description("searchconfig commands")
        .addCommand(getSearchConfigListCommand())
        .addCommand(getSearchConfigAppendCommand());
}
