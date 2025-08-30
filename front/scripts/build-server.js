import esbuild from 'esbuild'
import Config from "config"

let ctx = await esbuild.context({
	entryPoints: [
		'client/index.ts'
	],
	entryNames: '[name]',
	define: {
		ENV: '"' + Config.get("env") + '"',
		BACKEND_URL: '"' + Config.get("backend_url") + '"',
		BOT_USERNAME: '"' + Config.get("bot_username") + '"',
		BOT_VALIDATE_URL: '"' + Config.get("bot_validate_url") + '"',
		BOT_VALIDATE_AUTH_URL: '"' + Config.get("bot_validate_auth_url") + '"',
		HTTPS: '"' + Config.get("https") + '"',
	},
	bundle: true,
	platform: 'node',
	outdir: "build",
	outbase: "client",
	logLevel: "error",
	format: 'cjs',
})

await ctx.watch()
if (Config.get("env") != "development") {
	await ctx.dispose()
}