import Config from 'config';
import mustache from 'mustache';
import * as is from '../third_party/is';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import HomeBuilder from './homeBuilder';

let newTranslationFormTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/home/desktop/new_translation_form.mustache')), { encoding: 'utf-8' });
let newTranslationFormTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/home/mobile/new_translation_form.mustache')), { encoding: 'utf-8' });

class NewTranslationBuilder extends HomeBuilder {
	addNewTranslationForm(messageID: number) {
		this.newTranslationForm = mustache.render(this.isMobile ? newTranslationFormTemplateMobile : newTranslationFormTemplate, {
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
