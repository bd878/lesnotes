import React, {useEffect, useRef} from 'react'

function TgAuthComponent() {
	const ref = useRef(null)

	useEffect(() => {
		if (ref.current == null)
			return

		const scriptEl = document.createElement("script", {
			"async": true,
			"src": "https://telegram.org/js/telegram-widget.js?22",
			"data-telegram-login": `${BOT_USERNAME}`,
			"data-size": "small",
			"data-auth-url": `https://${BACKEND_URL}/tg_auth"`,
			"data-request-access": "write",
		});
		ref.current.append(scriptEl)
	}, [ref])

	return (
		<div id="telegram-login-widget" ref={ref}>
		</div>
	)
}

export default TgAuthComponent;
