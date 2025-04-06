import esbuild from 'esbuild'
import Config from "config"

let ctx = await esbuild.context({
	entryPoints: ['client/index.js'],
	entryNames: '[name]',
	define: {
		ENV: '"' + Config.get("env") + '"',
		BACKEND_URL: '"' + Config.get("backendurl") + '"',
		HTTPS: '"' + Config.get("https") + '"',
	},
	bundle: true,
	platform: 'node',
	outdir: "build",
	outbase: "client",
	format: 'cjs',
})

await ctx.watch()
if (Config.get("env") != "development") {
	await ctx.dispose()
}