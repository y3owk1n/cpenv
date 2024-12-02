import * as fs from "node:fs";
import * as fsPromise from "node:fs/promises";
import * as path from "node:path";
import * as util from "node:util";

export const copyFile = util.promisify(fs.copyFile);

export async function checkFileExists(
	destinationPath: string,
	fileName: string,
) {
	const filePath = path.join(destinationPath, fileName);

	try {
		await fsPromise.access(filePath, fs.constants.F_OK);
		return true; // File exists
	} catch (error) {
		return false; // File does not exist
	}
}
