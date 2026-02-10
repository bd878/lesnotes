import type { Message } from '../api/models';
import Config from 'config';
import mustache from 'mustache';
import api from '../api';
import * as is from '../third_party/is';
import { readFile } from 'node:fs/promises';
import { resolve, join } from 'node:path';
import HomeBuilder from './homeBuilder';

class MessageViewBuilder extends HomeBuilder {
	async addMessageView(userID: number, message?: Message) {
		if (is.empty(message)) {
			return
		}

		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/home/mobile/message_view.mustache' : 'templates/home/desktop/message_view.mustache'
		)), { encoding: 'utf-8' });

		const search = this.search

		this.messageView = mustache.render(template, {
			ID:               message.ID,
			title:            message.title,
			text:             message.text,
			name:             message.name,
			private:          message.private,
			filesSummary:     this.i18n("filesSummary"),
			newNoteHref:      function() { return "/home" + search; },
			editHref:         function() { return `/editor/messages/${message.ID}` + search; },
			deleteAction:     "/m/delete" + search,
			publishAction:    "/m/publish" + search,
			privateAction:    "/m/private" + search,
			newNoteButton:    this.i18n("newNote"),
			userID:           userID,
			domain:           Config.get("domain"),
		})
	}
}

export default MessageViewBuilder
