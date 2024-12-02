import { copyEnvFilesToVault } from "@/utils/env";
import { confirmCwd } from "../init/prompt";

export async function prepareBackup(): Promise<void> {
	await confirmCwd();
	await copyEnvFilesToVault();
}
