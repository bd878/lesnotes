import * as is from '../third_party/is';

function crop(str: string, size: number): string {
	if (is.undef(str)) {
		return ""
	}
	if (str.length > size) {
		return `${str.slice(0, size)}...`
	} else if (str.length <= size) {
		return str
	}
}

export default crop
