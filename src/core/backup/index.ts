import { copyEnvFilesToVault } from "@/utils/env";
import { confirmCwd } from "../init/prompt";

export async function prepareBackup() {
	await confirmCwd();
	await copyEnvFilesToVault();
}
