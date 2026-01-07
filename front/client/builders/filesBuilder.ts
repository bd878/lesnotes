import type { File, Paging } from '../api/models';
import Config from 'config';
import mustache from 'mustache';
import * as is from '../third_party/is';
import { readFile } from 'node:fs/promises';
import { resolve, join } from 'node:path';
import HomeBuilder from './homeBuilder';

class FilesBuilder extends HomeBuilder {
	async addFilesList(files: File[] = []) {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/files/mobile/files_list.mustache' : 'templates/files/desktop/files_list.mustache'
		)), { encoding: 'utf-8' });

		this.filesList = mustache.render(template, {
			noFiles:            this.i18n("noFiles"),
			files:              files,
		})
	}

	async addPagination(paging: Paging) {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/files/mobile/pagination.mustache' : 'templates/files/desktop/pagination.mustache'
		)), { encoding: 'utf-8' });

		const search = this.search
		const path = this.path

		const limit = parseInt(LIMIT)

		this.pagination = mustache.render(template, {
			prevPageHref:     function() { const params = new URLSearchParams(search); params.set("files", `${limit + this.offset}`); return path + "?" + params.toString(); },
			nextPageHref:     function() { const params = new URLSearchParams(search); params.set("files", `${Math.max(0, this.offset - limit)}`); return path + "?" + params.toString(); },
			isLastPage:       paging.isLastPage,
			isFirstPage:      paging.isFirstPage,
		})
	}

	async addFilesForm(files: File[] = []) {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/files/mobile/files_form.mustache' : 'templates/files/desktop/files_form.mustache'
		)), { encoding: 'utf-8' });

		this.filesForm = mustache.render(template, {
			noFiles:            this.i18n("noFiles"),
			files:              files,
			sendButton:         this.i18n("sendButton"),
			sendAction:         `/files/v1/upload`,
		})
	}
}

export default FilesBuilder
