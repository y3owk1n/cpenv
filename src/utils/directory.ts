import * as fs from "node:fs";
import * as url from "node:url";
import * as util from "node:util";

export type Directory = {
	name: string;
	value: string;
};

export const __dirname = url.fileURLToPath(new URL(".", import.meta.url));
export const readdir = util.promisify(fs.readdir);
export const mkdir = util.promisify(fs.mkdir);

export async function isFsDirectory(sourcePath: string) {
	return (await fs.promises.stat(sourcePath)).isDirectory();
}

export async function getDirectories(directory: string): Promise<Directory[]> {
	const dirents = await readdir(directory, { withFileTypes: true });
	return dirents
		.filter((dirent) => dirent.isDirectory())
		.map((dirent) => ({
			name: dirent.name,
			value: dirent.name,
		}));
}
