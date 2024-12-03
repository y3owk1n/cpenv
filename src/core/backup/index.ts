import { copyEnvFilesToVault } from "@/utils/env";
import { confirmCwdPrompt } from "./prompt";

export async function prepareBackup(): Promise<void> {
	await confirmCwdPrompt();
	await copyEnvFilesToVault();
}
