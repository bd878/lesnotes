import esbuild from 'esbuild'

await esbuild.build({
  entryPoints: ['client/index.js'],
  entryNames: '[name]',
  bundle: true,
  platform: 'node',
  outdir: "build",
  outbase: "client",
  format: 'cjs',
})
