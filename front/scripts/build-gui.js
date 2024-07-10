import {unlink, readdir} from "node:fs/promises";
import Config from "config";
import path from "node:path";
import esbuild from 'esbuild';
import {sassPlugin} from 'esbuild-sass-plugin';

/*remove stale files except favicon.ico*/
const files = await readdir('public', {withFileTypes: true});
for (const file of files) {
  const ext = path.extname(file.name)
  if (ext == ".js" || ext == ".css")
    await unlink('public/' + file.name)
}

await esbuild.build({
  entryPoints: ['client/gui/pages/**/index.jsx'],
  entryNames: '[dir]',
  define: {
    BACKENDURL: '"' + Config.get("backendurl") + '"',
    ENV: '"' + Config.get("env") + '"',
  },
  bundle: true,
  splitting: true,
  outdir: "public",
  format: 'esm',
  loader: { '.js': 'jsx' },
  outbase: 'client/gui/pages',
  plugins: [sassPlugin()],
})

await esbuild.build({
  entryPoints: ['client/gui/styles/*.sass'],
  bundle: true,
  outdir: "public",
  format: 'esm',
  outbase: 'client/gui/styles',
  plugins: [sassPlugin()],
})