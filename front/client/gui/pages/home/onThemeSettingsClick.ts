import * as is from '../../../third_party/is';
import api from '../../../api';

async function onThemeSettingsClick(_elems, e) {
	if (is.notEmpty(e.target.dataset.theme)) {
		const response = await api.changeTheme(e.target.dataset.theme)
		if (response.error.error) {
			console.error("[onThemeSettingsClick]: cannot update theme:", response)
			return
		}

		location.reload()
	} else {
		console.error("[onThemeSettingsClick]: data-theme is empty")
	}
}

export default onThemeSettingsClick
