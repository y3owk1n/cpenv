#!/usr/bin/env node

import ora from "ora";
import { startCopy } from "./core/copy";
import { commandOptions, commanderInit } from "./core/copy/command";
import {
	type Project,
	getProjectsList,
	selectProject,
} from "./core/copy/project";
import { confirmCwd, promptToSelectAction } from "./core/init/prompt";
import { copyEnvFilesToVault, envInit } from "./utils/env";
import { getCurrentVersion } from "./utils/version";

(async () => {
	try {
		console.log(`CLI Version: ${getCurrentVersion()}`);

		const options = commanderInit(commandOptions);

		await envInit();

		const { action } = await promptToSelectAction();

		if (action === "copy") {
			const projects: Project[] = await getProjectsList();
			const selectedProject: string = await selectProject(
				projects,
				options,
			);

			await startCopy(selectedProject, options);
		}

		if (action === "backup") {
			await confirmCwd();
			await copyEnvFilesToVault();
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
