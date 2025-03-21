import React, {lazy} from 'react';
import Tag from '../Tag';
import i18n from '../../i18n';
import {getFileDownloadUrl} from "../../api";

const List = lazy(() => import("../../components/List"));
const ListItem = lazy(() => import("../../components/ListItem"));

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
					<Tag
						el="li"
						css={liCss + " " + "br-12 p-8 bg-list-grey hover:bg-list-grey pointer"}
						key={`tag_${message.ID}`}
						onClick={(e) => {e.stopPropagation(); onListItemClick(message)}} 
					>
						{message.createUTCNano ? (
							<Tag>
								<Tag el="div">{message.createUTCNano}</Tag>
								<Tag el="div">{message.updateUTCNano}</Tag>
							</Tag>
						) : null}
						{(message.file && message.file.ID && message.file.name) ? <Tag
							el="a"
							href={getFileDownloadUrl(`/files/v1/${message.file.ID}`, false)}
							download={message.file.name}
							target="_blank"
						>
							{message.file.name}
						</Tag> : null}

						<ListItem key={`item_${message.ID}`}>{message.text}</ListItem>
						<Tag
							el="button"
							css="pointer"
							type="button"
							onClick={(e) => {e.stopPropagation();onDeleteClick(message)}}
						>{i18n("delete_message")}</Tag>
						<Tag
							el="button"
							css="pointer"
							type="button"
							onClick={(e) => {e.stopPropagation();onEditClick(message)}}
						>{i18n("edit_message")}</Tag>
					</Tag>
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
