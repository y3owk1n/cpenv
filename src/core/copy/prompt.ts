import type { Directory } from "@/utils/directory";
import { search } from "@inquirer/prompts";

export async function promptToSelectProject(projects: Directory[]): Promise<{
	project: string;
}> {
	const project = await search({
		message: "Select a project to copy .env files",
		source: async (input) => {
			if (!input) {
				return projects;
			}

			return projects.filter((project) =>
				project.value.toLocaleLowerCase().includes(input.toLowerCase()),
			);
		},
	});

	return { project };
}
