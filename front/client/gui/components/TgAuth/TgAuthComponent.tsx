import React from 'react'

function TgAuthComponent() {
	return (
		<script src="https://telegram.org/js/telegram-widget.js?22" data-telegram-login={`${BOT_USERNAME}`} data-size="small" data-auth-url={`https://${BACKEND_URL}/tg_auth"`} data-request-access="write"></script>
	)
}

export default TgAuthComponent;
