import commonjs from "@rollup/plugin-commonjs";
import json from "@rollup/plugin-json";
import resolve from "@rollup/plugin-node-resolve";
import replace from "@rollup/plugin-replace";
import terser from "@rollup/plugin-terser";
import typescript from "@rollup/plugin-typescript";
import { defineConfig } from "rollup";
import {
	getCurrentDescription,
	getCurrentName,
	getCurrentVersion,
} from "./src/utils/version";

export default defineConfig({
	input: "./src/index.ts",
	output: {
		dir: "./dist",
		format: "cjs",
		entryFileNames: "[name].cjs", // This ensures the output file has .cjs extension
	},
	plugins: [
		json(),
		commonjs(),
		resolve(), // Helps Rollup find external modules
		typescript({
			tsconfig: "./tsconfig.json", // Ensure this matches your TypeScript configuration
			rootDir: "./src",
			declaration: true,
			outDir: "./dist/types",
		}),
		terser(), // Minification
		replace({
			"process.env.VERSION": () => JSON.stringify(getCurrentVersion()),
			"process.env.NAME": () => JSON.stringify(getCurrentName()),
			"process.env.DESCRIPTION": () =>
				JSON.stringify(getCurrentDescription()),
		}),
	],
});
