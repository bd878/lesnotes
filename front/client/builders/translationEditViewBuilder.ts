import type { Message, Translation } from '../api/models';
import type { SelectedFile } from '../types';
import Config from 'config';
import mustache from 'mustache';
import api from '../api';
import * as is from '../third_party/is';
import { readFile } from 'node:fs/promises';
import { resolve, join } from 'node:path';
import HomeBuilder from './homeBuilder';

class TranslationEditViewBuilder extends HomeBuilder {
	async addTranslationEditForm(messageID: number, translation: Translation) {
		const template = await readFile(resolve(join(Config.get('basedir'), 
			this.isMobile ? 'templates/home/mobile/translation_edit_form.mustache' : 'templates/home/desktop/translation_edit_form.mustache'
		)), { encoding: 'utf-8' });

		this.translationEditForm = mustache.render(template, {
			message:            messageID,
			translation:        translation,
			titlePlaceholder:   this.i18n("titlePlaceholder"),
			textPlaceholder:    this.i18n("textPlaceholder"),
			cancelEditHref:     `/messages/${messageID}/${translation.lang}` + this.search,
			updateButton:       this.i18n("updateButton"),
			cancelButton:       this.i18n("cancelButton"),
			updateAction:       "/translation/update" + this.search,
			domain:             Config.get("domain"),
		})
	}
}

export default TranslationEditViewBuilder
