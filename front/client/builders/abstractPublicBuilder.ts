import type { Message, TranslationPreview, Translation, ThreadMessages } from '../api/models';
import type { FileWithMime } from '../types';
import Config from 'config';
import mustache from 'mustache';
import api from '../api';
import { unwrapPaging } from '../api/models/paging';
import * as is from '../third_party/is';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import AbstractBuilder from './abstractBuilder'

let signupTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/sidebar_vertical/desktop/signup.mustache')), { encoding: 'utf-8' });
let signupTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/sidebar_vertical/mobile/signup.mustache')), { encoding: 'utf-8' });

let logoutTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/sidebar_vertical/desktop/logout.mustache')), { encoding: 'utf-8' });
let logoutTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/sidebar_vertical/mobile/logout.mustache')), { encoding: 'utf-8' });

let messageNavigationTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/message_navigation/desktop/message_navigation.mustache')), { encoding: 'utf-8' });
let messageNavigationTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/message_navigation/mobile/message_navigation.mustache')), { encoding: 'utf-8' });

let newCommentTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/comments/desktop/new_comment.mustache')), { encoding: 'utf-8' });
let newCommentTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/comments/mobile/new_comment.mustache')), { encoding: 'utf-8' });

let commentsListTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/comments/desktop/comments_list.mustache')), { encoding: 'utf-8' });
let commentsListTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/comments/mobile/comments_list.mustache')), { encoding: 'utf-8' });

let commentsTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/comments/desktop/comments.mustache')), { encoding: 'utf-8' });
let commentsTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/comments/mobile/comments.mustache')), { encoding: 'utf-8' });

let filesViewTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/files_view/desktop/files_view.mustache')), { encoding: 'utf-8' });
let filesViewTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/files_view/mobile/files_view.mustache')), { encoding: 'utf-8' });

let messagesListTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/thread/desktop/messages_list.mustache')), { encoding: 'utf-8' });
let messagesListTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/thread/mobile/messages_list.mustache')), { encoding: 'utf-8' });

let messageViewTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/message/desktop/message_view.mustache')), { encoding: 'utf-8' });
let messageViewTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/message/mobile/message_view.mustache')), { encoding: 'utf-8' });

let translationViewTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/translation/desktop/translation_view.mustache')), { encoding: 'utf-8' });
let translationViewTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/translation/mobile/translation_view.mustache')), { encoding: 'utf-8' });

let searchTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/search/desktop/search_form.mustache')), { encoding: 'utf-8' });
let searchTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/search/mobile/search_form.mustache')), { encoding: 'utf-8' });

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

	addSignup() {
		const search = this.search

		this.signup = mustache.render(this.isMobile ? signupTemplateMobile : signupTemplate, {
			signup:           this.i18n("signup"),
			signupHref:       function() { const params = new URLSearchParams(search); params.delete("cwd"); params.delete("id"); /* TODO: delete pagination */ return "/signup?" + params.toString() },
		})
	}

	addLogout() {
		const search = this.search

		this.logout = mustache.render(this.isMobile ? logoutTemplateMobile : logoutTemplate, {
			logout:           this.i18n("logout"),
			logoutHref:       function() { const params = new URLSearchParams(search); params.delete("cwd"); params.delete("id"); /* TODO: delete pagination */ return "/logout?" + params.toString() },
		})
	}

	addMessageNavigation() {
		const search = this.search

		this.messageNavigation = mustache.render(this.isMobile ? messageNavigationTemplateMobile : messageNavigationTemplate, {
			attachments:      this.i18n("attachments"),
			comments:         this.i18n("comments"),
			attachmentsHref:  function() { const params = new URLSearchParams(search); params.set("msg", "files");     return "?" + params.toString(); },
			commentsHref:     function() { const params = new URLSearchParams(search); params.set("msg", "comments");  return "?" + params.toString(); },
		})
	}

	addNewComment(message: number | string) {
		this.newComment = mustache.render(this.isMobile ? newCommentTemplateMobile : newCommentTemplate, {
			commentPlaceholder:       this.i18n("commentPlaceholder"),
			newComment:               this.i18n("newComment"),
			redirectUrl:              this.path + this.search,
			message:                  message,
			sendAction:               "/comment/send" + this.search,
		})
	}

	addCommentsList(comments: Comment[]) {
		this.commentsList = mustache.render(this.isMobile ? commentsListTemplateMobile : commentsListTemplate, {
			comments: comments,
		})
	}

	addComments(message: number | string, comments: Comment[]) {
		this.addNewComment(message)
		this.addCommentsList(comments)

		this.comments = mustache.render(this.isMobile ? commentsTemplateMobile : commentsTemplate, {}, {
			commentsList:  this.commentsList,
			newComment:    this.newComment,
		})
	}

	addFilesView(files: FileWithMime[]) {
		this.filesView = mustache.render(this.isMobile ? filesViewTemplateMobile : filesViewTemplate, {
			files:    files,
			imgSrc:   function() { return `/files/v1/read/${this.name}` },
			fileHref: function() { return `/files/v1/download?id=${this.ID}` },
		})
	}

	addMessagesList(name: string, messages: ThreadMessages) {
		const search = this.search
		const path = this.path

		const limit = parseInt(LIMIT)

		this.messagesList = mustache.render(this.isMobile ? messagesListTemplateMobile : messagesListTemplate, {
			messages:         unwrapPaging(messages),
			limit:            limit,
			isSingle:         () => messages.messages.length == 1,
			messageHref:      function() { return `/t/${name}/${this.name}` + search; },
			prevPageHref:     function() { const params = new URLSearchParams(search); params.set(this.ID || 0, `${limit + this.offset}`); return path + "?" + params.toString(); },
			nextPageHref:     function() { const params = new URLSearchParams(search); params.set(this.ID || 0, `${Math.max(0, this.offset - limit)}`); return path + "?" + params.toString(); },
		})
	}

	addMessageView(message: Message) {
		this.messageView = mustache.render(this.isMobile ? messageViewTemplateMobile : messageViewTemplate, {
			message: message,
		})
	}

	addTranslationView(translation: Translation) {
		this.translationView = mustache.render(this.isMobile ? translationViewTemplateMobile : translationViewTemplate, {
			translation:      translation,
		})
	}

	addSearch() {
		this.searchForm = mustache.render(this.isMobile ? searchTemplateMobile : searchTemplate, {
			action:              "/search" + this.search,
			searchPlaceholder:   this.i18n("searchPlaceholder"),
			searchMessages:      this.i18n("search"),
		})
	}
}

export default AbstractPublicBuilder;
