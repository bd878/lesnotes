import url from 'url';

export default async function getAssetFilename(ctx, next) {
  const path = decodeURI(
    url.parse(ctx.url).pathname,
  );

  ctx.filename = path.split('/')[2];

  await next();
}
