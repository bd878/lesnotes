import type { File } from '../api/models';
import type { SelectedFile } from '../types';
import * as is from '../third_party/is';

async function selectMessageFiles(ctx, next) {
	console.log("--> selectMessageFiles")

	if (is.array(ctx.state.files.files) && is.array(ctx.state.message.files)) {
		ctx.state.files.files = ctx.state.files.files.map(searchSelected(ctx.state.message.files))
	}

	await next()

	console.log("<-- selectMessageFiles")
}

function searchSelected(selectedFiles: File[]): (file: File) => SelectedFile {
	return function(file: File): SelectedFile {
		let selected: SelectedFile = { ...file, isSelected: false }
		for (let i = 0; i < selectedFiles.length; i++) {
			if (selectedFiles[i].ID == file.ID) {
				selected.isSelected = true
				break
			}
		}
		return selected
	}
}

export default selectMessageFiles
