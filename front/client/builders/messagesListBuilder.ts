import type {MessagesList} from '../api/models';
import Config from 'config';
import mustache from 'mustache';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import AbstractPublicBuilder from './abstractPublicBuilder'

let messagesListTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/messages_tree/desktop/messages_list.mustache')), { encoding: 'utf-8' });
let messagesListTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/messages_tree/mobile/messages_list.mustache')), { encoding: 'utf-8' });

class MessagesListBuilder extends AbstractPublicBuilder {
	threadID = 0
	tree = undefined

	addThreadID(id: number) {
		this.threadID = id
		return this
	}

	addTree(tree: MessagesList) {
		this.tree = tree
		return this
	}

	build() {
		const search = this.search
		const path = this.path
		const parentName = this.parentName
		const isAuthed = this.isAuthed
		const isMobile = this.isMobile
		const threadID = this.threadID
		const tree = this.tree

		const close = ((new URLSearchParams(search)).get("close") || "").split(",").map(parseFloat).filter(v => !isNaN(v))

		const limit = parseInt(LIMIT)

		return mustache.render(this.isMobile ? messagesListTemplateMobile : messagesListTemplate, {
			isLastPage:        tree.isLastPage,
			isFirstPage:       tree.isFirstPage,
			total:             tree.total,
			offset:            tree.offset,
			count:             tree.count,
			messages:          tree.messages,
			threadID:          threadID,
			showThreadLink:    function() { return isAuthed || (!this.thread.private && this.count > 0) },

			hasMessages:       function() { return this.messages.messages.length > 0 },
			isFolded:          function() { return this.messages.messages.length > 0 },
			hasPagination:     function() { return !(this.isLastPage && this.isFirstPage) },
			noMessagesText:    this.i18n("noMessagesText"),
			showCounter:       function() { return this.count > 0 },

			messageHref:       function() {
				if (parentName == "") {
					// root
					return `/messages/${this.ID}` + search
				} else {
					const params = new URLSearchParams(search);
					return `/${parentName}/${this.name}?` + params.toString();
				}
			},
			openThreadHref:    function() { 
				if (isMobile) {
					return `/${this.name}/tree` + search
				} else {
					return isAuthed ? `/${this.name}/new` + search : `/${this.name}` + search
				}
			},
			unfoldHref:        function() { const params = new URLSearchParams(search); params.set(this.ID || 0, `${limit},0`); return path + "?" + params.toString(); },
			foldHref:          function() { const params = new URLSearchParams(search); params.delete(this.ID || 0); return path + "?" + params.toString(); },
			loadMoreHref:      function() { const params = new URLSearchParams(search); params.set(this.threadID || 0, `${this.count+limit},${this.offset}`); return path + "?" + params.toString(); },
		}, {
			list: this.isMobile ? messagesListTemplateMobile : messagesListTemplate,
		})
	}
}

export default MessagesListBuilder
