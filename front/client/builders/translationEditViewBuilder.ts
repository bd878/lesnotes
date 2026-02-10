import type { Message } from '../api/models';
import type { SelectedFile } from '../types';
import Config from 'config';
import mustache from 'mustache';
import api from '../api';
import * as is from '../third_party/is';
import { readFile } from 'node:fs/promises';
import { resolve, join } from 'node:path';
import HomeBuilder from './homeBuilder';

class TranslationEditViewBuilder extends HomeBuilder {
	async addTranslationEditForm() {
		const template = await readFile(resolve(join(Config.get('basedir'), 
			this.isMobile ? 'templates/home/mobile/translation_edit_form.mustache' : 'templates/home/desktop/translation_edit_form.mustache'
		)), { encoding: 'utf-8' });

		this.translationEditForm = mustache.render(template, {})
	}
}

export default TranslationEditViewBuilder
