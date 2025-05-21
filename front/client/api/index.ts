import loadMessages from './loadMessages';
import loadOneMessage from './loadOneMessage';
import sendMessage from './sendMessage';
import updateMessage from './updateMessage';
import deleteMessage from './deleteMessage';
import deleteMessages from './deleteMessages';
import publishMessages from './publishMessages'
import privateMessages from './privateMessages';
import moveMessage from './moveMessage';
import uploadFile from './uploadFile';
import auth from './auth';
import login from './login';
import logout from './logout';
import register from './register';
import {getFileDownloadUrl, getMessageLinkUrl} from './api';

export {getFileDownloadUrl, getMessageLinkUrl}
export default {
	register,
	login,
	logout,
	auth,
	loadMessages,
	loadOneMessage,
	sendMessage,
	updateMessage,
	deleteMessage,
	deleteMessages,
	publishMessages,
	privateMessages,
	moveMessage,
	uploadFile,
	getFileDownloadUrl,
	getMessageLinkUrl,
}