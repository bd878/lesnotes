import esbuild from 'esbuild'

await esbuild.build({
  entryPoints: ['front/client/gui/pages/**/index.jsx'],
  entryNames: '[dir]',
  bundle: true,
  splitting: true,
  outdir: "front/public",
  format: 'esm',
  loader: { '.js': 'jsx' },
  outbase: 'front/client/gui/pages'
})
