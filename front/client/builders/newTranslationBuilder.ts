import Config from 'config';
import mustache from 'mustache';
import * as is from '../third_party/is';
import { readFile } from 'node:fs/promises';
import { resolve, join } from 'node:path';
import HomeBuilder from './homeBuilder';

class NewTranslationBuilder extends HomeBuilder {
	async addNewTranslationForm(messageID: number) {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/home/mobile/new_translation_form.mustache' : 'templates/home/desktop/new_translation_form.mustache'
		)), { encoding: 'utf-8' });

		this.newTranslationForm = mustache.render(template, {
			titlePlaceholder:    this.i18n("titlePlaceholder"),
			textPlaceholder:     this.i18n("textPlaceholder"),
			defaultLang:         this.i18n("defaultLang"),
			sendButton:          this.i18n("sendButton"),
			sendAction:          "/translation/send" + this.search,
			messageID:           messageID,
			langs:               [{lang: "de", label: this.i18n("deLang")}, {lang: "en", label: this.i18n("enLang")}, {lang: "fr", label: this.i18n("frLang")}, {lang: "ru", label: this.i18n("ruLang")}],
		}, {})
	}
}

export default NewTranslationBuilder
