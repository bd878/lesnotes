import './node_fetch.ts'

import Koa from 'koa';
import pino from 'koa-pino-logger'
import mustache from 'mustache';
import Router from '@koa/router';
import Config from 'config';
import helmet from './handlers/helmet';
import errors from './handlers/errors';
import bodyParser from './handlers/bodyParser';
import useragent from './handlers/useragent';
import favicon from './handlers/favicon';
import etag from './handlers/etag';
import getMe from './handlers/getMe';
import notAuthed from './handlers/notAuthed';
import isAuthed from './handlers/isAuthed';
import authed from './handlers/authed';
import noCache from './handlers/noCache';
import loadTree from './handlers/loadTree';
import messageFeatures from './handlers/messageFeatures'
import loadMessagePath from './handlers/loadMessagePath';
import loadFiles from './handlers/loadFiles';
import loadMessage from './handlers/loadMessage';
import loadThread from './handlers/loadThread';
import formatTextarea from './handlers/formatTextarea';
import formatView from './handlers/formatView';
import getState from './handlers/getState';
import expireToken from './handlers/expireToken';
import redirectHome from './handlers/redirectHome';
import redirectSearch from './handlers/redirectSearch';
import redirectLogin from './handlers/redirectLogin';
import validateLogin from './handlers/validateLogin'
import validateSignup from './handlers/validateSignup';
import deleteMessage from './handlers/deleteMessage';
import deleteTranslation from './handlers/deleteTranslation';
import updateTranslation from './handlers/updateTranslation';
import publishMessage from './handlers/publishMessage';
import privateMessage from './handlers/privateMessage';
import sendMessage from './handlers/sendMessage';
import updateMessage from './handlers/updateMessage';
import sendComment from './handlers/sendComment';
import sendTranslation from './handlers/sendTranslation';
import loadParentMessage from './handlers/loadParentMessage';
import searchMessages from './handlers/searchMessages';
import threadPath from './handlers/threadPath';
import newNoteHref from './handlers/newNoteHref';
import newThreadNoteHref from './handlers/newThreadNoteHref';
import messageEditHref from './handlers/messageEditHref';
import searchMessageEditHref from './handlers/searchMessageEditHref';
import editHrefRedirect from './handlers/editHrefRedirect';
import deleteRedirectUrl from './handlers/deleteRedirectUrl';
import messageParentName from './handlers/messageParentName';
import searchMessageDeleteRedirectUrl from './handlers/searchMessageDeleteRedirectUrl';

import assets from './routes/assets';
import main from './routes/main';
import login from './routes/login';
import signup from './routes/signup';
import newMessage from './routes/newMessage';
import newRootMessage from './routes/newRootMessage';
import home from './routes/home';
import xxx from './routes/xxx';
import publicThreadOrMessage from './routes/publicThreadOrMessage';
import publicThreadMessage from './routes/publicThreadMessage';
import publicMessageTree from './routes/publicMessageTree'
import searchMessageView from './routes/searchMessageView';
import messageView from './routes/messageView';
import messageEdit from './routes/messageEdit';
import searchResult from './routes/searchResult';
import editThreadMessage from './routes/editThreadMessage';
import newThreadMessage from './routes/newThreadMessage'
import status from './routes/status';

const app = new Koa();
const router = new Router();

mustache.templateCache = undefined

