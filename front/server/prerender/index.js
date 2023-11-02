async function prerender(ctx) {
  ctx.status = 200;
  ctx.body = '<html><body><div>Bonjour tous les monde!</div></body></html>\n';
}

export default prerender;
