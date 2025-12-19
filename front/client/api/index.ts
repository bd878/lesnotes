import readMessages from './readMessages';
import readStackJson from './readStackJson';
import readMessagesJson from './readMessagesJson';
import readPathJson from './readPathJson';
import readBatchMessagesJson from './readBatchMessagesJson';
import readMessage from './readMessage';
import readMessageJson from './readMessageJson';
import sendMessage from './sendMessage';
import sendMessageJson from './sendMessageJson';
import searchMessagesJson from './searchMessagesJson';
import searchMessagesPathJson from './searchMessagesPathJson';
import updateMessage from './updateMessage';
import updateMessageJson from './updateMessageJson';
import privateThread from './privateThread';
import publishThread from './publishThread';
import deleteMessage from './deleteMessage';
import deleteMessageJson from './deleteMessageJson';
import deleteMessages from './deleteMessages';
import publishMessages from './publishMessages'
import publishMessageJson from './publishMessageJson'
import privateMessages from './privateMessages';
import privateMessageJson from './privateMessageJson';
import moveMessage from './moveMessage';
import uploadFile from './uploadFile';
import validateMiniappData from './validateMiniappData';
import validateTgAuthData from './validateTgAuthData';
import sendLog from './sendLog';
import auth from './auth';
import authJson from './authJson';
import reorderThread from './reorderThread';
import login from './login';
import logout from './logout';
import getMe from './getMe';
import getMeJson from './getMeJson';
import signup from './signup';
import {getFileDownloadUrl, getMessageLinkUrl} from './api';

export {getFileDownloadUrl, getMessageLinkUrl}
export default {
	signup,
	login,
	logout,
	auth,
	authJson,
	getMe,
	getMeJson,
	readMessages,
	readStackJson,
	readMessagesJson,
	readBatchMessagesJson,
	readPathJson,
	readMessageJson,
	readMessage,
	sendMessage,
	sendMessageJson,
	updateMessage,
	updateMessageJson,
	deleteMessage,
	deleteMessageJson,
	deleteMessages,
	publishMessages,
	publishMessageJson,
	privateMessages,
	privateMessageJson,
	privateThread,
	publishThread,
	moveMessage,
	uploadFile,
	getFileDownloadUrl,
	getMessageLinkUrl,
	validateMiniappData,
	validateTgAuthData,
	searchMessagesJson,
	searchMessagesPathJson,
	reorderThread,
	sendLog,
}