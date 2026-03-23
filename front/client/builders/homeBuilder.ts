import type { File, Message, Thread, Comment, ThreadMessages } from '../api/models';
import type { FileWithMime } from '../types';
import type { TranslationPreview } from '../api/models';
import Config from 'config';
import mustache from 'mustache';
import api from '../api';
import { unwrapPaging } from '../api/models/paging';
import * as is from '../third_party/is';
import i18n from '../i18n';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import AbstractBuilder from './abstractBuilder'

let messagesTreeTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/home/desktop/messages_tree.mustache')), { encoding: 'utf-8' });
let messagesTreeTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/home/mobile/messages_tree.mustache')), { encoding: 'utf-8' });

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

let logoutTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/sidebar_vertical/desktop/logout.mustache')), { encoding: 'utf-8' });
let logoutTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/sidebar_vertical/mobile/logout.mustache')), { encoding: 'utf-8' });

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
	messagesStack        = undefined;
	logout               = undefined;
	sidebar              = undefined;
	controlPanel         = undefined;
	navigation           = undefined;
	translations         = undefined;
	newTranslation       = undefined;
	messageNavigation    = undefined;
	newComment           = undefined;
	commentsList         = undefined;
	comments             = undefined;
	scripts              = ["/public/pages/home/homeScript.js"]

	addMessagesTree(stack: ThreadMessages[]) {
		const search = this.search
		const path = this.path

		const close = ((new URLSearchParams(search)).get("close") || "").split(",").map(parseFloat).filter(v => !isNaN(v))

		const limit = parseInt(LIMIT)

		this.messagesStack = mustache.render(this.isMobile ? messagesTreeTemplate : messagesTreeTemplateMobile, {
			stack:            stack.map(unwrapPaging),
			limit:            LIMIT,
			isSingle:         () => stack.length == 1,
			openHref:         function() { const params = new URLSearchParams(search); params.set("close", this.closeIDs); return path + "?" + params.toString() },
			closeHref:        function() { const params = new URLSearchParams(search); params.set("close", this.openIDs); return path + "?" + params.toString() },
			newMessageText:   this.i18n("newMessageText"),
			noMessagesText:   this.i18n("noMessagesText"),
			messageHref:      function() { return `/messages/${this.ID}` + search; },
			messageThreadHref: function() { const params = new URLSearchParams(search); params.set("cwd", `${this.ID}`); return path + "?" + params.toString(); },
			viewThreadHref:   function() { return `/threads/${this}` + search; /*context is ID, not thread*/ },
			closeThreadHref:  function() {
				const params = new URLSearchParams(search);
				const set = new Set(close);
				set.add(this.ID);
				(this.parentID == 0 ? params.delete("cwd") : params.set("cwd", `${this.parentID}`));
				params.set("close", Array.from(set).join(","));
				return path + "?" + params.toString();
			},
			rootThreadHref:   function() { const params = new URLSearchParams(search); params.delete("cwd"); return path + "?" + params.toString(); },
			prevPageHref:     function() { const params = new URLSearchParams(search); params.set(this.message.ID || 0, `${limit + this.offset}`); return path + "?" + params.toString(); },
			nextPageHref:     function() { const params = new URLSearchParams(search); params.set(this.message.ID || 0, `${Math.max(0, this.offset - limit)}`); return path + "?" + params.toString(); },
		})
	}

	addNewComment(message: number | string) {
		this.newComment = mustache.render(this.isMobile ? newCommentTemplate : newCommentTemplateMobile, {
			commentPlaceholder:       this.i18n("commentPlaceholder"),
			newComment:               this.i18n("newComment"),
			redirectUrl:              this.path + this.search,
			message:                  message,
			sendAction:               "/comment/send" + this.search,
		})
	}

	addCommentsList(comments: Comment[]) {
		this.commentsList = mustache.render(this.isMobile ? commentsListTemplate : commentsListTemplateMobile, {
			comments: comments,
		})
	}

	addComments(message: number | string, comments: Comment[]) {
		this.addNewComment(message)
		this.addCommentsList(comments)

		this.comments = mustache.render(this.isMobile ? commentsTemplate : commentsTemplateMobile, {}, {
			commentsList:  this.commentsList,
			newComment:    this.newComment,
		})
	}

	addMessageNavigation() {
		const search = this.search

		this.messageNavigation = mustache.render(this.isMobile ? messageNavigationTemplate : messageNavigationTemplateMobile, {
			attachments:      this.i18n("attachments"),
			comments:         this.i18n("comments"),
			attachmentsHref:  function() { const params = new URLSearchParams(search); params.set("msg", "files");     return "?" + params.toString(); },
			commentsHref:     function() { const params = new URLSearchParams(search); params.set("msg", "comments");  return "?" + params.toString(); },
		})
	}

	addTranslations(message: number | string, previews: TranslationPreview[]) {
		const search = this.search

		this.translations = mustache.render(this.isMobile ? translationsTemplate : translationsTemplateMobile, {
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

		this.newTranslation = mustache.render(this.isMobile ? newTranslationTemplate : newTranslationTemplateMobile, {
			newTranslation:        this.i18n("newTranslation"),
			newTranslationHref:    function() { return `/editor/messages/${message}/new_lang` + search },
		})
	}

	addControlPanel() {
		const search = this.search;

		this.controlPanel = mustache.render(this.isMobile ? controlPanelTemplate : controlPanelTemplateMobile, {
			newNoteHref:      function() { return "/home" + search; },
			newNoteButton:    this.i18n("newNote"),
			newFileHref:      function() { return "/files" + search; },
			newFileButton:    this.i18n("newFile"),
		})
	}

	addNavigation() {
		const search = this.search;

		this.navigation = mustache.render(this.isMobile ? navigationTemplate : navigationTemplateMobile, {
			messagesHref:          function() { return "/home" + search; },
			messagesSection:       this.i18n("messagesSection"),
			filesHref:             function() { return "/files" + search; },
			filesSection:          this.i18n("filesSection"),
		})
	}

	addLogout() {
		const search = this.search

		this.logout = mustache.render(this.isMobile ? logoutTemplate : logoutTemplateMobile, {
			logout:           this.i18n("logout"),
			logoutHref:       function() { const params = new URLSearchParams(search); params.delete("cwd"); params.delete("id"); /* TODO: delete pagination */ return "/logout?" + params.toString() },
		})
	}

	addFilesSelector(files: File[]) {
		this.filesSelector = mustache.render(this.isMobile ? filesSelectorTemplate : filesSelectorTemplateMobile, {
			files:             files,
			defaultFile:       this.i18n("defaultFile"),
		})
	}

	addFilesView(files: FileWithMime[]) {
		this.filesView = mustache.render(this.isMobile ? filesViewTemplate : filesViewTemplateMobile, {
			files:    files,
			imgSrc:   function() { return `/files/v1/read/${this.name}` },
			fileHref: function() { return `/files/v1/download?id=${this.ID}` },
		})
	}

	build() {
		return mustache.render(this.isMobile ? homeTemplate : homeTemplateMobile, {}, {
			messageEditForm:      this.messageEditForm,
			messageView:          this.messageView,
			threadView:           this.threadView,
			threadEditForm:       this.threadEditForm,
			newMessageForm:       this.newMessageForm,
			newTranslationForm:   this.newTranslationForm,
			translationEditForm:  this.translationEditForm,
			translationView:      this.translationView,
			messagesStack:        this.messagesStack,
			sidebar:              this.sidebar,
			pagination:           this.pagination,
			filesList:            this.filesList,
			filesForm:            this.filesForm,
			filesView:            this.filesView,
			controlPanel:         this.controlPanel,
			navigation:           this.navigation,
			translations:         this.translations,
			comments:             this.comments,
			messageNavigation:    this.messageNavigation,
		});
	}
}

export default HomeBuilder
