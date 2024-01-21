import { select } from "@inquirer/prompts";
import { getEnvFilesDirectory } from "./env";
import { readdir } from "./file";
import { commanderInit } from "@/core/copy/command";

export type Directory = {
	name: string;
	value: string;
};

/**
 * Asynchronously retrieves a list of directories within the envFilesDirectory.
 *
 * @param directory - The directory within the user's home directory where the environment files directory is located.
 * @returns A promise that resolves to an array of directory names.
 *
 * @throws If there is an issue reading the directory contents.
 *
 * @example
 * // Usage example:
 * const directories = await getDirectories("/path/to/directory");
 * console.log(directories);
 */
export async function getDirectories(directory: string): Promise<Directory[]> {
	const dirents = await readdir(directory, { withFileTypes: true });
	return dirents
		.filter((dirent) => dirent.isDirectory())
		.map((dirent) => ({
			name: dirent.name,
			value: dirent.name,
		}));
}
