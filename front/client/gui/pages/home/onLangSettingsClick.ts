import * as is from '../../../third_party/is';
import api from '../../../api';

async function onLangSettingsClick(_elems, e) {
	if (is.notEmpty(e.target.dataset.lang)) {
		const response = await api.changeLanguage(e.target.dataset.lang)
		if (response.error.error) {
			console.error("[onLangSettingsClick]: cannot update lang:", response)
			return
		}

		location.reload()
	} else {
		console.error("[onLangSettingsClick]: data-lang is empty")
	}
}

export default onLangSettingsClick
