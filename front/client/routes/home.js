import Config from 'config';
import mustache from 'mustache';
import path from 'path';
import { readFile } from 'node:fs/promises';
import { resolve } from 'node:path';

async function renderer(ctx) {
  try {
    const filePath = resolve(path.join(Config.get('basedir'), 'templates/index.mustache'));
    const template = await readFile(filePath, { encoding: 'utf-8' });

    ctx.set({ 'Cache-Control': 'no-cache, max-age=0' })
    ctx.body = mustache.render(template, {
      script: "/public/home.js",
      styles: [
        "/public/styles.css",
      ],
    });
    ctx.status = 200;
  } catch (err) {
    ctx.body = "<html>Pas de template</html>";
    ctx.status = 500;
    console.log("failed to return index template");
    throw Error(err);
  }
}

export default renderer;
