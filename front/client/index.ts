import './node_fetch'

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
import getMe from './handlers/getMe.js';
import getToken from './handlers/getToken.js';
import loadStack from './handlers/loadStack.js';
import loadMessage from './handlers/loadMessage.js';

import assets from './routes/assets';
import main from './routes/main';
import login from './routes/login';
import register from './routes/register';
import home from './routes/home';
import logout from './routes/logout';
import status from './routes/status';
import xxx from './routes/xxx';
import readUserMessage from './routes/readUserMessage';
import readPublicMessage from './routes/readPublicMessage';
import createNewMessage from './routes/createNewMessage';
import miniapp from './routes/miniapp';
import authTelegram from './routes/authTelegram';

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
	.get('/logout', etag, logout)
	.get('/signup', etag, register)
	.get('/home', etag, getToken, getMe, loadMessage, loadStack, home)
	.get('/status', status)
	.get("/tg_auth", authTelegram)
	.get("/new", etag, createNewMessage)
	.get("/m/:user/:id", etag, readUserMessage)
	.get("/m/:id", etag, readPublicMessage)
	.get("/miniapp", etag, miniapp)
	.get('/:any*', xxx)

app.use(router.routes());

const port = process.env.PORT || Config.get("port") || 8080;
const host = process.env.HOST || Config.get("addr") || "localhost";

app.listen(port, host, () => {
	console.log(`App is listening on ${port} port`);
});
