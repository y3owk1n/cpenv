import { confirm, select } from "@inquirer/prompts";

export async function promptToSelectAction(): Promise<{
	action: string;
}> {
	const action = await select({
		message: "Select an action",
		choices: [
			{
				name: "copy",
				value: "copy",
				description: "copy .env files to project",
			},
			{
				name: "backup",
				value: "backup",
				description: "backup .env files to your vault",
			},
		],
	});

	return { action };
}

export async function confirmCwd(): Promise<void> {
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
