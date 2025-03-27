import React, {lazy} from 'react';
import Form from '../Form';
import i18n from '../../i18n';

const Form = lazy(() => import("../../components/Form"));
const FormField = lazy(() => import("../../components/FormField"));
const Button = lazy(() => import("../../components/Button"));

function MessageFormComponent(props) {
	const {
		fileRef,
		text,
		onMessageChange,
		onFileChange,
		onSubmit,
		shouldShowCancelButton,
		onCancel,
	} = props

	return (
		<>
			<Form
				autoComplete="off"
				css="row items-end"
			>
				<FormField
					required
					el="textarea"
					name="message"
					type="input"
					value={text}
					onChange={onMessageChange}
				/>
				<FormField
					ref={fileRef}
					el="input"
					name="file"
					type="file"
					onChange={onFileChange}
				/>
			</Form>

			<Button
				type="button"
				text={i18n("msg_send_text")}
				onClick={onSubmit}
			/>
			{(shouldShowCancelButton && onCancel) ? (
				<Button
					type="button"
					text={i18n("msg_cancel_text")}
					onClick={onCancel}
				/>
			) : null}
		</>
	);
}

export default MessageFormComponent;
