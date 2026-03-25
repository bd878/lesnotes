import type { Message, Translation } from '../api/models';
import type { SelectedFile } from '../types';
import Config from 'config';
import mustache from 'mustache';
import api from '../api';
import * as is from '../third_party/is';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import HomeBuilder from './homeBuilder';

let translationEditViewTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/home/desktop/translation_edit_form.mustache')), { encoding: 'utf-8' });
let translationEditViewTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/home/mobile/translation_edit_form.mustache')), { encoding: 'utf-8' });

class TranslationEditViewBuilder extends HomeBuilder {
	addTranslationEditForm(messageID: number, translation: Translation) {
		this.translationEditForm = mustache.render(this.isMobile ? translationEditViewTemplateMobile : translationEditViewTemplate, {
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
