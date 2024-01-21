#!/usr/bin/env node

import { startCopy } from "./core/copy";
import { commandOptions, commanderInit } from "./core/copy/command";
import { Project, getProjectsList, selectProject } from "./core/copy/project";
import { envInit } from "./utils/env";
import { getCurrentVersion } from "./utils/version";

(async () => {
	try {
		console.log(`CLI Version: ${getCurrentVersion()}`);

		const options = commanderInit(commandOptions);

		await envInit();

		const projects: Project[] = await getProjectsList();
		const selectedProject: string = await selectProject(projects, options);

		await startCopy(selectedProject, options);
	} catch (error) {
		if (error instanceof Error) {
			if (error.message === "User force closed the prompt with 0 null") {
				console.log("Exiting...");
				return;
			}
			console.error("Error:", error.message);
			return;
		}
		console.error("Error:", error);
	}
})();
