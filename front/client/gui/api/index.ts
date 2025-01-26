import api, {getFileDownloadUrl} from "./api";
import loadMessages from './loadMessages';
import sendMessage from './sendMessage';
import auth from './auth';
import login from './login';
import register from './register';

export {getFileDownloadUrl}
export default {
  register,
  login,
  auth,
  loadMessages,
  sendMessage,
  getFileDownloadUrl,
  api,
}