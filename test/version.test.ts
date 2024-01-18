import { describe, expect, test } from "bun:test";
import {
	getCurrentDescription,
	getCurrentName,
	getCurrentPackageJsonData,
	getCurrentVersion,
} from "@/utils/version";

describe("Package JSON Utility Functions", () => {
	test("getCurrentPackageJsonData should return valid JSON data", () => {
		const packageJsonData = getCurrentPackageJsonData();
		expect(packageJsonData).toBeDefined();
		expect(packageJsonData).toHaveProperty("name");
		expect(packageJsonData).toHaveProperty("version");
		expect(packageJsonData).toHaveProperty("description");
		// Add more specific assertions if needed
	});

	test("getCurrentVersion should return a non-empty string", () => {
		const version = getCurrentVersion();
		expect(version).toBeDefined();
		expect(version).not.toEqual("");
	});

	test("getCurrentName should return a non-empty string", () => {
		const name = getCurrentName();
		expect(name).toBeDefined();
		expect(name).not.toEqual("");
	});

	test("getCurrentDescription should return a non-empty string", () => {
		const description = getCurrentDescription();
		expect(description).toBeDefined();
		expect(description).not.toEqual("");
	});
});
