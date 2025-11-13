import * as is from '../../../third_party/is';
import api from '../../../api';

async function onFontSizeSettingsClick(_elems, e) {
	if (is.notEmpty(e.target.dataset.fontSize)) {
		const response = await api.changeFontSize(e.target.dataset.fontSize)
		if (response.error.error) {
			console.error("[onFontSizeSettingsClick]: cannot update font size:", response)
			return
		}

		location.reload()
	} else {
		console.error("[onFontSizeSettingsClick]: data-font-size is empty")
	}
}

export default onFontSizeSettingsClick
