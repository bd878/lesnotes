import type { Message } from '../api/models';
import Config from 'config';
import mustache from 'mustache';
import api from '../api';
import * as is from '../third_party/is';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import HomeBuilder from './homeBuilder';

let messageViewTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/home/desktop/message_view.mustache')), { encoding: 'utf-8' });
let messageViewTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/home/mobile/message_view.mustache')), { encoding: 'utf-8' });

class MessageViewBuilder extends HomeBuilder {
	addMessageView(userID: number, message?: Message) {
		if (is.empty(message)) {
			return
		}

		const search = this.search

		this.messageView = mustache.render(this.isMobile ? messageViewTemplateMobile : messageViewTemplate, {
			ID:                    message.ID,
			title:                 message.title,
			text:                  message.text,
			name:                  message.name,
			private:               message.private,
			filesSummary:          this.i18n("filesSummary"),
			editHref:              `/editor/messages/${message.ID}` + search,
			threadHref:            `/editor/threads/${message.ID}` + search,
			deleteAction:          "/m/delete" + search,
			publishAction:         "/m/publish" + search,
			privateAction:         "/m/private" + search,
			newNoteButton:         this.i18n("newNote"),
			userID:                userID,
			domain:                Config.get("domain"),
		})
	}
}

export default MessageViewBuilder
