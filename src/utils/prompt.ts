import * as os from "node:os";
import * as path from "node:path";
import { confirm, input } from "@inquirer/prompts";

export async function promptForVaultDir(): Promise<{
	vaultDir: string;
}> {
	const vaultDir = await input({
		message:
			"Path from the root directory where the .env files should be stored, starts from your home `~/`:",
		default: ".env-files",
	});

	const confirmVaultDir = await confirm({
		message: `Are you sure you want to use ${path.join(
			os.homedir(),
			vaultDir,
		)} as the vault directory?`,
		default: true,
	});

	if (confirmVaultDir === false) {
		return await promptForVaultDir();
	}

	return { vaultDir };
}

export async function promptForGlobalOverwrite(): Promise<{
	overwriteAll: boolean;
}> {
	const overwriteAll = await confirm({
		message:
			"Do you want to overwrite all existing .env(s) in the current project if it exists?",
		default: false,
	});

	return { overwriteAll };
}

export async function promptForOverwrite(
	file: string,
	destinationPath: string,
): Promise<{ overwrite: boolean }> {
	const overwrite = await confirm({
		message: `Warning: ${file} already exists in ${destinationPath}. Do you want to overwrite it?`,
		default: false,
	});

	return { overwrite };
}
