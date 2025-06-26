import React, {useEffect, lazy} from 'react';
import ReactDOM from 'react-dom/client';
import {connect} from '../../../third_party/react-redux';
import {useRawInitData, useLaunchParams} from '@telegram-apps/sdk-react';
import Tag from '../../components/Tag';
import i18n from '../../../i18n';
import {
	selectIsLoading,
	selectIsValid,
	validateInitDataActionCreator,
} from '../../features/miniapp';
import {selectStack} from '../../features/stack';

const Thread = lazy(() => import("../Thread"));

function MiniappMainPageComponent(props) {
	const {
		stack,
		loading,
		valid,
		validate,
	} = props

	const initData = useRawInitData()
	const launchParams = useLaunchParams()

	useEffect(() => {
		validate(initData)
	}, [initData])

	console.log(launchParams)

	return (
		<Tag css="wrap dark">
			<Tag>{loading ? "loading..." : "loaded"}</Tag>
			<Tag css="bg-white">{"Miniapp"}</Tag>

			{valid ? stack.map((elem, index) => (
				<Thread
					css={index > 0 ? "ml-4" : ""}
					key={elem.ID}
					thread={elem}
					index={index}
					destroyThread={() => () => {}}
					openThread={() => () => {}}
					closeThread={() => () => {}}
					destroyContent={index === 0 ? "< " + i18n("logout") : ("X " + i18n("close_button_text"))}
				/>
			)) : "invalid"}
		</Tag>
	);
}

const mapStateToProps = state => ({
	loading: selectIsLoading(state),
	valid: selectIsValid(state),
	stack: selectStack(state),
})

const mapDispatchToProps = ({
	validate: validateInitDataActionCreator,
})

export default connect(mapStateToProps, mapDispatchToProps)(MiniappMainPageComponent);
