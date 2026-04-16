import type {Builder} from './builder'
import type {Comment} from '../api/models'
import type {FileWithMime} from '../types';
import Config from 'config';
import mustache from 'mustache';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import AbstractBuilder from './abstractBuilder';

let featuresTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/message_features/desktop/message_features.mustache')), { encoding: 'utf-8' });
let featuresTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/message_features/mobile/message_features.mustache')), { encoding: 'utf-8' });

let filesViewTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/files_view/desktop/files_view.mustache')), { encoding: 'utf-8' });
let filesViewTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/files_view/mobile/files_view.mustache')), { encoding: 'utf-8' });

let newCommentTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/comments/desktop/new_comment.mustache')), { encoding: 'utf-8' });
let newCommentTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/comments/mobile/new_comment.mustache')), { encoding: 'utf-8' });

let commentsListTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/comments/desktop/comments_list.mustache')), { encoding: 'utf-8' });
let commentsListTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/comments/mobile/comments_list.mustache')), { encoding: 'utf-8' });

let commentsTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/comments/desktop/comments.mustache')), { encoding: 'utf-8' });
let commentsTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/comments/mobile/comments.mustache')), { encoding: 'utf-8' });

class MessageFeaturesBuilder extends AbstractBuilder {
	navigation = undefined
	comments = undefined;
	commentsList = undefined;
	newComment = undefined;
	filesView = undefined;
	translations = undefined;

	addNavigation(nav: Builder) {
		this.navigation = nav.build()
		return this
	}

	addTranslations(translations: Builder) {
		this.translations = translations.build()
		return this
	}

	addCommentsList(comments: Comment[]) {
		this.commentsList = mustache.render(this.isMobile ? commentsListTemplateMobile : commentsListTemplate, {
			comments: comments,
		})
		return this
	}

	addNewComment(message: number | string) {
		this.newComment = mustache.render(this.isMobile ? newCommentTemplateMobile : newCommentTemplate, {
			commentPlaceholder:       this.i18n("commentPlaceholder"),
			newComment:               this.i18n("newComment"),
			redirectUrl:              this.path + this.search,
			message:                  message,
			sendAction:               "/comment/send" + this.search,
		})
		return this
	}

	addComments(message: number | string, comments: Comment[]) {
		this.addNewComment(message)
		this.addCommentsList(comments)

		this.comments = mustache.render(this.isMobile ? commentsTemplateMobile : commentsTemplate, {}, {
			commentsList:  this.commentsList,
			newComment:    this.newComment,
		})
		return this
	}

	addFilesView(files: FileWithMime[]) {
		this.filesView = mustache.render(this.isMobile ? filesViewTemplateMobile : filesViewTemplate, {
			files:    files,
			imgSrc:   function() { return `/files/v1/read/${this.ID}` },
			fileHref: function() { return `/files/v1/download?id=${this.ID}` },
		})
		return this
	}

	build() {
		return mustache.render(this.isMobile ? featuresTemplateMobile : featuresTemplate, {}, {
			translations: this.translations,
			navigation:   this.navigation,
			filesView:    this.filesView,
			comments:     this.comments,
		})
	}
}

export default MessageFeaturesBuilder
