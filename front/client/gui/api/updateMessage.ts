import i18n from '../i18n';
import api from './api';
import models from './models';

async function updateMessage(id = "", message = "") {
  let response = {};
  let result: UpdateMessageResult = {
    error: "",
    explain: "",
    ID: "",
    updateUTCNano: 0,
  }

  const form = new FormData()
  form.append("text", message);

  try {
    response = await api("/messages/v1/update", {
      queryParams: {
        id: id,
      },
      method: "POST",
      credentials: "include",
      body: form,
    });

    if (response.error != "") {
      result.error = response.error
      result.explain = response.explain
    } else {
      if (response.value && response.value.id) {
        const model = models.message({ID: response.value.id, update_utc_nano: response.value.update_utc_nano})
        result.ID = model.ID
        result.updateUTCNano = model.updateUTCNano
      }
    }
  } catch (e) {
    console.error(i18n("error_occured"), e);
    throw e
  }

  return result
}

export default updateMessage