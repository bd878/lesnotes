import * as is from '../../../third_party/is'

function onFilesListClick(elems, e) {
	const fileElem = document.getElementById(e.target.dataset.fileId)
	if (is.notEmpty(fileElem)) {
		fileElem.remove()
	}
}

export default onFilesListClick

