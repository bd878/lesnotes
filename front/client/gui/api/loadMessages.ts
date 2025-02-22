import i18n from '../i18n';
import api from './api';
import models from './models';

interface LoadMessagesResult {
  error: string;
  explain: string;
  messages: any[];
  isLastPage: boolean;
}

async function loadMessages(limit: number, offset: number, order: number): LoadMessagesResult {
  let response = {};
  let result: LoadMessagesResult = {
    error: "",
    explain: "",
    messages: [],
    isLastPage: false,
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
      result.messages = response.value.messages.map(models.message)
      result.isLastPage = response.value.is_last_page
    }
  } catch (e) {
    console.error(i18n("error_occured"), e);
    throw e
  }

  return result;
}

export default loadMessages;
