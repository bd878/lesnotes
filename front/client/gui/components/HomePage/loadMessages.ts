import i18n from '../../i18n';
import api from '../../api';

async function loadMessages(limit, offset, order) {
  let response = [];
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
  } catch (e) {
    console.error(i18n("error_occured"), e);
    throw e
  }

  if (Array.isArray(response.messages)) {
    return response
  }

  return {
    messages: [],
    islastpage: true,
  };
}

export default loadMessages;
