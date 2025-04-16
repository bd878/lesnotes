function equal(val1) {
	let inverted = false
	function equalFold(val2) {
		return inverted
			? val1 !== val2
			: val1 === val2
	}
	equalFold.not = () => {inverted = true; return equalFold;}
	return equalFold
}

export default equal;
