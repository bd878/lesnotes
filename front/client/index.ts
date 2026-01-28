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
import noCache from './handlers/noCache';
import loadStack from './handlers/loadStack';
import loadFiles from './handlers/loadFiles';
import selectMessageFiles from './handlers/selectMessageFiles';
import loadMessage from './handlers/loadMessage';
import loadThread from './handlers/loadThread';
import loadSearch from './handlers/loadSearch';
import formatTextarea from './handlers/formatTextarea';
import formatView from './handlers/formatView';
import loadSearchPath from './handlers/loadSearchPath';
import loadThreadMessages from './handlers/loadThreadMessages';
import getState from './handlers/getState';
import expireToken from './handlers/expireToken';
import redirectHome from './handlers/redirectHome';
import redirectLogin from './handlers/redirectLogin';
import validateLogin from './handlers/validateLogin'
import validateSignup from './handlers/validateSignup';
import deleteFile from './handlers/deleteFile';
import publishFile from './handlers/publishFile';
import privateFile from './handlers/privateFile';
import deleteMessage from './handlers/deleteMessage';
import publishMessage from './handlers/publishMessage';
import privateMessage from './handlers/privateMessage';
import publishThread from './handlers/publishThread';
import privateThread from './handlers/privateThread';
import sendMessage from './handlers/sendMessage';
import updateMessage from './handlers/updateMessage';
import updateThread from './handlers/updateThread';
import getSearchForm from './handlers/getSearchForm';
import getSearchQuery from './handlers/getSearchQuery';

import assets from './routes/assets';
import main from './routes/main';
import login from './routes/login';
import signup from './routes/signup';
import newMessage from './routes/newMessage';
import files from './routes/files';
import search from './routes/search';
import xxx from './routes/xxx';
import publicMessage from './routes/publicMessage';
import publicThread from './routes/publicThread';
import publicThreadMessage from './routes/publicThreadMessage';
import threadEdit from './routes/threadEdit';
import messageView from './routes/messageView';
import threadView from './routes/threadView';
import messageEdit from './routes/messageEdit';
import authTelegram from './routes/authTelegram';
import status from './routes/status';

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
	.get("main",                   "/",                             etag, noCache, getState, getToken, notAuthed, main)
	.get("login",                  "/login",                        etag, noCache, getState, getToken, notAuthed, login)
	.get("logout",                 "/logout",                       etag, noCache, getState, expireToken, redirectLogin)
	.get("signup",                 "/signup",                       etag, noCache, getState, getToken, notAuthed, signup)
	.get("home",                   "/home",                         etag, noCache, getToken, authed, getMe, getState, loadStack, loadFiles, newMessage)
	.get("files",                  "/files",                        etag, noCache, getToken, authed, getMe, getState, loadStack, loadFiles, files)
	.get("message",                "/messages/:id",                 etag, noCache, getToken, authed, getMe, getState, loadStack, loadMessage, formatView, messageView)
	.get("thread",                 "/threads/:id",                  etag, noCache, getToken, authed, getMe, getState, loadStack, loadThread, formatView, threadView)
	.get("editMessage",            "/editor/messages/:id",          etag, noCache, getToken, authed, getMe, getState, loadStack, loadMessage, loadFiles, selectMessageFiles, formatTextarea, messageEdit)
	.get("editThread",             "/editor/threads/:id",           etag, noCache, getToken, authed, getMe, getState, loadStack, loadThread, formatTextarea, threadEdit)
	.get("status",                 "/status",                       status, noCache, getState)
	.get("search",                 "/search",                       etag, noCache, getToken, authed, getState, getMe, getSearchQuery, loadSearch, loadSearchPath, search)
	.post("doLogin",               "/login",                        etag, getState, validateLogin, redirectHome)
	.post("doSignup",              "/signup",                       etag, getState, validateSignup, redirectHome)
	.post("doSearch",              "/search",                       etag, getToken, authed, getState, getMe, getSearchForm, loadSearch, loadSearchPath, search)
	.post("doDeleteFile",          "/f/delete",                     getToken, authed, getState, deleteFile)
	.post("doPublishFile",         "/f/publish",                    getToken, authed, getState, publishFile)
	.post("doPrivateFile",         "/f/private",                    getToken, authed, getState, privateFile)
	.post("doDeleteMessage",       "/m/delete",                     getToken, authed, getState, deleteMessage)
	.post("doPublishMessage",      "/m/publish",                    getToken, authed, getState, publishMessage)
	.post("doPrivateMessage",      "/m/private",                    getToken, authed, getState, privateMessage)
	.post("doUpdateMessage",       "/m/update",                     getToken, authed, getState, updateMessage)
	.post("doPublishThread",       "/t/publish",                    getToken, authed, getState, publishThread)
	.post("doPrivateThread",       "/t/private",                    getToken, authed, getState, privateThread)
	.post("doUpdateThread",        "/t/update",                     getToken, authed, getState, updateThread)
	.post("doSendMessage",         "/send",                         getToken, authed, getState, sendMessage)
	.get("tgAuth",                 "/tg_auth",                      authTelegram)
	.get("publicMessage",          "/m/:messageName",               etag, noCache, getState, getToken, getMe, loadMessage, formatView, publicMessage)
	.get("publicThread",           "/t/:threadName",                etag, noCache, getState, getToken, getMe, loadThread, loadThreadMessages, formatView, publicThread)
	.get("publicThreadMessage",    "/t/:threadName/:messageName",   etag, noCache, getState, getToken, getMe, loadThread, loadThreadMessages, loadMessage, formatView, publicThreadMessage)
	.get("any",                    "/:any*",                        getState, xxx)

app.use(router.routes());

const port = process.env.PORT || Config.get("port") || 8080;
const host = process.env.HOST || Config.get("addr") || "localhost";

app.listen(port, host, () => {
	console.log(`App is listening on ${port} port`);
});
