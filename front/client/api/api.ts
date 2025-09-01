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

function getFullUrl(url: string, isFullUrl: boolean) {
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

function getOptions(props): any {
	const {
		headers,
		body,
		method,
		credentials,
	} = props;

	return {
		headers: new Headers({ ...headers }),
		body: prepareBody(body, method),
		method,
		credentials,
	};
}

export {getFileDownloadUrl, getMessageLinkUrl};

export default function api(url, props: any = {}): Promise<any> {
	const { isFullUrl = false, queryParams } = props;

	let fullUrl = getFullUrl(url, isFullUrl);
	fullUrl = appendQueryParams(fullUrl, new URLSearchParams(queryParams));
	const options = getOptions(props);

	console.log(`[${url}]: api`, "fullUrl", fullUrl, "options", options)

	return fetch(fullUrl, options)
		.then(res => {
			console.log(`[${url}]: res:`, res)

			return res
				.text()
				.then(text => {
					try {
						const value = JSON.parse(text);

						if (is.notEmpty(value.error)) {
							console.error(`[${url}]: failed`, "props:", JSON.stringify(props), "res:", text)

							return Promise.resolve([null, {
								error:   true,
								status:  res.status,
								explain: text,
								code:    value.error.code,
							}])
						} else {
							console.log(`[${url}]: success`, "props:", JSON.stringify(props), "res:", text)

							return Promise.resolve([value.response, {
								error:   false,
								status:  res.status,
								explain: "",
								code:    994,
							}])
						}
					} catch (e) {
						console.error(`[${url}]: error:`, e.toString(), "props:", JSON.stringify(props))

						return Promise.resolve([null, {
							error:   true,
							status:  res.status,
							explain: text,
							code:    993,
						}])
					}
				})
				.catch(e => {
					console.error(`[${url}]: deserialize response error:`, e.toString())

					return Promise.reject([null, {
						error:   true,
						status:  500,
						code:    992,
						explain: e.toString(),
					}])
				})
			}
		)
		.catch(e => {
			console.error(`[${url}]: request error`, e.toString())

			return Promise.reject([null, {
				error:   true,
				status:  500,
				code:    991,
				explain: e.toString(),
			}])
		});
}
