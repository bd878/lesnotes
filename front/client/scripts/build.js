import esbuild from 'esbuild'

await esbuild.build({
  entryPoints: ['client/gui/index.js'],
  bundle: true,
  loader: { '.js': 'jsx' },
  outfile: 'public/gui.js',
})

await esbuild.build({
  entryPoints: ['client/index.js'],
  bundle: true,
  platform: 'node',
  loader: { '.js': 'jsx' },
  outfile: 'build/client.js',
})