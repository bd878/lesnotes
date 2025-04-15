import React, {lazy} from 'react';
import Tag from '../Tag';
import * as is from '../../../third_party/is';

const Thread = lazy(() => import("../Thread"));

function HomePageComponent(props) {
	const {
		stack,
	} = props;

	return (
		<Tag css="flex flex-row grow max-h-full pb-8">
			{stack.map((elem, i) => (
				<Thread key={i} thread={elem} index={i} />
			))}
		</Tag>
	)
}

export default HomePageComponent;
