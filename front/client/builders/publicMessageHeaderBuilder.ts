import type {Builder} from './builder'
import type { Identity } from '../api/models'
import Config from 'config';
import mustache from 'mustache';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import AbstractPublicBuilder from './abstractPublicBuilder'

let messageHeaderTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/message/desktop/message_header.mustache')), { encoding: 'utf-8' });
let messageHeaderTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/message/mobile/message_header.mustache')), { encoding: 'utf-8' });

class PublicMessageHeaderBuilder extends AbstractPublicBuilder {
	threadLink = undefined
	threadTitle = undefined

	addThreadLink(identity: Identity) {
		if (!this.isAuthed && identity.private) {
			return this
		}

		this.threadLink = "/" + identity.name + this.search
		this.threadTitle = identity.title

		return this
	}

	build() {
		return mustache.render(this.isMobile ? messageHeaderTemplateMobile : messageHeaderTemplate, {
			threadLink: this.threadLink,
			threadTitle: this.threadTitle,
		}, {})
	}
}

export default PublicMessageHeaderBuilder
