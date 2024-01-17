import { expect, spyOn, test } from "bun:test";

import * as commander from "commander";
import {
	commandOptions,
	commanderInit,
	configureCommand,
} from "../src/utils/commands";

test("should configure Commander program with specified command options", () => {
	const program = new commander.Command();
	const spyOption = spyOn(program, "option");

	configureCommand(program, commandOptions);

	// Check if spyOption was called at least once
	expect(spyOption).toHaveBeenCalled();

	// Check individual calls if any
	spyOption.mock.calls.forEach((call, index) => {
		expect(call[0]).toBe(commandOptions[index].flags);
		expect(call[1]).toBe(commandOptions[index].description);
	});
});

test("should parse command-line options and return the parsed options object", () => {
	// Mock process.argv to simulate command-line arguments
	const originalArgv = process.argv;
	process.argv = ["node", "script-name", "-p", "example-project", "-y"];

	const options = commanderInit();

	// Reset process.argv to its original value
	process.argv = originalArgv;

	// Verify that the parsed options object contains the expected values
	expect(options.project).toBe("example-project");
	expect(options.autoYes).toBe(true);
});
