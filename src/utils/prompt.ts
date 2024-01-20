import inquirer from "inquirer";

export async function promptForVaultDir(): Promise<{
	vaultDir: string;
}> {
	const response = await inquirer.prompt([
		{
			type: "input",
			name: "vaultDir",
			message:
				"Path to the root directory where the .env files should be stored, starts from your home `~/`:",
			default: ".env-files",
		},
	]);

	return { vaultDir: response.vaultDir };
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
	const response = await inquirer.prompt([
		{
			type: "confirm",
			name: "overwriteAll",
			message:
				"Do you want to overwrite all existing files in the current project if it exists?",
			default: false,
		},
	]);

	return { overwriteAll: response.overwriteAll };
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
	currentPath: string,
): Promise<{ overwrite: boolean }> {
	const response = await inquirer.prompt([
		{
			type: "confirm",
			name: "overwrite",
			message: `Warning: ${file} already exists in ${currentPath}. Do you want to overwrite it?`,
			default: false,
		},
	]);

	return { overwrite: response.overwrite };
}
