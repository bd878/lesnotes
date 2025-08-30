import readMessages from './readMessages';
import readMessagesJson from './readMessagesJson';
import readMessage from './readMessage';
import readMessageJson from './readMessageJson';
import sendMessage from './sendMessage';
import sendMessageJson from './sendMessageJson';
import updateMessage from './updateMessage';
import deleteMessage from './deleteMessage';
import deleteMessages from './deleteMessages';
import publishMessages from './publishMessages'
import privateMessages from './privateMessages';
import moveMessage from './moveMessage';
import uploadFile from './uploadFile';
import validateMiniappData from './validateMiniappData';
import validateTgAuthData from './validateTgAuthData';
import sendLog from './sendLog';
import auth from './auth';
import authJson from './authJson';
import login from './login';
import logout from './logout';
import getMe from './getMe';
import register from './register';
import {getFileDownloadUrl, getMessageLinkUrl} from './api';

export {getFileDownloadUrl, getMessageLinkUrl}
export default {
	register,
	login,
	logout,
	auth,
	authJson,
	getMe,
	readMessages,
	readMessagesJson,
	readMessageJson,
	readMessage,
	sendMessage,
	sendMessageJson,
	updateMessage,
	deleteMessage,
	deleteMessages,
	publishMessages,
	privateMessages,
	moveMessage,
	uploadFile,
	getFileDownloadUrl,
	getMessageLinkUrl,
	validateMiniappData,
	validateTgAuthData,
	sendLog,
}