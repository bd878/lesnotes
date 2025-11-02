function onFileInputChange(elems, e) {
	hideNoFilesListElem(elems)

	for (const file of e.target.files) {
		elems.filesListElem.appendChild(createFilesListElement(elems, file.name))
	}
}

function hideNoFilesListElem(elems) {
	elems.noFilesElem.classList.add(["hidden"])
}

function showNoFilesListElem(elems) {
	elems.noFilesElem.classList.remove("hidden")
}

function createFilesListElement(elems, fileName: string): HTMLDivElement {
	const elem = document.createElement("div")

	const textElem = document.createElement("span")
	const removeButton = document.createElement("button")

	removeButton.textContent = "X"
	removeButton.classList.add(...("cursor-pointer underline dark:hover:text-white hover:text-blue-600 mr-2").split(" "))

	removeButton.onclick = () => {
		elem.remove();

		if (elems.filesListElem.childElementCount == 0) {
			showNoFilesListElem(elems);
		}
	}

	textElem.textContent = fileName
	textElem.classList.add(...("overflow-hidden dark:text-white text-ellipsis").split(" "))

	elem.classList.add(...("mb-2 overflow-hidden text-ellipsis".split(" ")))
	elem.appendChild(removeButton)
	elem.appendChild(textElem)

	return elem
}

export default onFileInputChange;
