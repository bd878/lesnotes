import * as is from '../../../third_party/is';

async function onThemeSettingsClick(_elems, e) {
	if (is.notEmpty(e.target.dataset.theme)) {
		const params = new URLSearchParams(location.search)
		params.set("theme", e.target.dataset.theme)
		location.href = params.toString() ? ("/signup?" + params.toString()) : "/signup"
	} else {
		console.error("[onThemeSettingsClick]: data-theme is empty")
	}
}

export default onThemeSettingsClick
