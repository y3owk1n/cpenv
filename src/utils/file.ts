import * as fs from "fs";
import * as url from "url";
import * as util from "util";

export const __dirname = url.fileURLToPath(new URL(".", import.meta.url));
export const readdir = util.promisify(fs.readdir);
export const mkdir = util.promisify(fs.mkdir);
export const copyFile = util.promisify(fs.copyFile);

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
