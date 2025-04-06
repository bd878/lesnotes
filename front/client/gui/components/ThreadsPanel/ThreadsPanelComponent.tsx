import React, {lazy} from 'react'
import Tag from '../Tag';
import i18n from '../../../i18n';
import MessagesList from '../MessagesList';
import MainMessage from '../MainMessage';
import MessageForm from '../MessageForm';

const Button = lazy(() => import("../../components/Button"));

function ThreadsPanelComponent(props) {
	const {
		threadMessage,
		messages,
		close,
		send,
	} = props

	return (
		<Tag>
			<Button
				type="button"
				text={i18n("close_button_text")}
				onClick={close}
			/>
			<MainMessage message={threadMessage} />
			<MessagesList messages={messages} 	/>
			<MessageForm messageForEdit={{}} reset={() => {}} send={send} update={() => {}} edit={() => {}} />
		</Tag>
	)
}

export default ThreadsPanelComponent;