#!/usr/bin/env node

import { commanderInit } from "./utils/commands";
import { copyEnvFiles, getEnvConfigJsonData } from "./utils/env";
import { getProjectsList, selectProject } from "./utils/projects";
import { promptForGlobalOverwrite } from "./utils/prompt";
import { getCurrentVersion } from "./utils/version";

(async () => {
	try {
		console.log(`CLI Version: ${getCurrentVersion()}`);

		const config = await getEnvConfigJsonData();

		const projects: string[] = await getProjectsList();
		const selectedProject: string = await selectProject(projects);

		const options = commanderInit();

		if (!options.autoYes) {
			const overwriteAllAnswer = await promptForGlobalOverwrite();

			await copyEnvFiles(
				selectedProject,
				"",
				overwriteAllAnswer.overwriteAll,
			);
			console.log("Copy completed successfully!");
			return;
		}

		await copyEnvFiles(selectedProject, "", options.autoYes);
		console.log("Copy completed successfully!");
	} catch (error) {
		console.error("Error:", error instanceof Error ? error.message : error);
	}
})();
