import * as is from '../../../third_party/is';

async function onSearchFormSubmit(elems, e) {
	e.preventDefault()

	if (is.notEmpty(elems.searchFormElem.search.value)) {
		const query = elems.searchFormElem.search.value

		location.href = "/search?query=" + query
	} else {
		console.error("[onSearchFormSubmit]: search value is empty")
	}
}

export default onSearchFormSubmit
