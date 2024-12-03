import { Command } from "@commander-js/extra-typings";
import { commanderInit } from "./core";

process.on("uncaughtException", (error) => {
	if (error instanceof Error && error.name === "ExitPromptError") {
		console.log("👋 until next time!");
	} else {
		// Rethrow unknown errors
		throw error;
	}
});

function handleSigTerm() {
	process.stdout.write("\x1B[?25h");
	process.stdout.write("\n");
	process.exit(0);
}

process.on("SIGINT", handleSigTerm);
process.on("SIGTERM", handleSigTerm);
process.on("exit", handleSigTerm);

const program = new Command();

(async () => {
	commanderInit(program);

	await program.parseAsync(process.argv);
})();
