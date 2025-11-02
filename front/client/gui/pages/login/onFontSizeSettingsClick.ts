import * as is from '../../../third_party/is';

async function onFontSizeSettingsClick(_elems, e) {
	if (is.notEmpty(e.target.dataset.fontSize)) {
		const params = new URLSearchParams(location.search)
		params.set("size", e.target.dataset.fontSize)
		location.href = params.toString() ? ("/login?" + params.toString()) : "/login"
	} else {
		console.error("[onFontSizeSettingsClick]: data-font-size is empty")
	}
}

export default onFontSizeSettingsClick
