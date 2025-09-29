function onFileInputChange(elems, e) {
	for (const file of e.target.files) {
		elems.filesListElem.appendChild(createFilesListElement(file.name))
	}
}

function createFilesListElement(fileName: string): HTMLDivElement {
	const elem = document.createElement("div")

	const textElem = document.createElement("span")
	const removeButton = document.createElement("button")

	removeButton.textContent = "X"
	removeButton.classList.add(...("cursor-pointer underline hover:text-blue-600 mr-2").split(" "))

	removeButton.onclick = () => { elem.remove() }

	textElem.textContent = fileName
	textElem.classList.add(...("overflow-hidden text-ellipsis").split(" "))

	elem.classList.add(...("mb-2 overflow-hidden text-ellipsis".split(" ")))
	elem.appendChild(removeButton)
	elem.appendChild(textElem)

	return elem
}

export default onFileInputChange;
