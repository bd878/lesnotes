import Index from '../gui/pages/Index';

import mustache from 'mustache';
import React from 'react';
import ReactDOMServer from 'react-dom/server'; 
import { readFile } from 'node:fs/promises';
import { resolve } from 'node:path';

async function prerender(ctx) {
  try {
    const filePath = resolve('templates/index.mustache');
    const template = await readFile(filePath, { encoding: 'utf-8' });

    ctx.body = mustache.render(template, {
      html: ReactDOMServer.renderToString(<Index />),
      script: "public/gui.js",
    });
    ctx.status = 200;
  } catch (err) {
    ctx.body = "<html>Pas de template</html>";
    ctx.status = 500;
    console.log("failed to return index template");
    throw Error(err);
  }
}

export default prerender;
