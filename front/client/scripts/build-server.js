import esbuild from 'esbuild'

await esbuild.build({
  entryPoints: ['front/client/index.js'],
  entryNames: '[dir]',
  bundle: true,
  platform: 'node',
  outdir: "front/public",
  format: 'cjs',
  outbase: 'front'
})
