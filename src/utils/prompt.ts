import * as os from "os";
import * as path from "path";
import { confirm, input } from "@inquirer/prompts";

/**
 * Asynchronously prompts the user for the directory path where the environment files should be stored.
 *
 * @returns A Promise that resolves to an object containing the user-inputted vault directory.
 */
export async function promptForVaultDir(): Promise<{
	vaultDir: string;
}> {
	const vaultDir = await input({
		message:
			"Path from the root directory where the .env files should be stored, starts from your home `~/`:",
		default: ".env-files",
	});

	const confirmVaultDir = await confirm({
		message: `Are you sure you want to use ${path.join(
			os.homedir(),
			vaultDir,
		)} as the vault directory?`,
		default: true,
	});

	if (confirmVaultDir === false) {
		return await promptForVaultDir();
	}

	return { vaultDir };
}

/**
 * Function to prompt for global overwrite confirmation.
 *
 * @returns A promise resolving to an object with an 'overwriteAll' boolean property.
 *
 * @example
 * // Usage example:
 * const answer = await promptForGlobalOverwrite();
 * console.log(answer.overwriteAll); // Access the value of the 'overwriteAll' property
 */
export async function promptForGlobalOverwrite(): Promise<{
	overwriteAll: boolean;
}> {
	const overwriteAll = await confirm({
		message:
			"Do you want to overwrite all existing files in the current project if it exists?",
		default: false,
	});

	return { overwriteAll };
}

/**
 * Function to prompt for overwrite confirmation for a specific file.
 *
 * @param file - The name of the file for which overwrite confirmation is needed.
 * @param currentPath - The current path where the file is located.
 * @returns A promise resolving to an object with an 'overwrite' boolean property.
 *
 * @example
 * // Usage example:
 * const answer = await promptForOverwrite("example.txt", "/path/to/directory");
 * console.log(answer.overwrite); // Access the value of the 'overwrite' property
 */
export async function promptForOverwrite(
	file: string,
	destinationPath: string,
): Promise<{ overwrite: boolean }> {
	const overwrite = await confirm({
		message: `Warning: ${file} already exists in ${destinationPath}. Do you want to overwrite it?`,
		default: false,
	});

	return { overwrite };
}
