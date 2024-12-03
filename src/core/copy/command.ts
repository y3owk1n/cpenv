import { envInit } from "@/utils/env";
import type { Command } from "@commander-js/extra-typings";
import type { OptionValues as CommanderOptionValues } from "commander";
import { prepareCopy } from ".";

export interface CopyCommandOptionValues extends CommanderOptionValues {
	project?: string;
	autoYes?: boolean;
}

export function copyCommanderInit(program: Command): void {
	program
		.command("copy")
		.description("copy env(s) from vault to project")
		.option(
			"-p, --project <project>",
			"Select a project to copy .env files",
		)
		.option(
			"-y, --auto-yes",
			"Automatically overwrite files without prompting",
		)
		.action(async (options) => {
			await envInit();
			await prepareCopy(options);
		});
}
