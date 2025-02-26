import i18n from '../i18n';
import api from './api';

async function deleteMessage(id = "") {
  let response = {};
  let result: DeleteMessageResult = {
    error: "",
    explain: "",
    ID: "",
  }

  try {
    response = await api("/messages/v1/delete", {
      queryParams: {
        id: id,
      },
      method: "DELETE",
      credentials: "include",
    });

    if (response.error != "") {
      result.error = response.error
      result.explain = response.explain
    } else {
      if (response.value) {
        result.ID = id
      }
    }
  } catch (e) {
    console.error(i18n("error_occured"), e);
    throw e
  }

  return result
}

export default deleteMessage