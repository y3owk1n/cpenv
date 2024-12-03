import { envInit } from "@/utils/env";
import type { Command } from "commander";

export function setupCommanderInit(program: Command): void {
	program
		.command("setup")
		.description("initialize configurations for vault")
		.action(async () => {
			await envInit();
		});
}
