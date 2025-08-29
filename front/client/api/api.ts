import i18n from '../i18n';
import * as is from '../third_party/is'

const methodsWithBody = [
	'POST', 'PUT', 'DELETE', 'PATCH',
]

let proto = "http";
if (ENV && (ENV.includes("prod") || ENV.includes("stage"))) {
	proto = "https"
}
if (HTTPS && (HTTPS !== "")) {
	proto = "https"
}

function getFullUrl(url: string, isFullUrl: string) {
	return isFullUrl ? url : `${proto}://${BACKEND_URL}${url}`;
}

function getFileDownloadUrl(url: string) {
	return getFullUrl(url, false);
}

function getMessageLinkUrl(userId: string, id: string) {
	return getFullUrl(`/m/${userId}/${id}`, false)
}

function appendQueryParams(url: string, queryParams): string {
	if (!queryParams)
		return url;
	if (!(queryParams instanceof URLSearchParams))
		return url;
	if (queryParams.size == 0)
		return url;

	return url + "?" + queryParams.toString();
}

function prepareBody(body, method) {
	if (!methodsWithBody.includes(method))
		return;

	if (body instanceof URLSearchParams || body instanceof FormData)
		return body;

	return JSON.stringify(body);
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
		mode: 'cors',
		method,
		credentials,
	};
}

export {getFileDownloadUrl, getMessageLinkUrl};

export default function api(url: string, props = {}) {
	const { isFullUrl = false, queryParams } = props;

	let fullUrl = getFullUrl(url, isFullUrl);
	fullUrl = appendQueryParams(fullUrl, new URLSearchParams(queryParams));
	const options = getOptions(props);

	return fetch(fullUrl, options)
		.then(res => {
			if (!res.ok) {
				console.error(`[api] request to ${url} returned ${res.status} status`)
				return {
					value:   null,
					error:   i18n("bad_status_code"),
					explain: res,
				};
			}

			return res
				.text()
				.then(text => {
					try {
						const value = JSON.parse(text);
						if (value.status === "error" && is.notEmpty(value.error))
							return {
								value,
								error:   true,
								data:    value.data,
								code:    value.error.code,
								explain: value.error.explain,
							}
						else
							return {
								value:   value.response,
								data:    value.data,
								status:  value.status,
								error:   false,
								explain: "",
							}
					} catch (e) {
						console.error("[api]: error occured", e)
						return {
							value: text,
							error: i18n("bad_response"),
							explain: i18n("cannot_parse_response"),
						}
					}
				})
		})
		.catch(e => Promise.reject(e));
}
