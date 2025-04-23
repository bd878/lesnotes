function equal(val1) {
	let inverted = false
	function equalFold(val2) {
		return inverted
			? val1 !== val2
			: val1 === val2
	}
	equalFold.not = (val2) => {inverted = true; return equalFold(val2);}
	return equalFold
}

export default equal;
