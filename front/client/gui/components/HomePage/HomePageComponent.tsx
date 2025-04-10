import React, {lazy} from 'react';
import i18n from '../../../i18n';
import Button from '../Button';
import Tag from '../Tag';

const MessagesList = lazy(() => import("../MessagesList"));
const MessageForm = lazy(() => import("../MessageForm"));
const ThreadsPanel = lazy(() => import("../ThreadsPanel"));

function HomePageComponent(props) {
	const {
		listRef,
		onExitClick,
		onListScroll,
		onLoadMoreClick,
		isAllLoaded,
		error,
		messages,
		loading,
		sendMessage,
		updateMessage,
		resetEditMessage,
		messageForEdit,
	} = props;

	return (
		<>
			<Button
				content={"< " + i18n("logout")}
				onClick={onExitClick}
				css="btn"
			/>

			<Tag css="flex flex-row grow mt-2">
				<Tag css="flex flex-col items-start w-md w-full">
					<Button
						content={i18n("load_more")}
						css="btn w-full text-center mb-5"
						onClick={onLoadMoreClick}
						disabled={isAllLoaded}
					/>

					<Tag
						el="div"
						ref={listRef}
						css="grow w-full"
						onScroll={onListScroll}
					>
						<MessagesList
							css="w-full"
							liCss="mb-2"
							error={error}
							messages={messages}
							loading={loading}
						/>
					</Tag>

					<Tag css="w-full">
						<MessageForm
							send={sendMessage}
							update={updateMessage}
							reset={resetEditMessage}
							resetEdit={resetEditMessage}
							messageForEdit={messageForEdit}
						/>
					</Tag>
				</Tag>

				<Tag>
					<ThreadsPanel />
				</Tag>
			</Tag>
		</>
	)
}

export default HomePageComponent;
