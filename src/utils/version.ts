import * as fs from "fs";
import * as path from "path";

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
	const rootPath: string = process.cwd();
	const packageJsonPath: string = path.join(rootPath, "package.json");
	const packageJsonContent: string = fs.readFileSync(
		packageJsonPath,
		"utf-8",
	);
	const packageJson = JSON.parse(packageJsonContent);
	return packageJson;
}

export function getCurrentVersion(): string {
	const { version } = getCurrentPackageJsonData();
	return version;
}

export function getCurrentName(): string {
	const { name } = getCurrentPackageJsonData();
	return name;
}

export function getCurrentDescription(): string {
	const { description } = getCurrentPackageJsonData();
	return description;
}
