import {unlink, readdir, stat, rm} from "node:fs/promises";
// TODO: import ts from "typescript";
import Config from "config";
import path from "node:path";
import esbuild from 'esbuild';
import postcss from 'esbuild-postcss';

/*remove stale files except favicon.ico*/
const files = await readdir('public', {withFileTypes: true});
for (const file of files) {
	const st = await stat('public/' + file.name)
	if (st.isDirectory())
		await rm('public/' + file.name, { recursive: true, force: true })

	const ext = path.extname(file.name)
	if (ext == ".js" || ext == ".css" || ext == ".map")
		await unlink('public/' + file.name)
}

let ctx = await esbuild.context({
	entryPoints: [
		'client/gui/pages/**/*.tsx',
		'client/styles/styles.css'
	],
	entryNames: '[name]',
	define: {
		BACKEND_URL: '"' + Config.get("backend_url") + '"',
		BOT_VALIDATE_URL: '"' + Config.get("bot_validate_url") + '"',
		ENV: '"' + Config.get("env") + '"',
		HTTPS: '"' + Config.get("https") + '"',
	},
	sourcemap: true,
	bundle: true,
	splitting: true,
	outdir: "public",
	format: 'esm',
	outbase: 'client/gui/',
	plugins: [postcss()],
})

await ctx.watch()
if (Config.get("env") != "development") {
	await ctx.dispose()
}