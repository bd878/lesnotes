import type {MessagesList, Message} from '../api/models';
import Config from 'config';
import mustache from 'mustache';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import AbstractBuilder from './abstractBuilder';
import crop from '../utils/crop';

let messagesTreeTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/messages_tree/desktop/messages_tree.mustache')), { encoding: 'utf-8' });
let messagesTreeTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/messages_tree/mobile/messages_tree.mustache')), { encoding: 'utf-8' });

let messagesListTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/messages_tree/desktop/messages_list.mustache')), { encoding: 'utf-8' });
let messagesListTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/messages_tree/mobile/messages_list.mustache')), { encoding: 'utf-8' });

let threadPathTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/messages_tree/desktop/thread_path.mustache')), { encoding: 'utf-8' });
let threadPathTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/messages_tree/mobile/thread_path.mustache')), { encoding: 'utf-8' });

class MessagesTreeBuilder extends AbstractBuilder {
	list = undefined
	threadPath = undefined

	addList(tree: MessagesList) {
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

			hasMessages:       function() { return this.messages.messages.length > 0 },
			isFolded:          function() { return this.messages.messages.length > 0 },
			hasPagination:     function() { return !(this.isLastPage && this.isFirstPage) },
			noMessagesText:    this.i18n("noMessagesText"),

			messageHref:       function() { const params = new URLSearchParams(search); params.delete("nav"); params.delete("trans"); return `/messages/${this.ID}?` + params.toString(); },
			openThreadHref:    function() { const params = new URLSearchParams(search); params.delete("nav"); params.delete("trans"); params.set("cwd", this.ID); params.set(`${this.ID}`, `${limit},0`); return "/home?" + params.toString(); },
			unfoldHref:        function() { const params = new URLSearchParams(search); params.set(this.ID || 0, `${limit},0`); return path + "?" + params.toString(); },
			foldHref:          function() { const params = new URLSearchParams(search); params.delete(this.ID || 0); return path + "?" + params.toString(); },
			prevPageHref:      function() { const params = new URLSearchParams(search); params.set(this.ID || 0, `${limit},${limit + this.offset}`); return path + "?" + params.toString(); },
			nextPageHref:      function() { const params = new URLSearchParams(search); params.set(this.ID || 0, `${limit},${Math.max(0, this.offset - limit)}`); return path + "?" + params.toString(); },
		}, {
			list: this.isMobile ? messagesListTemplateMobile : messagesListTemplate,
		})
	}

	addThreadPath(threadPath: Message[]) {
		const search = this.search
		const path = this.path

		this.threadPath = mustache.render(this.isMobile ? threadPathTemplateMobile : threadPathTemplate, {
			path:     threadPath.slice(0, -1),
			thread:   threadPath.slice(-1),
			editThreadHref:  function() { return `/editor/threads/${this.ID}` + search },
			editThreadTitle: "<-- " + this.i18n("thread").toLowerCase(),
			threadTitle: function() { return `/${crop(this.title || this.text, 15)}` },
			threadHref: function() {
				const params = new URLSearchParams(search);
				switch (this.ID) {
				case 0:
					params.delete("cwd")
					return path + "?" + params.toString()
				default:
					params.set("cwd", this.ID)
					return path + "?" + params.toString()
				}
			},
		})
	}

	build() {
		return mustache.render(this.isMobile ? messagesTreeTemplateMobile : messagesTreeTemplate, {}, {
			list:       this.list,
			threadPath: this.threadPath,
		})
	}
}

export default MessagesTreeBuilder
