import React, {lazy} from 'react';
import Tag from '../Tag';
import i18n from '../../i18n';

const List = lazy(() => import("../../components/List"));
const MessageElement = lasy(() => import("../../components/MessageElement"))

function MessagesListComponent(props) {
	const {
		css,
		liCss,
		messages,
		loading,
		error,
		onListItemClick,
		onEditClick,
		onDeleteClick,
	} = props

	let content = <Tag></Tag>;
	if (error)
		content = <Tag>{error}</Tag>
	else
		content = (
			<List el="ul" css={css}>
				{messages.map(message => (
					<MessageElement
						css={liCss}
						message={message}
						onClick={onListItemClick}
						onEditClick={onEditClick}
						onDeleteClick={onDeleteClick}
					/>
				))}
			</List>
		)

	return (
		<>
			{loading ? <Tag>{i18n("loading")}</Tag> : null}
			{content}
		</>
	)
}

export default MessagesListComponent;
