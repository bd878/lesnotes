import Config from 'config';
import path from 'path';
import { readFile, stat } from 'node:fs/promises';
import { resolve } from 'node:path';

export default async function etagMiddleware(ctx, next) {
	try {
      const etagFile = resolve(path.join(Config.get("basedir"), "etag"));
      const etag = await readFile(etagFile, {encoding: 'utf-8'});
      const etagStat = await stat(etagFile)
      const mtime = etagStat.mtime

      ctx.set({ "ETag": etag })
      ctx.set({ "Last-Modified": mtime })
	} catch (e) {
		console.log("cannot set etag", e)
	}

	await next()
}