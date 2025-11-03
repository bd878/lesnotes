import * as is from '../../../third_party/is';

async function onLangSettingsClick(_elems, e) {
	if (is.notEmpty(e.target.dataset.lang)) {
		const params = new URLSearchParams(location.search)
		params.set("lang", e.target.dataset.lang)
		location.href = params.toString() ? ("/signup?" + params.toString()) : "/signup"
	} else {
		console.error("[onLangSettingsClick]: data-lang is empty")
	}
}

export default onLangSettingsClick
