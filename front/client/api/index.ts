import readMessages from './readMessages';
import readStackJson from './readStackJson';
import readMessagesAround from './readMessagesAround';
import readMessagesAroundJson from './readMessagesAroundJson';
import readMessagesJson from './readMessagesJson';
import readPathJson from './readPathJson';
import readBatchMessagesJson from './readBatchMessagesJson';
import readMessage from './readMessage';
import readMessageJson from './readMessageJson';
import sendMessage from './sendMessage';
import sendMessageJson from './sendMessageJson';
import updateMessage from './updateMessage';
import privateThread from './privateThread';
import publishThread from './publishThread';
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
import getMeJson from './getMeJson';
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
	getMeJson,
	readMessages,
	readStackJson,
	readMessagesAround,
	readMessagesAroundJson,
	readMessagesJson,
	readBatchMessagesJson,
	readPathJson,
	readMessageJson,
	readMessage,
	sendMessage,
	sendMessageJson,
	updateMessage,
	deleteMessage,
	deleteMessages,
	publishMessages,
	privateMessages,
	privateThread,
	publishThread,
	moveMessage,
	uploadFile,
	getFileDownloadUrl,
	getMessageLinkUrl,
	validateMiniappData,
	validateTgAuthData,
	sendLog,
}