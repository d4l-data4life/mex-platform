import { Command } from "commander";

import getAuthLoginCommand from "./auth/auth-token";

export default function (): Command {
    return new Command("auth").description("authentication commands").addCommand(getAuthLoginCommand());
}
