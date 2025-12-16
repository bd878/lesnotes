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
import authed from './handlers/authed';
import loadStack from './handlers/loadStack.js';
import loadMessage from './handlers/loadMessage.js';
import loadSearch from './handlers/loadSearch';
import getEditorMode from './handlers/getEditorMode';
import formatMessage from './handlers/formatMessage';
import loadSearchPath from './handlers/loadSearchPath';
import getLanguage from './handlers/getLanguage';
import getFontSize from './handlers/getFontSize';
import expireToken from './handlers/expireToken';
import redirectHome from './handlers/redirectHome';
import redirectLogin from './handlers/redirectLogin';
import validateLogin from './handlers/validateLogin'
import validateSignup from './handlers/validateSignup';
import deleteMessage from './handlers/deleteMessage';
import publishMessage from './handlers/publishMessage';
import privateMessage from './handlers/privateMessage';
import sendMessage from './handlers/sendMessage';
import updateMessage from './handlers/updateMessage';

import assets from './routes/assets/assets';
import main from './routes/main/main';
import login from './routes/login/login';
import signup from './routes/signup/signup';
import home from './routes/home/home';
import search from './routes/search/search';
import xxx from './routes/xxx/xxx';
import message from './routes/message/message';
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
	.get('/', etag, getLanguage, getTheme, getFontSize, getToken, notAuthed, main)
	.get('/login', etag, getLanguage, getFontSize, getTheme, getToken, notAuthed, login)
	.post('/login', etag, getLanguage, getFontSize, getTheme, validateLogin, redirectHome)
	.get('/logout', etag, getLanguage, getFontSize, getTheme, expireToken, redirectLogin)
	.get('/signup', etag, getLanguage, getFontSize, getToken, getTheme, notAuthed, signup)
	.post('/signup', etag, getLanguage, getFontSize, getTheme, validateSignup, redirectHome)
	.get('/home', etag, getToken, getMe, getLanguage, getFontSize, getTheme, loadMessage, loadStack, getEditorMode, formatMessage, home)
	.get('/search', etag, getToken, getMe, getLanguage, loadSearch, loadSearchPath, search)
	.get('/status', status, getLanguage)
	.post("/delete", getToken, authed, getLanguage, getFontSize, getTheme, deleteMessage, redirectHome)
	.post("/publish", getToken, authed, getLanguage, getFontSize, getTheme, publishMessage, redirectHome)
	.post("/private", getToken, authed, getLanguage, getFontSize, getTheme, privateMessage, redirectHome)
	.post("/send", getToken, authed, getLanguage, getFontSize, getTheme, sendMessage, redirectHome)
	.post("/update", getToken, authed, getLanguage, getFontSize, getTheme, updateMessage, redirectHome)
	.get("/tg_auth", authTelegram)
	.get("/m/:user/:id", etag, getLanguage, getFontSize, getTheme, getToken, loadMessage, message)
	.get("/m/:name", etag, getLanguage, getFontSize, getTheme, getToken, loadMessage, message)
	.get("/miniapp", etag, getLanguage, miniapp)
	.get('/:any*', getLanguage, xxx)

app.use(router.routes());

const port = process.env.PORT || Config.get("port") || 8080;
const host = process.env.HOST || Config.get("addr") || "localhost";

app.listen(port, host, () => {
	console.log(`App is listening on ${port} port`);
});
