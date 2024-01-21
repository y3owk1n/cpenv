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
		console.error("Error:", error instanceof Error ? error.message : error);
	}
})();
