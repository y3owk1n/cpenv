import { Command } from "commander";

export type CommandOption = {
	flags: string;
	description: string;
};

/**
 * Configures a Commander program with specified command options.
 *
 * @param program - The Commander program to configure.
 * @param options - An array of command options to be added to the program.
 *
 * @example
 * // Usage example:
 * const program = new commander.Command();
 * configureCommand(program, commonOptions);
 * program.parse(process.argv);
 */
export function configureCommand(
	program: Command,
	options: CommandOption[],
): void {
	for (const opt of options) {
		program.option(opt.flags, opt.description);
	}
	program.parse(process.argv);
}
