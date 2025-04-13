import React, {lazy} from 'react';
import Tag from '../Tag';
import i18n from '../../../i18n';

const List = lazy(() => import("../../components/List"));
const MessageElement = lazy(() => import("../../components/MessageElement"))

function MessagesListComponent(props) {
	const {
		css,
		liCss,
		messages,
		loading,
		error,
		onOpenThreadClick,
		onEditClick,
		onDeleteClick,
		onCopyClick,
	} = props

	return (
		<>
			{loading ? i18n("loading") : null}
			{error ? null : (
				<List el="ul" css={css}>
					{messages.map(message => (
						<MessageElement
							tabIndex="0"
							key={message.ID}
							css={liCss}
							message={message}
							onOpenThreadClick={onOpenThreadClick}
							onCopyClick={() => onCopyClick(message)}
							onEditClick={onEditClick}
							onDeleteClick={onDeleteClick}
						/>
					))}
				</List>
			)}
		</>
	)
}

export default MessagesListComponent;
