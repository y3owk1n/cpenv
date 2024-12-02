import { type Directory, getDirectories } from "@/utils/directory";
import { getEnvFilesDirectory } from "@/utils/env";
import ora from "ora";
import type { CopyCommandOptionValues } from "./command";
import { promptToSelectProject } from "./prompt";

export type Project = Directory;

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

export async function selectProject(
	projects: Directory[],
	options: CopyCommandOptionValues,
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
