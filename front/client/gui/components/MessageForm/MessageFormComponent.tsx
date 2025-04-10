import React, {lazy} from 'react';
import Form from '../Form';
import Tag from '../Tag';
import AttachIcon from '../../icons/AttachIcon';
import CheckmarkIcon from '../../icons/CheckmarkIcon';
import CrossIcon from '../../icons/CrossIcon';
import i18n from '../../../i18n';

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
		<Form autoComplete="off" css="flex flex-row justify-items-end">
			<Tag css="w-full border rounded-sm focus-within:outline focus-within:outline-offset-2">
				<FormField
					required
					css="focus:outline-none w-full resize-none px-2 py-1"
					el="textarea"
					name="message"
					type="input"
					value={text}
					onChange={onMessageChange}
				/>
				<Tag css="p-2 flex flex-row justify-between">
					<Tag css="flex">
						<Tag tabIndex="0" css="cursor-pointer inline-block" el="label" htmlFor="file">
							<AttachIcon width="24" height="24" />
						</Tag>
						<FormField
							ref={fileRef}
							css="hidden"
							el="input"
							id="file"
							name="file"
							type="file"
							onChange={onFileChange}
						/>
					</Tag>

					<Tag css="flex">
						{shouldShowCancelButton && onCancel ? (
							<Button
								type="button"
								content={
									<CrossIcon width="24" height="24" />
								}
								css="mr-2 cursor-pointer"
								onClick={onCancel}
							/>
						) : null}

						<Button
							type="button"
							content={
								<CheckmarkIcon width="24" height="24" />
							}
							css="cursor-pointer"
							onClick={onSubmit}
						/>
					</Tag>
				</Tag>
			</Tag>
		</Form>
	);
}

export default MessageFormComponent;
