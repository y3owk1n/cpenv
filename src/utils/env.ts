import * as os from "os";
import * as path from "path";
import { copyFile, isFsDirectory, mkdir, readdir } from "./file";
import { promptForOverwrite } from "./prompt";

export const envFilesDirectory: string = path.join(os.homedir(), "env-files");

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
