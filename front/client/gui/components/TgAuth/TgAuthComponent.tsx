import React from 'react'

function TgAuthComponent() {
	return (
		<script async src="https://telegram.org/js/telegram-widget.js?22" dataTelegramLogin={`${BOT_USERNAME}`} dataSize="small" dataAuthUrl={`https://${BACKEND_URL}/tg_auth"`} dataRequestAccess="write"></script>
	)
}

export default TgAuthComponent;
