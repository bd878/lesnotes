import esbuild from 'esbuild'

await esbuild.build({
  entryPoints: ['index.js'],
  entryNames: '[name]',
  bundle: true,
  platform: 'node',
  outdir: "../build",
  format: 'cjs',
})
