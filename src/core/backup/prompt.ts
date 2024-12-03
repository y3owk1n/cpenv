import { confirm } from "@inquirer/prompts";

export async function confirmCwdPrompt(): Promise<void> {
	const currentDir = process.cwd();

	const answer = await confirm({
		message: `Is this your root directory to perform the backup? (${currentDir})`,
	});

	if (!answer) {
		console.log("cd to your desired directory and restart the backup");
		process.exit(1);
	}

	return;
}
