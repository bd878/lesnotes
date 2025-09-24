import Config from 'config';
import path from 'path';
import mustache from 'mustache';
import { readFile } from 'node:fs/promises';
import { resolve } from 'node:path';

async function renderer(ctx) {
	console.log("failed to return index template");
	ctx.body = "<html>Pas de template</html>";
	ctx.status = 500;
}

export default renderer;
