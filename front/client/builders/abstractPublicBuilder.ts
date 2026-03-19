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
	messageNavigation  = undefined;
	newComment         = undefined;
	commentsList       = undefined;
	comments           = undefined;

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

	async addMessageNavigation() {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/message_navigation/mobile/message_navigation.mustache' : 'templates/message_navigation/desktop/message_navigation.mustache'
		)), { encoding: 'utf-8' });

		const search = this.search

		this.messageNavigation = mustache.render(template, {
			attachments:      this.i18n("attachments"),
			comments:         this.i18n("comments"),
			attachmentsHref:  function() { const params = new URLSearchParams(search); params.set("msg", "files");     return "?" + params.toString(); },
			commentsHref:     function() { const params = new URLSearchParams(search); params.set("msg", "comments");  return "?" + params.toString(); },
		})
	}

	async addNewComment(message: number | string) {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/comments/mobile/new_comment.mustache' : 'templates/comments/desktop/new_comment.mustache'
		)), { encoding: 'utf-8' });

		this.newComment = mustache.render(template, {
			commentPlaceholder:       this.i18n("commentPlaceholder"),
			newComment:               this.i18n("newComment"),
			redirectUrl:              this.path + this.search,
			message:                  message,
			sendAction:               "/comment/send" + this.search,
		})
	}

	async addCommentsList(comments: Comment[]) {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/comments/mobile/comments_list.mustache' : 'templates/comments/desktop/comments_list.mustache'
		)), { encoding: 'utf-8' });

		this.commentsList = mustache.render(template, {
			comments: comments,
		})
	}

	async addComments(message: number | string, comments: Comment[]) {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/comments/mobile/comments.mustache' : 'templates/comments/desktop/comments.mustache'
		)), { encoding: 'utf-8' });

		await this.addNewComment(message)
		await this.addCommentsList(comments)

		this.comments = mustache.render(template, {}, {
			commentsList:  this.commentsList,
			newComment:    this.newComment,
		})
	}

	async addFilesView(files: FileWithMime[]) {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/files_view/mobile/files_view.mustache' : 'templates/files_view/desktop/files_view.mustache'
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

	async addTranslationView(translation: Translation) {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/translation/mobile/translation_view.mustache' : 'templates/translation/desktop/translation_view.mustache'
		)), { encoding: 'utf-8' });

		this.translationView = mustache.render(template, {
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
