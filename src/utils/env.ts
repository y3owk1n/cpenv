import * as os from "os";
import * as path from "path";
import * as fs from "fs/promises";
import { isFsDirectory, mkdir, readdir } from "./directory";
import { copyFile } from "./file";
import { promptForOverwrite, promptForVaultDir } from "./prompt";

type ConfigJson = {
	vaultDir: string;
};

export const envConfigDirectory: string = path.join(
	os.homedir(),
	".env-files.json",
);

/**
 * Asynchronously creates the environment files directory if not found.
 *
 * @param vaultDir - The directory within the user's home directory where the environment files directory will be created.
 * @returns A Promise that resolves to the path of the environment files directory.
 * @throws If there is an issue accessing or creating the directory.
 */
export async function createEnvFilesDirectoryIfNotFound(
	vaultDir: string,
): Promise<string> {
	const envFilesDirectory = path.join(os.homedir(), vaultDir);

	try {
		await fs.access(envFilesDirectory);
	} catch (error) {
		console.log(
			`Env files directory not found. Creating a new one at ${envFilesDirectory}`,
		);
		await fs.mkdir(envFilesDirectory, { recursive: true });
		console.log("Env files directory created.");
	}

	return envFilesDirectory;
}

/**
 * Asynchronously retrieves the path to the environment files directory.
 *
 * @returns A Promise that resolves to the path of the environment files directory.
 * @throws If there is an issue obtaining the environment configuration data or creating the directory.
 */
export async function getEnvFilesDirectory(): Promise<string> {
	const { vaultDir } = await getEnvConfigJsonData();

	const envFilesDirectory = await createEnvFilesDirectoryIfNotFound(vaultDir);

	return envFilesDirectory;
}

/**
 * Asynchronously initializes the environment by checking and creating the configuration file if not found.
 *
 * @returns A Promise that resolves when the initialization is complete.
 * @throws If there is an issue accessing, creating, or loading the configuration file.
 */
export async function envInit(): Promise<void> {
	try {
		console.log("Checking config file...");

		try {
			// Use try-catch block to handle file not found error
			await fs.access(envConfigDirectory);
		} catch (error) {
			// If the file doesn't exist, create a new one with default content
			console.log(
				`Config file not found. Creating a new one at ${envConfigDirectory}`,
			);

			const { vaultDir } = await promptForVaultDir();

			const defaultConfig: ConfigJson = {
				vaultDir,
			};

			await fs.writeFile(
				envConfigDirectory,
				JSON.stringify(defaultConfig, null, 2),
				"utf-8",
			);
			console.log("Config file created.");
		}
	} catch (error) {
		console.error(
			"Error while loading config:",
			error instanceof Error ? error.message : error,
		);
		throw error;
	}
}

/**
 * Asynchronously reads and parses the environment configuration file content.
 *
 * @returns A Promise that resolves to the parsed content of the environment configuration file.
 * @throws If there is an issue reading or parsing the configuration file.
 */
export async function getEnvConfigJsonData(): Promise<ConfigJson> {
	try {
		const envConfigJsonContent: string = await fs.readFile(
			envConfigDirectory,
			"utf-8",
		);
		const envConfigJson = JSON.parse(envConfigJsonContent);
		return envConfigJson;
	} catch (error) {
		console.error(
			"Error while loading config:",
			error instanceof Error ? error.message : error,
		);
		throw error;
	}
}

/**
 * Asynchronously copies .env files from the specified project to the current working directory.
 *
 * @param project - The name of the project containing .env files.
 * @param [currentPath=""] - The current path within the project (used for recursive copying).
 * @param [autoYes=false] - Whether to automatically overwrite files without prompting.
 * @returns A promise that resolves once the copying process is complete.
 *
 * @throws If there is an issue reading or copying files.
 *
 * @example
 * // Usage example:
 * await copyEnvFiles("myProject");
 */
export async function copyEnvFiles(
	project: string,
	currentPath = "",
	autoYes = false,
): Promise<void> {
	const envFilesDirectory = await getEnvFilesDirectory();
	const projectPath: string = path.join(
		envFilesDirectory,
		project,
		currentPath,
	);

	const filesInProject: (string | Buffer)[] = await readdir(projectPath, {
		recursive: true,
	});

	for (const entry of filesInProject) {
		const file = entry.toString(); // Convert buffer entry to string
		const sourcePath: string = path.join(projectPath, file);
		const destinationPath: string = path.join(process.cwd(), currentPath);
		const destinationPathWithFile: string = path.join(
			process.cwd(),
			currentPath,
			file,
		);

		const filesInDestinationPath: (string | Buffer)[] =
			await readdir(destinationPath);

		if (await isFsDirectory(sourcePath)) {
			await mkdir(destinationPathWithFile, { recursive: true });
			await copyEnvFiles(project, path.join(currentPath, file), autoYes);
		} else if (file.endsWith(".env")) {
			// Check if there are any matching .env files in the current project folder
			const matchingEnvFiles = filesInDestinationPath.filter((f) =>
				f.toString().endsWith(".env"),
			);

			if (matchingEnvFiles.length === 0) {
				// No matching .env files found, copy the entire .env file to the project
				await copyFile(sourcePath, destinationPathWithFile);
			} else {
				// Matching .env files found, proceed with the regular copy process
				await handleEnvFileCopy(
					file,
					sourcePath,
					destinationPathWithFile,
					currentPath,
					autoYes,
				);
			}
		}
	}
}

/**
 * Asynchronously handles the copying of an .env file, considering overwrite options.
 *
 * @param file - The name of the .env file to copy.
 * @param sourcePath - The source path of the .env file.
 * @param destinationPath - The destination path for the .env file.
 * @param currentPath - The current path within the project.
 * @param autoYes - Whether to automatically overwrite files without prompting.
 * @returns A promise that resolves once the copying process is complete.
 *
 * @throws If there is an issue reading or copying files.
 *
 * @example
 * // Usage example:
 * await handleEnvFileCopy("example.env", "/path/to/source", "/path/to/destination", "/current/path", false);
 */
export async function handleEnvFileCopy(
	file: string,
	sourcePath: string,
	destinationPath: string,
	currentPath: string,
	autoYes: boolean,
): Promise<void> {
	let autoYesForRemainingFiles = autoYes;

	const shouldOverwrite = autoYesForRemainingFiles;

	if (shouldOverwrite) {
		await copyFile(sourcePath, destinationPath);
		console.log(`Successfully copied ${file} to the ${currentPath}.`);
	} else {
		if (!autoYesForRemainingFiles) {
			const overwriteAnswer = await promptForOverwrite(file, currentPath);

			if (overwriteAnswer.overwrite) {
				await copyFile(sourcePath, destinationPath);
				console.log(`Successfully copied ${file} to ${currentPath}.`);
			} else {
				console.log(`Skipped copying ${file} to ${currentPath}.`);
				autoYesForRemainingFiles = false;
			}
		} else {
			console.log(`Skipped copying ${file} to ${currentPath}.`);
		}
	}
}
