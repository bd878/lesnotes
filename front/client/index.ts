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
import language from './handlers/language';

import assets from './routes/assets/assets';
import main from './routes/main/main';
import login from './routes/login/login';
import register from './routes/register/register';
import home from './routes/home/home';
import logout from './routes/logout/logout';
import xxx from './routes/xxx/xxx';
import message from './routes/message/message';
import newMessage from './routes/new_message/new_message';
import miniapp from './routes/miniapp/miniapp';
import authTelegram from './routes/auth_telegram/auth_telegram';
import status from './routes/status/status';

const app = new Koa();
const router = new Router();

// app.use(helmet);
app.use(errors);
app.use(logger);
app.use(bodyParser);
app.use(useragent);
app.use(favicon);

router
	.get('/public/:path*', etag, assets)
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
	.get("/new", etag, newMessage)
	.get("/m/:user/:id", etag, getToken, loadMessage, message)
	.get("/m/:name", etag, getToken, loadMessage, message)
	.get("/miniapp", etag, miniapp)
	.get('/:any*', xxx)

app.use(router.routes());

const port = process.env.PORT || Config.get("port") || 8080;
const host = process.env.HOST || Config.get("addr") || "localhost";

app.listen(port, host, () => {
	console.log(`App is listening on ${port} port`);
});
