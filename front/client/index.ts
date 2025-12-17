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
import notAuthed from './handlers/notAuthed';
import authed from './handlers/authed';
import loadStack from './handlers/loadStack.js';
import loadMessage from './handlers/loadMessage.js';
import loadSearch from './handlers/loadSearch';
import formatMessage from './handlers/formatMessage';
import loadSearchPath from './handlers/loadSearchPath';
import getState from './handlers/getState';
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
import publicMessage from './routes/publicMessage/publicMessage';
import messageView from './routes/messageView/messageView';
import messageEdit from './routes/messageEdit/messageEdit';
import miniapp from './routes/miniapp/miniapp';
import authTelegram from './routes/authTelegram/authTelegram';
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
	.get('main', '/', etag, getState, getToken, notAuthed, main)
	.get('login', '/login', etag, getState, getToken, notAuthed, login)
	.post('/login', etag, getState, validateLogin, redirectHome)
	.get('logout', '/logout', etag, getState, expireToken, redirectLogin)
	.get('signup', '/signup', etag, getState, getToken, notAuthed, signup)
	.post('/signup', etag, getState, validateSignup, redirectHome)
	.get('home', '/home', etag, getToken, getMe, getState, loadStack, home)
	.get('message', '/messages/:id', etag, getToken, getMe, getState, loadStack, loadMessage, formatMessage, messageView)
	.get('editMessage', '/editor/messages/:id', etag, getToken, getMe, getState, loadStack, loadMessage, formatMessage, messageEdit)
	.get('search', '/search', etag, getToken, getMe, getState, loadSearch, loadSearchPath, search)
	.get('status', '/status', status, getState)
	.post("/delete", getToken, authed, getState, deleteMessage)
	.post("/publish", getToken, authed, getState, publishMessage)
	.post("/private", getToken, authed, getState, privateMessage)
	.post("/send", getToken, authed, getState, sendMessage)
	.post("/update", getToken, authed, getState, updateMessage)
	.get("/tg_auth", authTelegram)
	.get("/m/:user/:id", etag, getState, getToken, loadMessage, publicMessage)
	.get("/m/:name", etag, getState, getToken, loadMessage, publicMessage)
	.get("/miniapp", etag, getState, miniapp)
	.get('/:any*', getState, xxx)

app.use(router.routes());

const port = process.env.PORT || Config.get("port") || 8080;
const host = process.env.HOST || Config.get("addr") || "localhost";

app.listen(port, host, () => {
	console.log(`App is listening on ${port} port`);
});
