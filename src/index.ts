#!/usr/bin/env node

import ora from "ora";
import { prepareBackup } from "./core/backup";
import { prepareCopy } from "./core/copy";
import { commandOptions, commanderInit } from "./core/copy/command";
import { promptToSelectAction } from "./core/init/prompt";
import { envInit } from "./utils/env";
import { getCurrentVersion } from "./utils/version";

(async () => {
	try {
		console.log(`CLI Version: ${getCurrentVersion()}`);

		const options = commanderInit(commandOptions);

		await envInit();

		if (options.project || options.autoYes) {
			await prepareCopy(options);
			return;
		}

		const { action } = await promptToSelectAction();

		if (action === "copy") {
			await prepareCopy(options);
			return;
		}

		if (action === "backup") {
			await prepareBackup();
			return;
		}
	} catch (error) {
		if (error instanceof Error) {
			if (error.message === "User force closed the prompt with 0 null") {
				console.log("Exiting...");
				return;
			}
			ora(error.message).fail();
			return;
		}
		console.error("Error:", error);
	}
})();
