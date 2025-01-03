import i18n from '../../i18n';
import api from '../../api';

async function loadMessages(limit, offset, order) {
  let response = {};
  let result = {
    error: "",
    explain: "",
    messages: [],
    islastpage: false,
  }

  try {
    response = await api('/messages/v1/read', {
      queryParams: {
        limit: limit,
        offset: offset,
        asc: order,
      },
      method: "GET",
      credentials: 'include',
    });

    if (response.error != "") {
      console.error('[loadMessages]: /read response returned error')
      result.error = response.error
      result.explain = response.explain
    } else {
      result.messages = response.value.messages
      result.islastpage = response.value.islastpage
    }
  } catch (e) {
    console.error(i18n("error_occured"), e);
    throw e
  }

  if (Array.isArray(result.messages)) {
    return result
  }

  return result;
}

export default loadMessages;
