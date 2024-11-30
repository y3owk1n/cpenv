import * as fs from "node:fs/promises";
import * as os from "node:os";
import * as path from "node:path";
import ora from "ora";
import { getBackupTimestamp } from "./date";
import { isFsDirectory, mkdir, readdir } from "./directory";
import { checkFileExists, copyFile } from "./file";
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

	const spinner = ora("Checking if env files directory exists");
	spinner.start();

	try {
		await fs.access(envFilesDirectory);
		spinner.succeed(`Env files directory exists at ${envFilesDirectory}`);
	} catch (error) {
		spinner.indent = 2;
		spinner.info(
			`Env files directory not found. Creating a new one at ${envFilesDirectory}`,
		);
		await fs.mkdir(envFilesDirectory, { recursive: true });
		spinner.succeed(`Env files directory created at ${envFilesDirectory}`);
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
	await loadEnvConfig(envConfigDirectory);
}

/**
 * Checks if the environment configuration file exists in the specified directory.
 *
 * @param directory - The path to the directory where the environment configuration file is expected.
 * @returns A Promise that resolves to a boolean indicating whether the configuration file exists.
 *
 * @throws Throws an error if there is an issue accessing the file system.
 *
 * @example
 * // Example usage:
 * const directoryPath = '/path/to/directory';
 * try {
 *   const exists = await envConfigExists(directoryPath);
 *   console.log(`Environment configuration file exists: ${exists}`);
 * } catch (error) {
 *   console.error(`Error checking for environment configuration file: ${error.message}`);
 * }
 */
async function envConfigExists(directory: string): Promise<boolean> {
	try {
		await fs.access(directory);
		return true;
	} catch (error) {
		return false;
	}
}

/**
 * Creates a new environment configuration file in the specified directory with default content.
 *
 * @param directory - The path to the directory where the environment configuration file will be created.
 * @returns A Promise that resolves when the configuration file is successfully created.
 *
 * @throws Throws an error if there is an issue writing to the file system or obtaining user input.
 *
 * @example
 * // Example usage:
 * const directoryPath = '/path/to/directory';
 * try {
 *   await createEnvConfigFile(directoryPath);
 *   console.log('Environment configuration file created successfully.');
 * } catch (error) {
 *   console.error(`Error creating environment configuration file: ${error.message}`);
 * }
 */
async function createEnvConfigFile(directory: string): Promise<void> {
	const { vaultDir } = await promptForVaultDir();

	const defaultConfig: ConfigJson = {
		vaultDir,
	};

	await fs.writeFile(
		directory,
		JSON.stringify(defaultConfig, null, 2),
		"utf-8",
	);
}

/**
 * Loads the environment configuration from the specified directory, creating a new one if it doesn't exist.
 *
 * @param envConfigDirectory - The path to the directory where the environment configuration file is expected.
 * @returns A Promise that resolves to the loaded environment configuration.
 *
 * @throws Throws an error if there is an issue accessing or creating the configuration file,
 *                  or if there is an error parsing the file content.
 *
 * @example
 * // Example usage:
 * const directoryPath = '/path/to/directory';
 * try {
 *   const loadedConfig = await loadEnvConfig(directoryPath);
 *   console.log('Environment configuration loaded:', loadedConfig);
 * } catch (error) {
 *   console.error(`Error loading environment configuration: ${error.message}`);
 * }
 */
