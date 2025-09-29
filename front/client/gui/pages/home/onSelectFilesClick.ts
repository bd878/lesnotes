import * as is from '../../../third_party/is';

function onSelectFilesClick(elems, e) {
	if (is.notEmpty(elems.filesInputElem.id)) {
		elems.filesInputElem.click()
	}
}

export default onSelectFilesClick
