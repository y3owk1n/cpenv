import {
	getCurrentDescription,
	getCurrentName,
	getCurrentVersion,
} from "@/utils/version";
import type { Command } from "commander";
import { prepareBackup } from "../backup";
import { backupCommanderInit } from "../backup/command";
import { prepareCopy } from "../copy";
import { copyCommanderInit } from "../copy/command";
import { promptToSelectAction } from "./prompt";

export function commanderInit(program: Command) {
	const name = getCurrentName();
	const description = getCurrentDescription();
	const version = getCurrentVersion();

	program
		.name(name)
		.description(description)
		.version(version)
		.action(async () => {
			const { action } = await promptToSelectAction();

			if (action === "copy") {
				await prepareCopy(program.opts);
				return;
			}

			if (action === "backup") {
				await prepareBackup();
				return;
			}
		});

	copyCommanderInit(program);

	backupCommanderInit(program);
}
