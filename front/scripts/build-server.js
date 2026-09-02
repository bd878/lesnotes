import esbuild from 'esbuild'
import Config from "config"

let ctx = await esbuild.context({
	entryPoints: [
		'client/index.ts',
	],
	entryNames: '[name]',
	define: {
		ENV: '"' + Config.get("env") + '"',
		DOMAIN: '"' + Config.get("domain") + '"',
		PUBLIC_USER_ID: Config.get("public_user_id"),
		BACKEND_URL: '"' + Config.get("backend_url") + '"',
		HTTPS: '"' + Config.get("https") + '"',
		LIMIT: "24"
	},
	external: ["pino", "koa-pino-logger", "pino-opentelemetry-transport", "thread-stream"],
	bundle: true,
	minify: false,
	platform: 'node',
	outdir: "build",
	outbase: "client",
	logLevel: "error",
	format: 'cjs',
})

await ctx.watch()
if (Config.get("env") != "development") {
	console.log("dispose")
	await ctx.dispose()
}