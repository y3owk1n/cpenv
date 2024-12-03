import {
	getCurrentDescription,
	getCurrentName,
	getCurrentVersion,
} from "@/utils/version";
import type { Command } from "commander";
import { backupCommanderInit } from "./backup/command";
import { copyCommanderInit } from "./copy/command";
import { setupCommanderInit } from "./init/command";

export function commanderInit(program: Command) {
	const name = getCurrentName();
	const description = getCurrentDescription();
	const version = getCurrentVersion();

	program.name(name).description(description).version(version);

	setupCommanderInit(program);

	copyCommanderInit(program);

	backupCommanderInit(program);
}
