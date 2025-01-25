import i18n from '../i18n';
import api from './api';

async function sendMessage(message = "", file = null) {
  let response = {};
  let result: SendMessageResult = {
    error: "",
    explain: "",
    message: "",
  }

  const form = new FormData()
  form.append("message", message);
  if (file != null && file.name != "") {
    form.append('file', file, file.name);
  }

  try {
    const response = await api("/messages/v1/send", {
      method: "POST",
      credentials: "include",
      body: form,
    });
    console.log("[sendMessage] response: ", response);
    if (response.error != "") {
      console.error("[sendMessage]: /send returned error", response.error, response.explain)
      result.error = response.error
      result.explain = response.explain
    } else {
      if (response.value.message != undefined) { 
        result.message = response.value.message
      }
    }
  } catch (e) {
    console.error(i18n("error_occured"), e);
    throw e
  }

  return result
}

export default sendMessage