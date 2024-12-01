import { defineConfig } from "rollup";
import commonjs from "@rollup/plugin-commonjs";
import json from "@rollup/plugin-json";
import resolve from "@rollup/plugin-node-resolve";
import terser from "@rollup/plugin-terser";
import typescript from "@rollup/plugin-typescript";

export default defineConfig({
	input: "src/index.ts",
	output: {
		dir: "dist",
		format: "cjs",
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
	],
});
