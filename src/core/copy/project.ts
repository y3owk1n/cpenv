import { getEnvFilesDirectory } from "@/utils/env";
import { Directory, getDirectories } from "@/utils/projects";
import { commanderInit } from "./command";
import { select } from "@inquirer/prompts";

export type Project = Directory;

export async function getProjectsList(): Promise<Directory[]> {
	const envFilesDirectory = await getEnvFilesDirectory();

	const projects = await getDirectories(envFilesDirectory);
	return projects;
}

/**
 * Asynchronously prompts the user to select a project from the available list or uses the project specified through command-line options.
 *
 * @param projects - An array of available project names.
 * @returns A promise that resolves to the selected project name.
 *
 * @throws If there is an issue with the user prompt or if the specified project is not found.
 *
 * @example
 * // Usage example:
 * const projects = ["project1", "project2", ...];
 * const selectedProject = await selectProject(projects);
 * console.log(selectedProject); // 'project1' or 'project2'
 */
export async function selectProject(projects: Directory[]): Promise<string> {
	const options = commanderInit();

	if (!options.project) {
		const project = await select({
			message: "Select a project to copy .env files:",
			choices: projects,
		});
		return project;
	}

	if (!projects.includes(options.project)) {
		console.log("Error: Specified project not found in the directory.");
		process.exit(1); // Exit the process with an error code
	}

	return options.project;
}
