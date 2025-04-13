import React, {lazy} from 'react';
import Form from '../Form';
import Tag from '../Tag';
import AttachIcon from '../../icons/AttachIcon';
import CheckmarkIcon from '../../icons/CheckmarkIcon';
import CrossIcon from '../../icons/CrossIcon';
import * as is from '../../../third_party/is'
import i18n from '../../../i18n';

const Form = lazy(() => import("../../components/Form"));
const FormField = lazy(() => import("../../components/FormField"));
const Button = lazy(() => import("../../components/Button"));

function MessageFormComponent(props) {
	const {
		fileRef,
		text,
		file,
		onMessageChange,
		onFileChange,
		onSubmit,
		onFileClick,
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

				<Tag css="py-2 pr-2 pl-1 flex flex-row justify-between">
					<Tag css="group flex overflow-hidden whitespace-nowrap text-ellipsis mr-2">
						<Button
							type="button"
							tabIndex="0"
							content={
								<AttachIcon width="24" height="24" />
							}
							css="m-1 cursor-pointer"
							onClick={onFileClick}
						/>
						<FormField
							ref={fileRef}
							css="hidden"
							el="input"
							id="file"
							name="file"
							type="file"
							onChange={onFileChange}
						/>
						{is.notUndef(file) ? (
							<Tag css="w-full flex items-center text-sm italic overflow-hidden">
								<Tag css="px-2 overflow-hidden whitespace-nowrap text-ellipsis">{file.name}</Tag>
							</Tag>
						) : null}
					</Tag>

					<Tag css="flex items-center">
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
