import React, {lazy} from 'react'
import i18n from '../../../i18n';
import Tag from '../../components/Tag';

const Form = lazy(() => import("../../components/Form"));
const FormField = lazy(() => import("../../components/FormField"));
const Button = lazy(() => import("../../components/Button"));

function LoginFormComponent(props) {
	const {
		name,
		onNameChange,
		password,
		onPasswordChange,
		sendLoginRequest,
	} = props

	return (
		<>
			<Form
				autoComplete="off"
				name="login-form"
			>
				<FormField
					required
					el="input"
					name="name"
					type="text"
					value={name}
					onChange={onNameChange}
				/>
				<FormField
					required
					el="input"
					name="password"
					type="password"
					value={password}
					onChange={onPasswordChange}
				/>
			</Form>

			<Button
				type="button"
				text={i18n("login")}
				onClick={sendLoginRequest}
			/>

			<Tag
				el="a"
				href="/register"
				target="_self"
			>
				{i18n("register")}
			</Tag>
		</>
	)
}

export default LoginFormComponent;
