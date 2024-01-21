#!/usr/bin/env node

import { startCpCli } from "./core/copy";
import { commanderInit } from "./core/copy/command";
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

		await startCpCli(selectedProject, options);
	} catch (error) {
		console.error("Error:", error instanceof Error ? error.message : error);
	}
})();
