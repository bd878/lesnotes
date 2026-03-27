import './node_fetch'

import Koa from 'koa';
import Router from '@koa/router';
import Config from 'config';
import helmet from './handlers/helmet';
import errors from './handlers/errors';
import logger from './handlers/logger';
import bodyParser from './handlers/bodyParser';
import useragent from './handlers/useragent';
import favicon from './handlers/favicon';
import etag from './handlers/etag';
import getMe from './handlers/getMe';
import notAuthed from './handlers/notAuthed';
import authed from './handlers/authed';
import noCache from './handlers/noCache';
import loadTree from './handlers/loadTree';
import loadComments from './handlers/loadComments';
import messageFeatures from './handlers/messageFeatures'
import messageTranslations from './handlers/messageTranslations'
import loadPath from './handlers/loadPath';
import loadFiles from './handlers/loadFiles';
import selectMessageFiles from './handlers/selectMessageFiles';
import loadTranslation from './handlers/loadTranslation';
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
import deleteTranslation from './handlers/deleteTranslation';
import updateTranslation from './handlers/updateTranslation';
import publishMessage from './handlers/publishMessage';
import privateMessage from './handlers/privateMessage';
import publishThread from './handlers/publishThread';
import privateThread from './handlers/privateThread';
import sendMessage from './handlers/sendMessage';
import sendComment from './handlers/sendComment';
import sendTranslation from './handlers/sendTranslation';
import updateMessage from './handlers/updateMessage';
import updateThread from './handlers/updateThread';
import getSearchForm from './handlers/getSearchForm';
import getSearchQuery from './handlers/getSearchQuery';
import parseMessageName from './handlers/parseMessageName';
import parseThreadName from './handlers/parseThreadName';

import assets from './routes/assets';
import main from './routes/main';
import login from './routes/login';
import signup from './routes/signup';
import newMessage from './routes/newMessage';
import search from './routes/search';
import xxx from './routes/xxx';
import publicMessage from './routes/publicMessage';
import publicTranslation from './routes/publicTranslation';
import publicThread from './routes/publicThread';
import publicThreadMessage from './routes/publicThreadMessage';
import threadEdit from './routes/threadEdit';
import messageView from './routes/messageView';
import threadView from './routes/threadView';
import messageEdit from './routes/messageEdit';
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
	.get("main",                   "/",                             etag, noCache, getState, notAuthed, main)
	.get("login",                  "/login",                        etag, noCache, getState, notAuthed, login)
	.get("logout",                 "/logout",                       etag, noCache, getState, expireToken, redirectLogin)
	.get("signup",                 "/signup",                       etag, noCache, getState, notAuthed, signup)
	.get("home",                   "/home",                         etag, noCache, getState, authed, getMe, loadTree, loadFiles, newMessage)
	.get("message",                "/messages/:id",                 etag, noCache, getState, authed, getMe, loadTree, loadPath, loadThread, loadMessage, loadComments, loadTranslation, formatView, messageFeatures, messageTranslations, messageView)
	.get("editMessage",            "/editor/messages/:id",          etag, noCache, getState, authed, getMe, loadTree, loadPath, loadMessage, loadFiles, loadComments, selectMessageFiles, loadTranslation, formatTextarea, messageFeatures, messageTranslations, messageEdit)
	.get("thread",                 "/threads/:id",                  etag, noCache, getState, authed, getMe, loadTree, loadPath, loadThread, formatView, threadView)
	.get("newThreadMessage",       "/editor/messages/:id/new",      etag, noCache, getState, authed, getMe, loadTree, loadPath, newMessage)
	.get("editThread",             "/editor/threads/:id",           etag, noCache, getState, authed, getMe, loadTree, loadPath, loadThread, formatTextarea, threadEdit)
	.get("status",                 "/status",                       status, noCache, getState)
	.get("search",                 "/search",                       etag, noCache, getState, authed, getMe, getSearchQuery, loadSearch, loadSearchPath, search)
	.post("doLogin",               "/login",                        etag, getState, validateLogin, redirectHome)
	.post("doSignup",              "/signup",                       etag, getState, validateSignup, redirectHome)
	.post("doSearch",              "/search",                       etag, getState, authed, getMe, getSearchForm, loadSearch, loadSearchPath, search)
	.post("doDeleteFile",          "/f/delete",                     getState, authed, deleteFile)
	.post("doPublishFile",         "/f/publish",                    getState, authed, publishFile)
	.post("doPrivateFile",         "/f/private",                    getState, authed, privateFile)
	.post("doDeleteMessage",       "/m/delete",                     getState, authed, deleteMessage)
	.post("doDeleteTranslation",   "/translation/delete",           getState, authed, deleteTranslation)
	.post("doPublishMessage",      "/m/publish",                    getState, authed, publishMessage)
	.post("doPrivateMessage",      "/m/private",                    getState, authed, privateMessage)
	.post("doUpdateMessage",       "/m/update",                     getState, authed, updateMessage)
	.post("doUpdateTranslation",   "/translation/update",           getState, authed, updateTranslation)
	.post("doPublishThread",       "/t/publish",                    getState, authed, publishThread)
	.post("doPrivateThread",       "/t/private",                    getState, authed, privateThread)
	.post("doUpdateThread",        "/t/update",                     getState, authed, updateThread)
	.post("doSendComment",         "/comment/send",                 getState, sendComment)
	.post("doSendMessage",         "/send",                         getState, authed, sendMessage)
	.post("doSendTranslation",     "/translation/send",             getState, authed, sendTranslation)
	.get("publicMessage",          "/m/:messageName",               etag, noCache, getState, getMe, loadMessage, loadComments, parseMessageName, formatView, publicMessage)
	.get("publicTranslation",      "/m/:messageName/:lang",         etag, noCache, getState, getMe, loadMessage, loadComments, parseMessageName, loadTranslation, formatView, publicTranslation)
	.get("publicThread",           "/t/:threadName",                etag, noCache, getState, getMe, loadThread, parseThreadName, loadThreadMessages, formatView, publicThread)
	.get("publicThreadMessage",    "/t/:threadName/:messageName",   etag, noCache, getState, getMe, loadThread, loadComments, parseThreadName, parseMessageName, loadThreadMessages, loadMessage, formatView, publicThreadMessage)
	.get("any",                    "/:any*",                        getState, xxx)

app.use(router.routes());

const port = process.env.PORT || Config.get("port") || 8080;
const host = process.env.HOST || Config.get("addr") || "localhost";

app.listen(port, host, () => {
	console.log(`App is listening on ${port} port`);
});
