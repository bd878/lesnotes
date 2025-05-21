import React, {useRef, useCallback} from 'react'
import MessageListElementComponent from './MessageListElementComponent'

function MessageListElementContainer(props) {
	const {
		css,
		message,
		index,
		isMyThreadOpen,
		isSelected,
		isEdit,
		isPublic,
		isAnyThreadOpen,
		onToggleThreadClick,
		onEditClick,
		onCopyClick,
		onCopyLinkClick,
		onSelectClick,
		onUnselectClick,
		onResetEditClick,
	} = props

	const listRef = useRef(null)

	const onDragStart = useCallback(event => {
		if (listRef.current == null)
			return

		let rect = listRef.current.getBoundingClientRect()
		let shiftX = event.clientX - rect.left;
		let shiftY = event.clientY - rect.top;

		event.dataTransfer.effectAllowed = "move"
		event.dataTransfer.setData("text/plain", JSON.stringify([index, message]))
		event.dataTransfer.setDragImage(listRef.current, shiftX, shiftY)
	}, [listRef, message, index])

	return (
		<MessageListElementComponent
			css={css}
			listRef={listRef}
			onDragStart={onDragStart}
			message={message}
			isMyThreadOpen={isMyThreadOpen}
			isSelected={isSelected}
			isEdit={isEdit}
			isPublic={isPublic}
			isAnyThreadOpen={isAnyThreadOpen}
			onToggleThreadClick={onToggleThreadClick}
			onEditClick={onEditClick}
			onCopyClick={onCopyClick}
			onCopyLinkClick={onCopyLinkClick}
			onSelectClick={onSelectClick}
			onUnselectClick={onUnselectClick}
			onResetEditClick={onResetEditClick}
		/>
	)
}

export default MessageListElementContainer