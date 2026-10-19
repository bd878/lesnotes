import type {MessagesList, Message} from '../api/models';
import Config from 'config';
import mustache from 'mustache';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import AbstractPublicBuilder from './abstractPublicBuilder';
import crop from '../utils/crop';

let messagesTreeTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/messages_tree/desktop/messages_tree.mustache')), { encoding: 'utf-8' });
let messagesTreeTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/messages_tree/mobile/messages_tree.mustache')), { encoding: 'utf-8' });

let messagesListTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/messages_tree/desktop/messages_list.mustache')), { encoding: 'utf-8' });
let messagesListTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/messages_tree/mobile/messages_list.mustache')), { encoding: 'utf-8' });

let threadPathTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/messages_tree/desktop/thread_path.mustache')), { encoding: 'utf-8' });
let threadPathTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/messages_tree/mobile/thread_path.mustache')), { encoding: 'utf-8' });

class MessagesTreeBuilder extends AbstractPublicBuilder {
	list = undefined
	threadPath = undefined

	addList(tree: MessagesList, threadID: number) {
		const search = this.search
		const path = this.path

		const close = ((new URLSearchParams(search)).get("close") || "").split(",").map(parseFloat).filter(v => !isNaN(v))

		const limit = parseInt(LIMIT)

		this.list = mustache.render(this.isMobile ? messagesListTemplateMobile : messagesListTemplate, {
			isLastPage:        tree.isLastPage,
			isFirstPage:       tree.isFirstPage,
			total:             tree.total,
			offset:            tree.offset,
			count:             tree.count,
			messages:          tree.messages,
			threadID:          threadID,
			showThreadLink:    true,

			hasMessages:       function() { return this.messages.messages.length > 0 },
			isFolded:          function() { return this.messages.messages.length > 0 },
			hasPagination:     function() { return !(this.isLastPage && this.isFirstPage) },
			noMessagesText:    this.i18n("noMessagesText"),
			showCounter:       function() { return this.count > 0 },

			messageHref:       function() { const params = new URLSearchParams(search); params.delete("nav"); params.delete("trans"); return `/messages/${this.ID}?` + params.toString(); },
			openThreadHref:    function() { return `/${this.name}/new` },
			unfoldHref:        function() { const params = new URLSearchParams(search); params.set(this.ID || 0, `${limit},0`); return path + "?" + params.toString(); },
			foldHref:          function() { const params = new URLSearchParams(search); params.delete(this.ID || 0); return path + "?" + params.toString(); },
			loadMoreHref:      function() { const params = new URLSearchParams(search); params.set(this.threadID || 0, `${this.count+limit},${this.offset}`); return path + "?" + params.toString(); },
		}, {
			list: this.isMobile ? messagesListTemplateMobile : messagesListTemplate,
		})
		return this
	}

	addThreadPath(threadPath: Message[]) {
		const search = this.search
		const urlPath = this.path
		const isAuthed = this.isAuthed

		this.threadPath = mustache.render(this.isMobile ? threadPathTemplateMobile : threadPathTemplate, {
			path:     threadPath.slice(0, -1),
			lastMessage:   threadPath.slice(-1),
			editPublicThread: this.isAuthed,
			threadTitle: function() { return `/${crop(this.title || this.text, 30)}` },
			lastThreadTarget: "_self",
			lastThreadExternal: false,
			lastThreadHref:  function() { return `/${this.name}` },
			pathThreadHref: function() { return `/${this.name}` + search },
		})
		return this
	}

	build() {
		return mustache.render(this.isMobile ? messagesTreeTemplateMobile : messagesTreeTemplate, {}, {
			list:       this.list,
			threadPath: this.threadPath,
		})
	}
}

export default MessagesTreeBuilder
