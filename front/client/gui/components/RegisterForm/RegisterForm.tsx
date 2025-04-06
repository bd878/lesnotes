import React, {lazy, useState, useCallback} from 'react';
import {connect} from '../../../third_party/react-redux';
import {registerActionCreator} from '../../features/me'
import Tag from '../Tag';
import i18n from '../../../i18n';

const Form = lazy(() => import("../../components/Form"));
const FormField = lazy(() => import("../../components/FormField"));
const Button = lazy(() => import("../../components/Button"));

function RegisterForm(props) {
	const {
		register,
	} = props

	const [name, setName] = useState("")
	const [password, setPassword] = useState("")

	const onNameChange = useCallback(e => {
		setName(e.target.value)
	}, [setName])

	const onPasswordChange = useCallback(e => {
		setPassword(e.target.value)
	}, [setPassword])

	const sendRegisterRequest = useCallback(e => {
		e.preventDefault();

		if (!name) {console.error(i18n("name_required_err")); return;}
		if (!password) {console.error(i18n("pass_required_err")); return;}

		register(name, password)
	}, [register, name, password])

	return (
		<Tag>
			<Form autoComplete="off" name="register-form">
				<FormField required el="input" name="name" type="text" value={name} onChange={onNameChange} />
				<FormField required el="input" name="password" type="password" value={password} onChange={onPasswordChange} />
			</Form>
			<Button type="button" text={i18n("register")} onClick={sendRegisterRequest} />
			<Tag el="a" href="/login" target="_self">{i18n("login")}</Tag>
		</Tag>
	);
}

const mapStateToProps = () => {}

const mapDispatchToProps = ({
	register: registerActionCreator,
})

export default connect(mapStateToProps, mapDispatchToProps)(RegisterForm);
