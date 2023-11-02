import fs from 'fs';
import request from 'request';

function getAssets(ctx) {
  ctx.set({ 'Content-Type': "test/plain" });
  ctx.body = "test";
  ctx.status = 200;
}

export default getAssets;
