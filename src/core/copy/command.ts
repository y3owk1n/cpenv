import { type CommandOption, configureCommand } from "@/utils/command";
import {
	getCurrentDescription,
	getCurrentName,
	getCurrentVersion,
} from "@/utils/version";
import type { OptionValues as CommanderOptionValues } from "commander";
import { Command } from "commander";

export interface OptionValues extends CommanderOptionValues {
	project?: string;
	autoYes?: boolean;
}

export const commandOptions: CommandOption[] = [
	{
		flags: "-p, --project <project>",
		description: "Select a project to copy .env files",
	},
	{
		flags: "-y, --auto-yes",
		description: "Automatically overwrite files without prompting",
	},
];

/**
 * Initializes and configures a Commander program with specified command options.
 *
 * @param commandOptions - An array of command options to be added to the program.
 * @returns An object containing the parsed command-line options.
 *
 * @example
 * // Usage example:
 * const options = commanderInit();
 * console.log(options.project); // Access the value of the 'project' option
 */
export function commanderInit(
	commandOptions: CommandOption[] = [],
): OptionValues {
	const program = new Command();

	const name = getCurrentName();
	const description = getCurrentDescription();
	const version = getCurrentVersion();

	program.name(name).description(description).version(version);
	configureCommand(program, commandOptions);

	const options = program.opts();

	return options;
}
