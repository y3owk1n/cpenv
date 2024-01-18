await Bun.build({
	entrypoints: ["./src/index.ts"],
	format: "esm",
	outdir: "./dist",
	minify: true,
	target: "node",
	splitting: true,
	external: ["*"],
	// sourcemap: "external",
});
