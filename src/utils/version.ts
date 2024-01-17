import * as path from "path";
import * as fs from "fs";
import { __dirname } from "./file";

/**
 * Retrieves the current version from the package.json file.
 *
 * @returns The version string from the package.json file.
 * @throws If there is an issue reading or parsing the package.json file.
 *
 * @example
 * const version = getCurrentVersion();
 * console.log(`Current version: ${version}`);
 */
export function getCurrentVersion(): string {
  const packageJsonPath: string = path.join(__dirname, "..", "package.json");
  console.log(packageJsonPath);
  const packageJsonContent: string = fs.readFileSync(packageJsonPath, "utf-8");
  const packageJson = JSON.parse(packageJsonContent);
  return packageJson.version;
}
