import { envInit } from "@/utils/env";
import type { Command } from "commander";
import { prepareBackup } from ".";

export function backupCommanderInit(program: Command): void {
	program
		.command("backup")
		.description("backup current project env(s) to vault")
		.action(async () => {
			await envInit();
			await prepareBackup();
		});
}
