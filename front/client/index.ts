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
import getTheme from './handlers/getTheme.js';
import notAuthed from './handlers/notAuthed';
import loadStack from './handlers/loadStack.js';
import loadMessage from './handlers/loadMessage.js';
import loadSearch from './handlers/loadSearch';
import getEditorMode from './handlers/getEditorMode';
import formatMessage from './handlers/formatMessage';
import loadSearchPath from './handlers/loadSearchPath';
import getLanguage from './handlers/getLanguage';
import getFontSize from './handlers/getFontSize';
import getQuery from './handlers/getQuery';

import assets from './routes/assets/assets';
import main from './routes/main/main';
import login from './routes/login/login';
import register from './routes/register/register';
import home from './routes/home/home';
import search from './routes/search/search';
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
	.get('/', etag, getLanguage, getTheme, getFontSize, getToken, getQuery, notAuthed, main)
	.get('/login', etag, getLanguage, getFontSize, getTheme, getQuery, getToken, notAuthed, login)
	.get('/logout', etag, getLanguage, getFontSize, getTheme, getQuery, logout)
	.get('/signup', etag, getLanguage, getFontSize, getToken, getTheme, getQuery, notAuthed, register)
	.get('/home', etag, getLanguage, getToken, getMe, loadMessage, loadStack, getEditorMode, formatMessage, home)
	.get('/search', etag, getLanguage, getToken, getMe, loadSearch, loadSearchPath, search)
	.get('/status', status, getLanguage)
	.get("/tg_auth", authTelegram)
	.get("/new", etag, getLanguage, newMessage)
	.get("/m/:user/:id", etag, getLanguage, getToken, loadMessage, message)
	.get("/m/:name", etag, getLanguage, getFontSize, getTheme, getToken, loadMessage, message)
	.get("/miniapp", etag, getLanguage, miniapp)
	.get('/:any*', getLanguage, xxx)

app.use(router.routes());

const port = process.env.PORT || Config.get("port") || 8080;
const host = process.env.HOST || Config.get("addr") || "localhost";

app.listen(port, host, () => {
	console.log(`App is listening on ${port} port`);
});
