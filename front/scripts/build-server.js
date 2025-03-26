import esbuild from 'esbuild'
import Config from "config"

await esbuild.build({
	entryPoints: ['client/index.js'],
	entryNames: '[name]',
	define: {
		ENV: '"' + Config.get("env") + '"',
	},
	bundle: true,
	platform: 'node',
	outdir: "build",
	outbase: "client",
	format: 'cjs',
})
