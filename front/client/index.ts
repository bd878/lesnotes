import './node_fetch.ts'

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
import isAuthed from './handlers/isAuthed';
import authed from './handlers/authed';
import noCache from './handlers/noCache';
import loadTree from './handlers/loadTree';
import loadComments from './handlers/loadComments';
import messageFeatures from './handlers/messageFeatures'
import loadMessagePath from './handlers/loadMessagePath';
import loadCwdPath from './handlers/loadCwdPath';
import loadCwdThread from './handlers/loadCwdThread';
import loadFiles from './handlers/loadFiles';
import selectMessageFiles from './handlers/selectMessageFiles';
import loadTranslation from './handlers/loadTranslation';
import loadMessage from './handlers/loadMessage';
import loadThread from './handlers/loadThread';
import formatTextarea from './handlers/formatTextarea';
import formatView from './handlers/formatView';
import getState from './handlers/getState';
import expireToken from './handlers/expireToken';
import redirectHome from './handlers/redirectHome';
import redirectLogin from './handlers/redirectLogin';
import validateLogin from './handlers/validateLogin'
import validateSignup from './handlers/validateSignup';
import deleteMessage from './handlers/deleteMessage';
import deleteTranslation from './handlers/deleteTranslation';
import updateTranslation from './handlers/updateTranslation';
import publishMessage from './handlers/publishMessage';
import privateMessage from './handlers/privateMessage';
import sendMessage from './handlers/sendMessage';
import sendComment from './handlers/sendComment';
import sendTranslation from './handlers/sendTranslation';
import loadParentMessage from './handlers/loadParentMessage';

import assets from './routes/assets';
import main from './routes/main';
import login from './routes/login';
import signup from './routes/signup';
import newMessage from './routes/newMessage';
import xxx from './routes/xxx';
import publicThreadOrMessage from './routes/publicThreadOrMessage';
import publicThreadMessage from './routes/publicThreadMessage';
import messageView from './routes/messageView';
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
	.get("home",                   "/home",                         etag, noCache, getMe, getState, authed, loadCwdThread, loadTree, loadCwdPath, loadFiles, newMessage)
	.get("message",                "/messages/:idOrName",           etag, noCache, getMe, getState, authed, loadCwdThread, loadTree, loadMessagePath, loadCwdPath, loadThread, loadMessage, loadComments, loadTranslation, formatView, messageFeatures, messageView)
	.get("editMessage",            "/editor/messages/:idOrName",    etag, noCache, getMe, getState, authed, loadCwdThread, loadTree, loadMessagePath, loadCwdPath, loadMessage, loadFiles, loadComments, selectMessageFiles, loadTranslation, formatTextarea, messageFeatures, messageEdit)
	.get("status",                 "/status",                       status, noCache, getState)
	.post("doLogin",               "/login",                        etag, getState, validateLogin, redirectHome)
	.post("doSignup",              "/signup",                       etag, getState, validateSignup, redirectHome)
	.post("doDeleteMessage",       "/message/delete",               getState, authed, deleteMessage)
	.post("doPublishMessage",      "/message/publish",              getState, authed, publishMessage)
	.post("doPrivateMessage",      "/message/private",              getState, authed, privateMessage)
	.post("doDeleteTranslation",   "/translation/delete",           getState, authed, deleteTranslation)
	.post("doUpdateTranslation",   "/translation/update",           getState, authed, updateTranslation)
	.post("doSendComment",         "/comment/send",                 getState, sendComment)
	.post("doSendTranslation",     "/translation/send",             getState, authed, sendTranslation)
	.get("publicMessage",          "/:messageOrParentName",         etag, noCache, getMe, getState, loadMessage, isAuthed(loadMessagePath), loadComments, loadTree, formatView, messageFeatures, publicThreadOrMessage)
	.get("publicThreadMessage",    "/:parentName/:messageName",     etag, noCache, getMe, getState, loadMessage, isAuthed(loadMessagePath), loadComments, loadParentMessage, loadTree, formatView, messageFeatures, publicThreadMessage)
	.get("any",                    "/:any*",                        getState, xxx)

app.use(router.routes());

const port = process.env.PORT || Config.get("port") || 8080;
const host = process.env.HOST || Config.get("addr") || "localhost";

app.listen(port, host, () => {
	console.log(`App is listening on ${port} port`);
});
