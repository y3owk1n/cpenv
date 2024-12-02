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

export async function getEnvFilesDirectory(): Promise<string> {
	const { vaultDir } = await getEnvConfigJsonData();

	const envFilesDirectory = await createEnvFilesDirectoryIfNotFound(vaultDir);

	return envFilesDirectory;
}

export async function envInit(): Promise<void> {
	await loadEnvConfig(envConfigDirectory);
}

async function envConfigExists(directory: string): Promise<boolean> {
	try {
		await fs.access(directory);
		return true;
	} catch (error) {
		return false;
	}
}

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
			destinationPath,
			file,
		);

		if (await isFsDirectory(sourcePath)) {
			await mkdir(destinationPathWithFile, { recursive: true });
		} else if (file.includes(".env")) {
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