// app.use(helmet);
app.use(pino({
	transport: {target: "pino-opentelemetry-transport"},
}));
app.use(errors);
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
	.get("home",                   "/home",                         etag, noCache, getMe, getState, authed, newNoteHref, loadTree, loadFiles, home)
	.get("newRootMessage",         "/new",                          etag, noCache, getMe, getState, authed, newNoteHref, loadTree, newRootMessage)
	.get("message",                "/messages/:idOrName",           etag, noCache, getMe, getState, authed, newNoteHref, loadTree, loadMessagePath, loadThread, loadMessage, messageParentName, threadPath, formatView, messageFeatures, messageView)
	.get("editMessage",            "/messages/:idOrName/edit",      etag, noCache, getMe, getState, authed, newNoteHref, loadTree, loadMessagePath, loadMessage, loadFiles, formatTextarea, messageFeatures, messageEdit)
	.get("status",                 "/status",                       status, noCache, getState)
	.get("search",                 "/search",                       etag, noCache, getMe, getState, newNoteHref, searchMessages, searchResult)
	.get("searchMessage",          "/search/:idOrName",             etag, noCache, getMe, getState, newNoteHref, searchMessages, loadMessagePath, loadMessage, messageParentName, threadPath, loadFiles, searchMessageEditHref, searchMessageDeleteRedirectUrl, formatView, messageFeatures, searchMessageView)
	.post("doLogin",               "/login",                        etag, getState, validateLogin, redirectHome)
	.post("doSignup",              "/signup",                       etag, getState, validateSignup, redirectHome)
	.post("doSendMessage",         "/message/send",                 getState, authed, sendMessage)
	.post("doSearch",              "/message/search",               etag, getState, redirectSearch)
	.post("doDeleteMessage",       "/message/delete",               getState, authed, deleteMessage)
	.post("doPublishMessage",      "/message/publish",              getState, authed, publishMessage)
	.post("doPrivateMessage",      "/message/private",              getState, authed, privateMessage)
	.post("doUpdateMessage",       "/message/update",               getState, authed, updateMessage)
	.post("doDeleteTranslation",   "/translation/delete",           getState, authed, deleteTranslation)
	.post("doUpdateTranslation",   "/translation/update",           getState, authed, updateTranslation)
	.post("doSendComment",         "/comment/send",                 getState, sendComment)
	.post("doSendTranslation",     "/translation/send",             getState, authed, sendTranslation)
	.get("publicMessage",          "/:messageOrParentName",         etag, noCache, getMe, getState, loadMessage, threadPath, newThreadNoteHref, messageEditHref, isAuthed(loadMessagePath), deleteRedirectUrl, loadTree, formatView, messageFeatures, publicThreadOrMessage)
	.get("publicMessageTree",      "/:messageOrParentName/tree",    etag, noCache, getMe, getState, loadMessage, threadPath, newThreadNoteHref, messageEditHref, isAuthed(loadMessagePath), deleteRedirectUrl, loadTree, formatView, messageFeatures, publicMessageTree)
	.get("publicMessageEdit",      "/:messageOrParentName/edit",    etag, noCache, getMe, getState, authed, loadMessage, threadPath, newThreadNoteHref, isAuthed(loadMessagePath), editHrefRedirect, loadTree, loadFiles, formatView, messageFeatures, formatTextarea, editThreadMessage)
	.get("newMessage",             "/:messageOrParentName/new",     etag, noCache, getMe, getState, authed, loadMessage, threadPath, newThreadNoteHref, loadTree, formatView, newThreadMessage)
	.get("publicThreadMessage",    "/:parentName/:messageName",         etag, noCache, getMe, getState, loadMessage, threadPath, newThreadNoteHref, messageEditHref, isAuthed(loadMessagePath), deleteRedirectUrl, loadParentMessage, loadTree, formatView, messageFeatures, publicThreadMessage)
	.get("publicThreadMessageEdit",  "/:parentName/:messageName/edit",  etag, noCache, getMe, getState, authed, loadMessage, threadPath, newThreadNoteHref, isAuthed(loadMessagePath), editHrefRedirect, loadTree, loadFiles, messageFeatures, formatTextarea, editThreadMessage)
	.get("any",                    "/:any*",                        getState, xxx)

app.use(router.routes());

const port = process.env.PORT || Config.get("port") || 8080;
const host = process.env.HOST || Config.get("addr") || "localhost";

app.listen(port, host, () => {
	console.log(`App is listening on ${port} port`)
});
