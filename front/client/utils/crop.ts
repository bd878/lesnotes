
function crop(str: string, size: number): string {
	if (str.length > size) {
		return `${str.slice(0, size)}...`
	} else if (str.length <= size) {
		return str
	}
}

export default crop
