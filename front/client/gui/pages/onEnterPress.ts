function onEnterPress(cbk) {
	return function (e) {
		if (e.keyCode == 13) {
			cbk(e)
		}
	}
}

export default onEnterPress