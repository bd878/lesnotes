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
				css="flex flex-col"
			>
				<Tag>
					<Tag el="label" htmlFor="name">{i18n("username")}</Tag>
					<FormField
						required
						el="input"
						name="name"
						type="text"
						css="block w-full border-solid border-1"
						value={name}
						onChange={onNameChange}
					/>
				</Tag>
				<Tag css="mt-2">
					<Tag el="label" htmlFor="password">{i18n("password")}</Tag>
					<FormField
						required
						el="input"
						name="password"
						type="password"
						css="block w-full border-solid border-1"
						value={password}
						onChange={onPasswordChange}
					/>
				</Tag>
			</Form>

			<Tag css="flex flex-row justify-between mt-3">
				<Tag
					el="a"
					href="/register"
					target="_self"
					css="underline italic text-blue-600 text-left cursor-pointer"
				>
					{i18n("register")}
				</Tag>

				<Button
					type="button"
					text={i18n("login") + " >"}
					css="border-1 px-4 border-solid text-left cursor-pointer"
					onClick={sendLoginRequest}
				/>
			</Tag>
		</>
	)
}

export default LoginFormComponent;
