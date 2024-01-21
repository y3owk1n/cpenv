import { copyEnvFiles } from "@/utils/env";
import { promptForGlobalOverwrite } from "@/utils/prompt";
import { OptionValues } from "commander";

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
export async function startCpCli(
	selectedProject: string,
	options: OptionValues,
): Promise<void> {
	if (!options.autoYes) {
		const overwriteAllAnswer = await promptForGlobalOverwrite();

		await copyEnvFiles(
			selectedProject,
			"",
			overwriteAllAnswer.overwriteAll,
		);
		console.log("Copy completed successfully!");
		return;
	}

	await copyEnvFiles(selectedProject, "", options.autoYes);
	console.log("Copy completed successfully!");
}
