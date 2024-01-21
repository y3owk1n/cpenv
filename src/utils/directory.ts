import * as fs from "fs";
import * as url from "url";
import * as util from "util";

export type Directory = {
	name: string;
	value: string;
};

export const __dirname = url.fileURLToPath(new URL(".", import.meta.url));
export const readdir = util.promisify(fs.readdir);
export const mkdir = util.promisify(fs.mkdir);

/**
 * Asynchronously checks if the specified path corresponds to a directory.
 *
 * @param sourcePath - The path to the file or directory to check.
 * @returns A promise that resolves to true if the path corresponds to a directory, false otherwise.
 *
 * @throws If there is an issue checking the file or directory status.
 *
 * @example
 * // Usage example:
 * const isDirectory = await isFsDirectory("/path/to/directory");
 * console.log(isDirectory); // true or false
 */
export async function isFsDirectory(sourcePath: string) {
	return (await fs.promises.stat(sourcePath)).isDirectory();
}

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