async function loadEnvConfig(envConfigDirectory: string): Promise<ConfigJson> {
	const spinner = ora("Checking if config exists").start();
	const configExists = await envConfigExists(envConfigDirectory);

	if (!configExists) {
		spinner.fail(
			`Config file not found. Creating a new one at ${envConfigDirectory}`,
		);

		await createEnvConfigFile(envConfigDirectory);
		spinner.indent = 2;
		spinner.succeed(`Config file created at ${envConfigDirectory}`);
	}
	spinner.start("Loading config...");
	// Load the config (this could be a separate function if needed)
	const configContent = await fs.readFile(envConfigDirectory, "utf-8");
	const config: ConfigJson = JSON.parse(configContent);
	spinner.succeed("Config loaded!");

	return config;
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
export async function copyEnvFilesToProject(
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

		if (await isFsDirectory(sourcePath)) {
			await mkdir(destinationPathWithFile, { recursive: true });
		} else if (file.endsWith(".env")) {
			// Check if file already exists in the destination path
			const fileExists = await checkFileExists(destinationPath, file);

			if (!fileExists) {
				const spinner = ora(
					`Copying ${file} to ${destinationPathWithFile}`,
				).start();
				// No matching .env files found, copy the entire .env file to the project
				await copyFile(sourcePath, destinationPathWithFile);
				spinner.succeed(
					`Successfully copied ${file} to the ${destinationPathWithFile}.`,
				);
			} else {
				// Matching .env files found, proceed with the regular copy process
				await handleEnvFileCopy(
					file,
					sourcePath,
					destinationPathWithFile,
					autoYes,
				);
			}
		}
	}
}

export async function copyEnvFilesToVault(): Promise<void> {
	const envFilesDirectory = await getEnvFilesDirectory();

	const currentDir = process.cwd();

	const currentProjectFolderName = currentDir.split("/").at(-1);

	if (!currentProjectFolderName) {
		console.error("failed to parse the folder name, try again...");
		process.exit(1);
	}

	const currentProjectFolderNameWithTimestamp = `${currentProjectFolderName}-${getBackupTimestamp()}`;

	const filesInProject: (string | Buffer)[] = await readdir(currentDir, {
		recursive: true,
	});

	const destinationPath: string = path.join(
		envFilesDirectory,
		currentProjectFolderNameWithTimestamp,
	);

	await mkdir(destinationPath);

	for (const entry of filesInProject) {
		const file = entry.toString(); // Convert buffer entry to string
		const sourcePath: string = path.join(currentDir, file);

		if (
			file.includes("node_modules") ||
			file.includes(".template") ||
			file.includes(".example")
		) {
			continue;
		}

		const destinationPathWithFile: string = path.join(
			destinationPath,
			file,
		);

		if (file.includes(".env")) {
			const spinner = ora(
				`Copying ${file} to ${destinationPathWithFile}`,
			).start();
			await copyFile(sourcePath, destinationPathWithFile);
			spinner.succeed(
				`Successfully copied ${file} to the ${destinationPathWithFile}.`,
			);
		}
	}
}

/**
 * Asynchronously handles the copying of an .env file, considering overwrite options.
 *
 * @param file - The name of the .env file to copy.
 * @param sourcePath - The source path of the .env file.
 * @param destinationPath - The destination path for the .env file.
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
	autoYes: boolean,
): Promise<void> {
	let autoYesForRemainingFiles = autoYes;

	const shouldOverwrite = autoYesForRemainingFiles;

	if (shouldOverwrite) {
		const spinner = ora(`Copying ${file} to ${destinationPath}`);
		spinner.indent = 2;
		spinner.start();
		await copyFile(sourcePath, destinationPath);
		spinner.succeed(
			`Successfully copied ${file} to the ${destinationPath}.`,
		);
	} else {
		if (!autoYesForRemainingFiles) {
			const overwriteAnswer = await promptForOverwrite(
				file,
				destinationPath,
			);

			const spinner = ora(`Copying ${file} to ${destinationPath}`);
			spinner.indent = 2;
			spinner.start();
			if (overwriteAnswer.overwrite) {
				await copyFile(sourcePath, destinationPath);
				spinner.succeed(
					`Successfully copied ${file} to ${destinationPath}.`,
				);
			} else {
				spinner.info(`Skipped copying ${file} to ${destinationPath}.`);
				autoYesForRemainingFiles = false;
			}
		} else {
			ora(`Skipped copying ${file} to ${destinationPath}.`).info();
		}
	}
}
