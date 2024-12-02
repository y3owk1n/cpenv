import { copyEnvFilesToProject } from "@/utils/env";
import { promptForGlobalOverwrite } from "@/utils/prompt";
import ora from "ora";
import type { CopyCommandOptionValues } from "./command";
import { type Project, getProjectsList, selectProject } from "./project";

export async function prepareCopy(
	options: CopyCommandOptionValues,
): Promise<void> {
	const projects: Project[] = await getProjectsList();
	const selectedProject: string = await selectProject(projects, options);

	await startCopy(selectedProject, options);
}

export async function startCopy(
	selectedProject: string,
	options: CopyCommandOptionValues,
): Promise<void> {
	if (!options.autoYes) {
		const overwriteAllAnswer = await promptForGlobalOverwrite();

		await copyEnvFilesToProject(
			selectedProject,
			"",
			overwriteAllAnswer.overwriteAll,
		);
		ora("Copy completed successfully!").succeed();
		return;
	}

	await copyEnvFilesToProject(selectedProject, "", options.autoYes);
	ora("Copy completed successfully!").succeed();
}
