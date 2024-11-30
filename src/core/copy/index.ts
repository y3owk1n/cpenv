import { copyEnvFilesToProject } from "@/utils/env";
import { promptForGlobalOverwrite } from "@/utils/prompt";
import ora from "ora";
import type { OptionValues } from "./command";
import { type Project, getProjectsList, selectProject } from "./project";

export async function prepareCopy(options: OptionValues): Promise<void> {
	const projects: Project[] = await getProjectsList();
	const selectedProject: string = await selectProject(projects, options);

	await startCopy(selectedProject, options);
}

/**
 * Asynchronously copies .env files from the specified project to the current working directory.
 *
 * @param selectedProject - The name of the project containing .env files.
 * @param options - An object containing command-line options.
 * @returns A promise that resolves once the copying process is complete.
 *
 * @throws If there is an issue reading or copying files.
 *
 * @example
 * // Usage example:
 * await copyEnvFiles("myProject");
 */
export async function startCopy(
	selectedProject: string,
	options: OptionValues,
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
