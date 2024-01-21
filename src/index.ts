#!/usr/bin/env node

import { commanderInit } from "./core/copy/command";
import { startCopy } from "./core/copy";
import { Project, getProjectsList, selectProject } from "./core/copy/project";
import { envInit } from "./utils/env";
import { getCurrentVersion } from "./utils/version";

(async () => {
	try {
		console.log(`CLI Version: ${getCurrentVersion()}`);

		await envInit();

		const projects: Project[] = await getProjectsList();
		const selectedProject: string = await selectProject(projects);

		const options = commanderInit();

		await startCopy(selectedProject, options);
	} catch (error) {
		console.error("Error:", error instanceof Error ? error.message : error);
	}
})();
