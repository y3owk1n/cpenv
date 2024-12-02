import { Command } from "commander";
import { commanderInit } from "./core/init/command";

process.on("uncaughtException", (error) => {
	if (error instanceof Error && error.name === "ExitPromptError") {
		console.log("👋 until next time!");
	} else {
		// Rethrow unknown errors
		throw error;
	}
});

const program = new Command();

(async () => {
	commanderInit(program);

	await program.parseAsync(process.argv);
})();
