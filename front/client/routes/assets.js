import fs from 'node:fs';
import mime from 'mime';
import request from 'request';

function getAssets(ctx) {
  ctx.set({ 'Content-Type': mime.getType(ctx.params.filename) || 'text/plain' });
  ctx.set({ 'Cache-Control': 'max-age=604800', 'ETag': '2' })
  ctx.body = fs.createReadStream(`front/public/${ctx.params.filename}`);
  ctx.status = 200;
}

export default getAssets;
