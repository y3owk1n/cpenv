import * as commander from "commander";

type CommandOption = {
	flags: string;
	description: string;
};

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
 * @returns An object containing the parsed command-line options.
 *
 * @example
 * // Usage example:
 * const options = commanderInit();
 * console.log(options.project); // Access the value of the 'project' option
 */
export function commanderInit() {
	const program = new commander.Command();
	configureCommand(program, commandOptions);

	const options = program.opts();

	return options;
}

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
	program: commander.Command,
	options: CommandOption[],
): void {
	for (const opt of options) {
		program.option(opt.flags, opt.description);
	}
	program.parse(process.argv);
}
