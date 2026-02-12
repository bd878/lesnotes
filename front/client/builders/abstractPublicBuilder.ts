import type { Message, TranslationPreview, Translation, ThreadMessages } from '../api/models';
import type { FileWithMime } from '../types';
import Config from 'config';
import mustache from 'mustache';
import api from '../api';
import { unwrapPaging } from '../api/models/paging';
import * as is from '../third_party/is';
import { readFile } from 'node:fs/promises';
import { resolve, join } from 'node:path';
import AbstractBuilder from './abstractBuilder'

abstract class AbstractPublicBuilder extends AbstractBuilder {
	signup             = undefined;
	logout             = undefined;
	sidebar            = undefined;
	translations       = undefined;
	filesView          = undefined;
	messageView        = undefined;
	searchForm         = undefined;
	messagesList       = undefined;
	translationView    = undefined;

	async addSignup() {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/sidebar_vertical/mobile/signup.mustache' : 'templates/sidebar_vertical/desktop/signup.mustache'
		)), { encoding: 'utf-8' });

		const search = this.search

		this.signup = mustache.render(template, {
			signup:           this.i18n("signup"),
			signupHref:       function() { const params = new URLSearchParams(search); params.delete("cwd"); params.delete("id"); /* TODO: delete pagination */ return "/signup?" + params.toString() },
		})
	}

	async addLogout() {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/sidebar_vertical/mobile/logout.mustache' : 'templates/sidebar_vertical/desktop/logout.mustache'
		)), { encoding: 'utf-8' });

		const search = this.search

		this.logout = mustache.render(template, {
			logout:           this.i18n("logout"),
			logoutHref:       function() { const params = new URLSearchParams(search); params.delete("cwd"); params.delete("id"); /* TODO: delete pagination */ return "/logout?" + params.toString() },
		})
	}

	async addSidebar() {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/sidebar_vertical/mobile/sidebar_vertical.mustache' : 'templates/sidebar_vertical/desktop/sidebar_vertical.mustache'
		)), { encoding: 'utf-8' });

		this.sidebar = mustache.render(template, {
			settingsHeader: this.i18n("settingsHeader"),
			mainHref:       "/" + this.search,
		}, {
			settings:       this.settings,
			searchForm:     this.searchForm,
		})
	}

	async addTranslations(messageID: number, previews: TranslationPreview[]) {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/translations/mobile/translations.mustache' : 'templates/translations/desktop/translations.mustache'
		)), { encoding: 'utf-8' });

		const search = this.search

		this.translations = mustache.render(template, {
			newTranslation:        this.i18n("newTranslation"),
			mainMessage:           this.i18n("mainMessage"),
			mainMessageHref:       function() { return `/messages/${messageID}` },
			newTranslationHref:    function() { return `/editor/messages/${messageID}/new_lang` },
			translationHref:       function() { return `/messages/${messageID}/${this.lang}` },
			translations:          previews,
			hasTranslations:       () => previews.length > 0,
		})
	}

	async addFilesView(files: FileWithMime[]) {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/message/mobile/files_view.mustache' : 'templates/message/desktop/files_view.mustache'
		)), { encoding: 'utf-8' });

		this.filesView = mustache.render(template, {
			files:    files,
			imgSrc:   function() { return `/files/v1/read/${this.name}` },
			fileHref: function() { return `/files/v1/download?id=${this.ID}` },
		})
	}

	async addMessagesList(name: string, messages: ThreadMessages) {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/thread/mobile/messages_list.mustache' : 'templates/thread/desktop/messages_list.mustache'
		)), { encoding: 'utf-8' });

		const search = this.search
		const path = this.path

		const limit = parseInt(LIMIT)

		this.messagesList = mustache.render(template, {
			messages:         unwrapPaging(messages),
			limit:            limit,
			isSingle:         () => messages.messages.length == 1,
			messageHref:      function() { return `/t/${name}/${this.name}` + search; },
			prevPageHref:     function() { const params = new URLSearchParams(search); params.set(this.ID || 0, `${limit + this.offset}`); return path + "?" + params.toString(); },
			nextPageHref:     function() { const params = new URLSearchParams(search); params.set(this.ID || 0, `${Math.max(0, this.offset - limit)}`); return path + "?" + params.toString(); },
		})
	}

	async addMessageView(message: Message) {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/message/mobile/message_view.mustache' : 'templates/message/desktop/message_view.mustache'
		)), { encoding: 'utf-8' });

		this.messageView = mustache.render(template, {
			message: message,
		})
	}

	async addTranslationView(messageID: number, translation: Translation) {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/translation/mobile/translation_view.mustache' : 'templates/translation/desktop/translation_view.mustache'
		)), { encoding: 'utf-8' });

		this.translationView = mustache.render(template, {
			messageID:        messageID,
			translation:      translation,
		})
	}

	async addSearch() {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/search_form/mobile/search_form.mustache' : 'templates/search_form/desktop/search_form.mustache'
		)), { encoding: 'utf-8' });

		this.searchForm = mustache.render(template, {
			action:              "/search" + this.search,
			searchPlaceholder:   this.i18n("searchPlaceholder"),
			searchMessages:      this.i18n("search"),
		})
	}
}

export default AbstractPublicBuilder;
