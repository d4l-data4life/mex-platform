import { Command } from "commander";

import getCroissantGenerateCommand from "./croissant/croissant-generate";

export default function (): Command {
    return new Command("croissant").description("Croissant metadata commands").addCommand(getCroissantGenerateCommand());
}
