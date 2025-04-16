import React, {lazy} from 'react';
import Tag from '../Tag';
import * as is from '../../../third_party/is';

const Thread = lazy(() => import("../Thread"));

function HomePageComponent(props) {
	const {
		stack,
		openThread,
		closeThread,
	} = props;

	return (
		<Tag css="flex flex-row grow max-h-full pb-8">
			{stack.map((elem, index) => (
				<Thread
					key={elem.ID}
					thread={elem}
					index={index}
					openThread={openThread(index)}
					closeThread={closeThread(index)}
				/>
			))}
		</Tag>
	)
}

export default HomePageComponent;
