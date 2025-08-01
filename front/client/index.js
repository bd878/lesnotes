import './node_fetch.ts'

import Koa from 'koa';
import Router from '@koa/router';
import Config from 'config';
import helmet from './handlers/helmet.js';
import errors from './handlers/errors.js';
import logger from './handlers/logger.js';
import bodyParser from './handlers/bodyParser.js';
import useragent from './handlers/useragent.js';
import favicon from './handlers/favicon.js';
import etag from './handlers/etag.js';

import assets from './routes/assets.js';
import main from './routes/main.js';
import login from './routes/login.js';
import register from './routes/register.js';
import home from './routes/home.js';
import status from './routes/status.js';
import xxx from './routes/xxx.js';
import message from './routes/message.js';
import miniapp from './routes/miniapp.js';
import authTelegram from './routes/authTelegram.js';

const app = new Koa();
const router = new Router();

// app.use(helmet);
app.use(errors);
app.use(logger);
app.use(bodyParser);
app.use(useragent);
app.use(favicon);

router
	.get('/public/:filename', etag, assets)
	.get('/index', ctx => {
		ctx.redirect('/')
		ctx.status = 301
	})
	.get('/', etag, main)
	.get('/login', etag, login)
	.get('/signup', etag, register)
	.get('/home', etag, home)
	.get('/status', status)
	.get("/new", etag, message)
	.get("/tg_auth", authTelegram)
	.get("/m/:id", etag, message)
	.get("/miniapp", etag, miniapp)
	.get('/:any*', xxx)

app.use(router.routes());

const port = process.env.PORT || Config.get("port") || 8080;
const host = process.env.HOST || Config.get("addr") || "localhost";

app.listen(port, host, () => {
	console.log(`App is listening on ${port} port`);
});
