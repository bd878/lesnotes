import React, {lazy} from 'react';
import i18n from '../../../i18n';
import Button from '../Button';
import Tag from '../Tag';

const MessagesList = lazy(() => import("../MessagesList"));
const MessageForm = lazy(() => import("../MessageForm"));
const MainMessage = lazy(() => import("../MainMessage"));

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

		shouldShowThreadsPanel,
		threadMessage,
		threadMessages,
		closeThread,
		sendThreadMessage,
	} = props;

	return (
		<>
			<Tag css="flex flex-row grow max-h-full pb-8">
				<Tag css="flex flex-col items-start w-md w-full">
					<Button
						content={"< " + i18n("logout")}
						onClick={onExitClick}
						css="btn mb-2"
					/>

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

				{shouldShowThreadsPanel ? (
					<Tag css="flex flex-col items-start w-md w-full">
						<Button
							type="button"
							content={i18n("close_button_text")}
							onClick={closeThread}
						/>
						<MainMessage message={threadMessage} />
						<MessagesList messages={threadMessages} 	/>
						<MessageForm messageForEdit={{}} reset={() => {}} send={sendThreadMessage} update={() => {}} edit={() => {}} />
					</Tag>
				) : null}
			</Tag>
		</>
	)
}

export default HomePageComponent;
