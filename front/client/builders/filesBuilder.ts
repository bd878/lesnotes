import type { File, Paging } from '../api/models';
import Config from 'config';
import mustache from 'mustache';
import * as is from '../third_party/is';
import { readFile } from 'node:fs/promises';
import { resolve, join } from 'node:path';
import HomeBuilder from './homeBuilder';

class FilesBuilder extends HomeBuilder {
	async addFilesList(files: File[] = [], paging: Paging) {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/home/mobile/files_list.mustache' : 'templates/home/desktop/files_list.mustache'
		)), { encoding: 'utf-8' });

		this.filesList = mustache.render(template, {
			noFiles:            this.i18n("noFiles"),
			files:              files,
			paging:             paging,
		})
	}

	async addFilesInput(files: File[] = []) {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/home/mobile/files_input.mustache' : 'templates/home/desktop/files_input.mustache'
		)), { encoding: 'utf-8' });

		this.filesInput = mustache.render(template, {
			noFiles:            this.i18n("noFiles"),
			files:              files,
		})
	}
}

export default FilesBuilder
