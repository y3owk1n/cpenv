import * as fs from "node:fs";
import * as path from "node:path";
import { __dirname } from "./directory";

type PackageJson = {
	name: string;
	version: string;
	description: string;
};

/**
 * Retrieves the data from the package.json file.
 *
 * @returns The current data from the package.json file.
 * @throws If there is an issue reading or parsing the package.json file.
 *
 * @example
 * const jsonData = getCurrentPackageJsonData();
 * console.log(`Current data: ${jsonData}`);
 */
export function getCurrentPackageJsonData(): PackageJson {
	const packageJsonPath: string = path.join(process.cwd(), "package.json");
	const packageJsonContent: string = fs.readFileSync(
		packageJsonPath,
		"utf-8",
	);
	const packageJson = JSON.parse(packageJsonContent);
	return packageJson;
}

export function getCurrentVersion(): string {
	const envVer = process.env.VERSION;
	if (envVer) return envVer;
	const { version } = getCurrentPackageJsonData();
	return version;
}

export function getCurrentName(): string {
	const envName = process.env.NAME;
	if (envName) return envName;
	const { name } = getCurrentPackageJsonData();
	return name;
}

export function getCurrentDescription(): string {
	const envDesc = process.env.DESCRIPTION;
	if (envDesc) return envDesc;
	const { description } = getCurrentPackageJsonData();
	return description;
}
