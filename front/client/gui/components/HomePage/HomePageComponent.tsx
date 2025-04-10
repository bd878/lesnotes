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

			<Tag css="flex flex-row grow mt-2 max-h-full pb-8">
				<Tag css="flex flex-col items-start w-md w-full">
					<Button
						tabIndex="0"
						content={i18n("load_more")}
						css="disabled:opacity-30 btn w-full text-center mb-5"
						onClick={onLoadMoreClick}
						disabled={isAllLoaded}
					/>

					<Tag
						el="div"
						ref={listRef}
						css="grow w-full h-full overflow-x-hidden overflow-y-scroll"
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

					<Tag css="w-full mt-5">
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
