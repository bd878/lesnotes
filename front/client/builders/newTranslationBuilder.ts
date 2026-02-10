import Config from 'config';
import mustache from 'mustache';
import * as is from '../third_party/is';
import { readFile } from 'node:fs/promises';
import { resolve, join } from 'node:path';
import HomeBuilder from './homeBuilder';

class NewTranslationBuilder extends HomeBuilder {
	async addNewTranslationForm() {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/home/mobile/new_translation_form.mustache' : 'templates/home/desktop/new_translation_form.mustache'
		)), { encoding: 'utf-8' });

		this.newTranslationForm = mustache.render(template, {}, {})
	}
}

export default NewTranslationBuilder
