import {unlink, readdir, stat, rm} from "node:fs/promises";
// TODO: import ts from "typescript";
import Config from "config";
import path from "node:path";
import esbuild from 'esbuild';
import {sassPlugin} from 'esbuild-sass-plugin';

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
    'client/gui/styles/*.sass',
  ],
  entryNames: '[name]',
  define: {
    BACKENDURL: '"' + Config.get("backendurl") + '"',
    ENV: '"' + Config.get("env") + '"',
  },
  sourcemap: true,
  bundle: true,
  splitting: true,
  outdir: "public",
  format: 'esm',
  outbase: 'client/gui/',
  plugins: [sassPlugin()],
})

await ctx.watch()
if (Config.get("env") == "development") {
  console.log("watching...")
} else {
  await ctx.dispose()
}