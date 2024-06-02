import esbuild from 'esbuild'

await esbuild.build({
  entryPoints: ['client/gui/pages/**/index.jsx'],
  entryNames: '[dir]',
  bundle: true,
  splitting: true,
  outdir: "public",
  format: 'esm',
  loader: { '.js': 'jsx' },
  outbase: 'client/gui/pages'
})
