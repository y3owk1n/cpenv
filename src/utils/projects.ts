import { rawlist } from "@inquirer/prompts";
import { commanderInit } from "./commands";
import { getEnvFilesDirectory } from "./env";
import { readdir } from "./file";

/**
 * Asynchronously retrieves a list of directories within the envFilesDirectory.
 *
 * @returns A promise that resolves to an array of directory names.
 *
 * @throws If there is an issue reading the directory contents.
 *
 * @example
 * // Usage example:
 * const projects = await getProjectsList();
 * console.log(projects); // ['project1', 'project2', ...]
 */
export async function getProjectsList(): Promise<
	{ name: string; value: string }[]
> {
	const envFilesDirectory = await getEnvFilesDirectory();
	const dirents = await readdir(envFilesDirectory, { withFileTypes: true });
	return dirents
		.filter((dirent) => dirent.isDirectory())
		.map((dirent) => ({
			name: dirent.name,
			value: dirent.name,
		}));
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
export async function selectProject(
	projects: { name: string; value: string }[],
): Promise<string> {
	const options = commanderInit();

	if (!options.project) {
		const project = await rawlist({
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
