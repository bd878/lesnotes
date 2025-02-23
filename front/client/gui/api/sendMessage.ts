import i18n from '../i18n';
import api from './api';
import models from './models';

async function sendMessage(message = "", file = null) {
  let response = {};
  let result: SendMessageResult = {
    error: "",
    explain: "",
    message: models.message(),
  }

  const form = new FormData()
  form.append("text", message);
  if (file != null && file.name != "") {
    form.append('file', file, file.name);
  }

  try {
    response = await api("/messages/v1/send", {
      method: "POST",
      credentials: "include",
      body: form,
    });

    if (response.error != "") {
      result.error = response.error
      result.explain = response.explain
    } else {
      if (response.value != undefined && response.value.message != undefined) { 
        result.message = models.message(response.value.message)
      }
    }
  } catch (e) {
    console.error(i18n("error_occured"), e);
    throw e
  }

  return result
}

export default sendMessage