import esbuild from 'esbuild'

await esbuild.build({
  entryPoints: ['client/gui/index.js'],
  bundle: true,
  splitting: true,
  outdir: "public",
  format: 'esm',
  loader: { '.js': 'jsx' }
})
