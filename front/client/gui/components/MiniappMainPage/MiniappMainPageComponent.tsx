import React, {useEffect} from 'react';
import ReactDOM from 'react-dom/client';
import {connect} from '../../../third_party/react-redux';
import {useRawInitData, useLaunchParams, themeParams} from '@telegram-apps/sdk-react';
import Tag from '../../components/Tag';
import i18n from '../../../i18n';
import api from '../../../api';
import {
	selectIsLoading,
	selectIsValid,
	validateInitDataActionCreator,
} from '../../features/miniapp';
import {selectStack} from '../../features/stack';
import MiniappThread from '../MiniappThread'

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

	useEffect(() => {
		api.sendLog("send launch params")
		api.sendLog(JSON.stringify(launchParams.tgWebAppThemeParams))
	}, [launchParams])

	useEffect(() => {
		console.log("theme params state:", themeParams.state())
	}, [])

	return (
		<Tag css="w-full h-full">
			{valid ? stack.map((elem, index) => (
				<MiniappThread
					key={elem.ID}
					index={index}
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
