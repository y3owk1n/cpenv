#!/usr/bin/env node

import { startCpCli } from "./core/copy";
import { commanderInit } from "./utils/commands";
import { envInit } from "./utils/env";
import { Project, getProjectsList, selectProject } from "./utils/projects";
import { getCurrentVersion } from "./utils/version";

(async () => {
	try {
		console.log(`CLI Version: ${getCurrentVersion()}`);

		await envInit();

		const projects: Project[] = await getProjectsList();
		const selectedProject: string = await selectProject(projects);

		const options = commanderInit();

		await startCpCli(selectedProject, options);
	} catch (error) {
		console.error("Error:", error instanceof Error ? error.message : error);
	}
})();
