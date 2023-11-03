import fs from 'node:fs';
import mime from 'mime/lite';
import request from 'request';

function getAssets(ctx) {
  ctx.set({ 'Content-Type': mime.getType(ctx.params.filename) || 'text/plain' });
  ctx.body = fs.createReadStream(`public/${ctx.params.filename}`);
  ctx.status = 200;
}

export default getAssets;
