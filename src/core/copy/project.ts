import { Directory, getDirectories } from "@/utils/directory";
import { getEnvFilesDirectory } from "@/utils/env";
import ora from "ora";
import { OptionValues } from "./command";
import { promptToSelectProject } from "./prompt";

export type Project = Directory;

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
 * console.log(projects);
 * // [{ name: "project1", value: "project1" }, { name: "project2", value: "project2" }, ...]
 */
export async function getProjectsList(): Promise<Directory[]> {
	const envFilesDirectory = await getEnvFilesDirectory();

	const spinner = ora("Loading available projects");
	const projects = await getDirectories(envFilesDirectory);

	if (projects.length === 0) {
		spinner.warn(
			`No projects found in ${envFilesDirectory}, add a project first.`,
		);
		spinner.warn(
			"If you indeed have projects but the directory is wrong, reconfigure the vault path at ~/.env-files.json",
		);
		process.exit(1);
	}

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
export async function selectProject(
	projects: Directory[],
	options: OptionValues,
): Promise<string> {
	if (!options.project) {
		const { project } = await promptToSelectProject(projects);
		return project;
	}

	const projectsInStrArr = projects.map((project) => project.value);

	if (projectsInStrArr.length === 0) {
		ora("No projects found in the vault, add a project first.").fail();
		ora(
			"If you indeed have projects but the directory is wrong, reconfigure the vault path at ~/.env-files.json",
		).fail();
	}

	if (!projectsInStrArr.includes(options.project)) {
		ora("Specified project not found in the directory.").fail();
		ora(`Available options: ${projectsInStrArr.join(", ")}.`).fail();
		process.exit(1); // Exit the process with an error code
	}

	return options.project;
}
