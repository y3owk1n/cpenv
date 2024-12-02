import { prepareBackup } from "./core/backup";
import { prepareCopy } from "./core/copy";
import { commandOptions, commanderInit } from "./core/copy/command";
import { promptToSelectAction } from "./core/init/prompt";
import { envInit } from "./utils/env";
import { getCurrentVersion } from "./utils/version";

process.on("uncaughtException", (error) => {
	if (error instanceof Error && error.name === "ExitPromptError") {
		console.log("👋 until next time!");
	} else {
		// Rethrow unknown errors
		throw error;
	}
});

(async () => {
	console.log(`CLI Version: ${getCurrentVersion()}`);

	const options = commanderInit(commandOptions);

	await envInit();

	if (options.project || options.autoYes) {
		await prepareCopy(options);
		return;
	}

	const { action } = await promptToSelectAction();

	if (action === "copy") {
		await prepareCopy(options);
		return;
	}

	if (action === "backup") {
		await prepareBackup();
		return;
	}
})();
