import type {Builder} from './builder'
import type {TranslationPreview, Translation} from '../api/models';
import Config from 'config';
import mustache from 'mustache';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import AbstractBuilder from './abstractBuilder';

let translationsTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/translations/desktop/translations.mustache')), { encoding: 'utf-8' });
let translationsTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/translations/mobile/translations.mustache')), { encoding: 'utf-8' });

let translationsListTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/translations/desktop/translations_list.mustache')), { encoding: 'utf-8' });
let translationsListTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/translations/mobile/translations_list.mustache')), { encoding: 'utf-8' });

let addTranslationTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/translations/desktop/add_translation.mustache')), { encoding: 'utf-8' });
let addTranslationTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/translations/mobile/add_translation.mustache')), { encoding: 'utf-8' });

let newTranslationTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/translations/desktop/new_translation_form.mustache')), { encoding: 'utf-8' });
let newTranslationTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/translations/mobile/new_translation_form.mustache')), { encoding: 'utf-8' });

let translationViewTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/translations/desktop/translation_view.mustache')), { encoding: 'utf-8' });
let translationViewTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/translations/mobile/translation_view.mustache')), { encoding: 'utf-8' });

let translationEditViewTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/translations/desktop/translation_edit_form.mustache')), { encoding: 'utf-8' });
let translationEditViewTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/translations/mobile/translation_edit_form.mustache')), { encoding: 'utf-8' });

class TranslationsBuilder extends AbstractBuilder {
	translationsList = undefined
	addTranslation = undefined
	newTranslationForm = undefined
	translationEditForm = undefined
	translationView = undefined

	addTranslationsList(message: number, previews: TranslationPreview[]) {
		const search = this.search
		const path = this.path
		this.translationsList = mustache.render(this.isMobile ? translationsListTemplateMobile : translationsListTemplate, {
			newTranslation:        this.i18n("newTranslation"),
			mainMessage:           this.i18n("mainMessage"),
			mainMessageHref:       `/messages/${message}` + this.search,
			newTranslationHref:    function() { const params = new URLSearchParams(search); params.set("trans", "new"); return path + "?" + params.toString() },
			translationHref:       function() { const params = new URLSearchParams(search); params.set("trans", this.lang + ",view"); return path + "?" + params.toString() },
			translations:          previews,
		})
		return this
	}

	addNewTranslation(message: number | string) {
		const search = this.search
		const path = this.path
		this.addTranslation = mustache.render(this.isMobile ? addTranslationTemplateMobile : addTranslationTemplate, {
			newTranslation:        this.i18n("newTranslation"),
			newTranslationHref:    function() { const params = new URLSearchParams(search); params.set("trans", "new"); return path + "?" + params.toString() },
		})
		return this
	}

	addTranslationEditForm(messageID: number, translation: Translation) {
		const search = this.search
		const path = this.path
		this.translationEditForm = mustache.render(this.isMobile ? translationEditViewTemplateMobile : translationEditViewTemplate, {
			message:            messageID,
			translation:        translation,
			titlePlaceholder:   this.i18n("titlePlaceholder"),
			textPlaceholder:    this.i18n("textPlaceholder"),
			redirectUrl:        function() { const params = new URLSearchParams(search); params.set("trans", this.lang + ",view"); return path + "?" + params.toString() },
			cancelEditHref:     function() { const params = new URLSearchParams(search); params.set("trans", this.lang + ",view"); return path + "?" + params.toString() },
			updateButton:       this.i18n("updateButton"),
			cancelButton:       this.i18n("cancelButton"),
			updateAction:       "/translation/update" + this.search,
			domain:             Config.get("domain"),
		})
		return this
	}

	addTranslationView(messageID: number, translation: Translation) {
		const search = this.search
		const path = this.path
		this.translationView = mustache.render(this.isMobile ? translationViewTemplateMobile : translationViewTemplate, {
			messageID:        messageID,
			translation:      translation,
			redirectUrl:      function() { const params = new URLSearchParams(search); params.delete("trans"); return path + "?" + params.toString(); },
			editHref:         function() { const params = new URLSearchParams(search); params.set("trans", this.lang + ",edit"); return path + "?" + params.toString() },
			deleteAction:     "/translation/delete" + search,
		})
		return this
	}

	addTranslationNewForm(messageID: number) {
		const search = this.search
		const path = this.path
		this.newTranslationForm = mustache.render(this.isMobile ? newTranslationTemplateMobile : newTranslationTemplate, {
			titlePlaceholder:    this.i18n("titlePlaceholder"),
			textPlaceholder:     this.i18n("translationPlaceholder"),
			defaultLang:         this.i18n("defaultLang"),
			sendButton:          this.i18n("sendButton"),
			cancelButton:        this.i18n("cancelButton"),
			sendAction:          "/translation/send" + this.search,
			redirectUrl:         function() { const params = new URLSearchParams(search); params.delete("trans"); return path + "?" + params.toString(); },
			cancelHref:          function() { const params = new URLSearchParams(search); params.delete("trans"); return path + "?" + params.toString(); },
			messageID:           messageID,
			langs:               [{lang: "de", label: this.i18n("deLang")}, {lang: "en", label: this.i18n("enLang")}, {lang: "fr", label: this.i18n("frLang")}, {lang: "ru", label: this.i18n("ruLang")}],
		})
		return this
	}

	build() {
		const hasContent = this.translationEditForm || this.newTranslationForm || this.translationView
		return mustache.render(this.isMobile ? translationsTemplateMobile : translationsTemplate, {
			hasContent: hasContent,
		}, {
			translationsList:    this.translationsList,
			addTranslation:      this.addTranslation,
			translationEditForm: this.translationEditForm,
			newTranslationForm:  this.newTranslationForm,
			translationView:     this.translationView,
		})
	}
}

export default TranslationsBuilder
