import api, {getFileDownloadUrl} from "./api";
import loadMessages from './loadMessages';
import sendMessage from './sendMessage';
import auth from './auth';
import login from './login';
import logout from './logout';
import register from './register';

export {getFileDownloadUrl}
export default {
  register,
  login,
  logout,
  auth,
  loadMessages,
  sendMessage,
  getFileDownloadUrl,
  api,
}