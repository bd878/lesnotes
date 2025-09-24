import Config from 'config';
import path from 'path';
import fs from 'node:fs';
import mime from 'mime';

function getAssets(ctx) {
	ctx.set({ 'Content-Type': mime.getType(ctx.params.path) || 'text/plain' });
	ctx.set({ 'Cache-Control': 'no-cache, max-age=0' })
	ctx.body = fs.createReadStream(path.join(Config.get('basedir'), `public/${ctx.params.path}`));
	ctx.status = 200;
}

export default getAssets;
