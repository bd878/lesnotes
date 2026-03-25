import type {MessagesList} from '../api/models';
import Config from 'config';
import mustache from 'mustache';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import AbstractBuilder from './abstractBuilder';

let messagesTreeTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/messages_tree/desktop/messages_tree.mustache')), { encoding: 'utf-8' });
let messagesTreeTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/messages_tree/mobile/messages_tree.mustache')), { encoding: 'utf-8' });

let messagesListTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/messages_tree/desktop/messages_list.mustache')), { encoding: 'utf-8' });
let messagesListTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/messages_tree/mobile/messages_list.mustache')), { encoding: 'utf-8' });

class MessagesTreeBuilder extends AbstractBuilder {
	list = undefined

	addList(tree: MessagesList) {
		const search = this.search
		const path = this.path

		const close = ((new URLSearchParams(search)).get("close") || "").split(",").map(parseFloat).filter(v => !isNaN(v))

		const limit = parseInt(LIMIT)

		this.list = mustache.render(this.isMobile ? messagesListTemplateMobile : messagesListTemplate, {
			isLastPage:       tree.isLastPage,
			isFirstPage:      tree.isFirstPage,
			total:            tree.total,
			count:            tree.count,
			messages:         tree.messages,

			hasMessages:      function() { return this.messages.messages.length > 0 },
			hasPagination:    function() { return !(this.isLastPage && this.isFirstPage) },
			noMessagesText:   this.i18n("noMessagesText"),

			messageHref:      function() { return `/messages/${this.ID}` + search; },
			messageThreadHref: function() { const params = new URLSearchParams(search); params.set(this.ID || 0, `${limit},0`); return path + "?" + params.toString(); },
			viewThreadHref:   function() { return `/threads/${this}` + search; /*context is ID, not thread*/ },
			prevPageHref:     function() { const params = new URLSearchParams(search); params.set(this.ID || 0, `${limit},${limit + this.offset}`); return path + "?" + params.toString(); },
			nextPageHref:     function() { const params = new URLSearchParams(search); params.set(this.ID || 0, `${limit},${Math.max(0, this.offset - limit)}`); return path + "?" + params.toString(); },
		}, {
			list: this.isMobile ? messagesListTemplateMobile : messagesListTemplate,
		})
	}


	build() {
		return mustache.render(this.isMobile ? messagesTreeTemplateMobile : messagesTreeTemplate, {}, {
			list:  this.list,
		})
	}
}

export default MessagesTreeBuilder
