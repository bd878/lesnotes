import type {Message} from '../api/models';
import Config from 'config';
import mustache from 'mustache';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import AbstractBuilder from './abstractBuilder'
import crop from '../utils/crop';

let threadPathTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/messages_tree/desktop/thread_path.mustache')), { encoding: 'utf-8' });
let threadPathTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/messages_tree/mobile/thread_path.mustache')), { encoding: 'utf-8' });

class ThreadPathBuilder extends AbstractBuilder {
	path = undefined

	addPath(path: Message[]) {
		this.path = path
		return this
	}

	build() {
		const search = this.search
		const isMobile = this.isMobile

		return mustache.render(this.isMobile ? threadPathTemplateMobile : threadPathTemplate, {
			path:     this.path.slice(0, -1),
			lastMessage:   this.path.slice(-1),
			threadPrivate: function() { return this.private },
			editPublicThread: this.isAuthed,
			editThreadHref: function() { return `/messages/${this.ID}/edit` + search },

			lastThreadTarget: "_self",
			lastThreadExternal: false,
			lastThreadHref: function() { return (isMobile ? this.name ? `/${this.name}/tree` : "/" : `/${this.name}`) + search },
			threadTitle: function() { return `/${crop(this.title || this.text, 30)}` },
			pathThreadHref: function() { return (isMobile ? this.name ? `/${this.name}/tree` : "/" : `/${this.name}`) + search }
		})
	}
}

export default ThreadPathBuilder
