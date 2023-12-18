import { Command } from "commander";

import getConfigUpdateCommand from "./config/config-update";
import getConfigSendCommand from "./config/config-send";

export default function (): Command {
    return new Command("config").description("config commands").addCommand(getConfigUpdateCommand()).addCommand(getConfigSendCommand());
}
