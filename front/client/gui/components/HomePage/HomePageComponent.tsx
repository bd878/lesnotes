import React, {lazy} from 'react';
import Tag from '../Tag';
import i18n from '../../../i18n';
import * as is from '../../../third_party/is';

const Thread = lazy(() => import("../Thread"));

function HomePageComponent(props) {
	const {
		stack,
		openThread,
		closeThread,
		destroyThread,
	} = props;

	return (
		<Tag css="flex flex-row grow max-h-full pb-8">
			{stack.map((elem, index) => (
				<Thread
					css={index > 0 ? "ml-4" : ""}
					key={elem.ID}
					index={index}
					destroyThread={destroyThread(index)}
					openThread={openThread(index)}
					closeThread={closeThread(index)}
					destroyContent={index === 0 ? "< " + i18n("logout") : ("X " + i18n("close_button_text"))}
				/>
			))}
		</Tag>
	)
}

export default HomePageComponent;
