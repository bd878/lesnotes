import React, {useEffect} from 'react';
import ReactDOM from 'react-dom/client';
import {connect} from '../../../third_party/react-redux';
import {useRawInitData} from '@telegram-apps/sdk-react';
import Tag from '../../components/Tag';
import i18n from '../../../i18n';
import {
	selectIsLoading,
	selectIsValid,
	validateInitDataActionCreator,
} from '../../features/miniapp';

function MiniappMainPageComponent(props) {
	const {
		loading,
		valid,
		validate,
	} = props

	const initData = useRawInitData()

	useEffect(() => {
		validate(initData)
	}, [initData])

	return (
		<Tag css="wrap dark">
			<Tag>{loading ? "loading..." : "loaded"}</Tag>
			<Tag>{valid ? "valid" : "invalid"}</Tag>
			<Tag css="bg-white">{"Miniapp"}</Tag>
		</Tag>
	);
}

const mapStateToProps = state => ({
	loading: selectIsLoading(state),
	valid: selectIsValid(state),
})

const mapDispatchToProps = ({
	validate: validateInitDataActionCreator,
})

export default connect(mapStateToProps, mapDispatchToProps)(MiniappMainPageComponent);
