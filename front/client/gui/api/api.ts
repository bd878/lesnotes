const methodsWithBody = [
  'POST', 'PUT', 'DELETE', 'PATCH',
]

let proto = "http";
if (ENV == "production" || ENV == "staging") {
  proto = "https"
}

function getFullUrl(url: string, isFullUrl: string) {
  return isFullUrl ? url : `${proto}://${BACKENDURL}${url}`;
}

function getFileDownloadUrl(url: string) {
  return getFullUrl(url, false);
}

function appendQueryParams(url: string, queryParams): string {
  if (!queryParams) {
    return url;
  }
  if (!(queryParams instanceof URLSearchParams)) {
    return url;
  }
  if (queryParams.size == 0) {
    return url;
  }

  return url + "?" + queryParams.toString();
}

function prepareBody(body, method) {
  if (!methodsWithBody.includes(method)) {
    return;
  }

  if (body instanceof URLSearchParams || body instanceof FormData) {
    return body;
  }

  return JSON.stringify({ ...body });
}

function getOptions(props) {
  const {
    headers,
    body,
    method,
    credentials,
  } = props;

  return {
    headers: new Headers({ ...headers }),
    body: prepareBody(body, method),
    mode: 'no-cors',
    method,
    credentials,
  };
}

export {getFileDownloadUrl};

export default function api(url, props = {}) {
  const { isFullUrl = false, queryParams } = props;

  let fullUrl = getFullUrl(url, isFullUrl);
  fullUrl = appendQueryParams(fullUrl, new URLSearchParams(queryParams));
  const options = getOptions(props);

  return fetch(fullUrl, options)
    .then(res => {
      return res
        .text()
        .then(text => {
          try {
            return JSON.parse(text);
          } catch (e) {
            return text;
          }
        })
    })
    .catch(e => Promise.reject({
      error: e,
      message: 'Cannot parse the response',
    }));
}
