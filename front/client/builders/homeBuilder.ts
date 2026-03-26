import type {Builder} from './builder';
import type {File, Message, Thread, Comment, MessagesList} from '../api/models';
import type {FileWithMime} from '../types';
import type {TranslationPreview} from '../api/models';
import Config from 'config';
import mustache from 'mustache';
import api from '../api';
import { unwrapPaging } from '../api/models/paging';
import * as is from '../third_party/is';
import i18n from '../i18n';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import AbstractBuilder from './abstractBuilder'

let newCommentTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/comments/desktop/new_comment.mustache')), { encoding: 'utf-8' });
let newCommentTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/comments/mobile/new_comment.mustache')), { encoding: 'utf-8' });

let commentsListTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/comments/desktop/comments_list.mustache')), { encoding: 'utf-8' });
let commentsListTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/comments/mobile/comments_list.mustache')), { encoding: 'utf-8' });

let commentsTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/comments/desktop/comments.mustache')), { encoding: 'utf-8' });
let commentsTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/comments/mobile/comments.mustache')), { encoding: 'utf-8' });

let messageNavigationTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/message_navigation/desktop/message_navigation.mustache')), { encoding: 'utf-8' });
let messageNavigationTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/message_navigation/mobile/message_navigation.mustache')), { encoding: 'utf-8' });

let translationsTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/translations/desktop/translations.mustache')), { encoding: 'utf-8' });
let translationsTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/translations/mobile/translations.mustache')), { encoding: 'utf-8' });

let newTranslationTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/translations/desktop/new_translation.mustache')), { encoding: 'utf-8' });
let newTranslationTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/translations/mobile/new_translation.mustache')), { encoding: 'utf-8' });

let controlPanelTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/home/desktop/control_panel.mustache')), { encoding: 'utf-8' });
let controlPanelTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/home/mobile/control_panel.mustache')), { encoding: 'utf-8' });

let navigationTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/home/desktop/navigation.mustache')), { encoding: 'utf-8' });
let navigationTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/home/mobile/navigation.mustache')), { encoding: 'utf-8' });

let filesSelectorTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/home/desktop/files_selector.mustache')), { encoding: 'utf-8' });
let filesSelectorTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/home/mobile/files_selector.mustache')), { encoding: 'utf-8' });

let filesViewTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/files_view/desktop/files_view.mustache')), { encoding: 'utf-8' });
let filesViewTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/files_view/mobile/files_view.mustache')), { encoding: 'utf-8' });

let homeTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/home/desktop/home.mustache')), { encoding: 'utf-8' });
let homeTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/home/mobile/home.mustache')), { encoding: 'utf-8' });

class HomeBuilder extends AbstractBuilder {
	messageEditForm      = undefined;
	messageView          = undefined;
	newMessageForm       = undefined;
	newTranslationForm   = undefined;
	translationEditForm  = undefined;
	translationView      = undefined;
	threadView           = undefined;
	threadEditForm       = undefined;
	pagination           = undefined;
	filesSelector        = undefined;
	filesForm            = undefined;
	filesView            = undefined;
	filesList            = undefined;
	header               = undefined;
	messagesTree         = undefined;
	sidebar              = undefined;
	controlPanel         = undefined;
	navigation           = undefined;
	translations         = undefined;
	newTranslation       = undefined;
	messageNavigation    = undefined;
	newComment           = undefined;
	commentsList         = undefined;
	comments             = undefined;
	logout               = undefined;
	messageHeader        = undefined;
	scripts              = ["/public/pages/home/homeScript.js"]

	addMessagesTree(tree: Builder) {
		this.messagesTree = tree.build()
	}

	addMessageHeader(header: Builder) {
		this.messageHeader = header.build()
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

	addMessageNavigation() {
		const search = this.search

		this.messageNavigation = mustache.render(this.isMobile ? messageNavigationTemplateMobile : messageNavigationTemplate, {
			attachments:      this.i18n("attachments"),
			comments:         this.i18n("comments"),
			attachmentsHref:  function() { const params = new URLSearchParams(search); params.set("msg", "files");     return "?" + params.toString(); },
			commentsHref:     function() { const params = new URLSearchParams(search); params.set("msg", "comments");  return "?" + params.toString(); },
		})
	}

	addTranslations(message: number | string, previews: TranslationPreview[]) {
		const search = this.search

		this.translations = mustache.render(this.isMobile ? translationsTemplateMobile : translationsTemplate, {
			newTranslation:        this.i18n("newTranslation"),
			mainMessage:           this.i18n("mainMessage"),
			mainMessageHref:       function() { return `/messages/${message}` + search },
			newTranslationHref:    function() { return `/editor/messages/${message}/new_lang` + search },
			translationHref:       function() { return `/messages/${message}/${this.lang}` + search },
			translations:          previews,
			hasTranslations:       () => previews.length > 0,
		}, {
			newTranslation:        this.newTranslation,
		})
	}

	addNewTranslation(message: number | string) {
		const search = this.search

		this.newTranslation = mustache.render(this.isMobile ? newTranslationTemplateMobile : newTranslationTemplate, {
			newTranslation:        this.i18n("newTranslation"),
			newTranslationHref:    function() { return `/editor/messages/${message}/new_lang` + search },
		})
	}

	addControlPanel() {
		const search = this.search;

		this.controlPanel = mustache.render(this.isMobile ? controlPanelTemplateMobile : controlPanelTemplate, {}, {
			logout:           this.logout,
		})
	}

	addNavigation() {
		const search = this.search;

		this.navigation = mustache.render(this.isMobile ? navigationTemplateMobile : navigationTemplate, {
			messagesHref:          function() { return "/home" + search; },
			messagesSection:       this.i18n("messagesSection"),
			filesHref:             function() { return "/files" + search; },
			filesSection:          this.i18n("filesSection"),
		})
	}

	addLogout(logout: Builder) {
		this.logout = logout.build()
	}

	addFilesSelector(files: File[]) {
		this.filesSelector = mustache.render(this.isMobile ? filesSelectorTemplateMobile : filesSelectorTemplate, {
			files:             files,
			defaultFile:       this.i18n("defaultFile"),
		})
	}

	addFilesView(files: FileWithMime[]) {
		this.filesView = mustache.render(this.isMobile ? filesViewTemplateMobile : filesViewTemplate, {
			files:    files,
			imgSrc:   function() { return `/files/v1/read/${this.name}` },
			fileHref: function() { return `/files/v1/download?id=${this.ID}` },
		})
	}

	addHeader(header: Builder) {
		this.header = header.build()
	}

	build() {
		return mustache.render(this.isMobile ? homeTemplateMobile : homeTemplate, {}, {
			messageEditForm:      this.messageEditForm,
			messageView:          this.messageView,
			threadView:           this.threadView,
			header:               this.header,
			threadEditForm:       this.threadEditForm,
			newMessageForm:       this.newMessageForm,
			newTranslationForm:   this.newTranslationForm,
			translationEditForm:  this.translationEditForm,
			translationView:      this.translationView,
			messagesTree:         this.messagesTree,
			sidebar:              this.sidebar,
			pagination:           this.pagination,
			filesList:            this.filesList,
			filesForm:            this.filesForm,
			filesView:            this.filesView,
			controlPanel:         this.controlPanel,
			navigation:           this.navigation,
			translations:         this.translations,
			messageHeader:        this.messageHeader,
			comments:             this.comments,
			messageNavigation:    this.messageNavigation,
		});
	}
}

export default HomeBuilder
