import type {Builder} from './builder'
import type {SearchMessage} from '../api/models'
import Config from 'config';
import mustache from 'mustache';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import AbstractBuilder from './abstractBuilder'

let searchTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/search/desktop/search.mustache')), { encoding: 'utf-8' });
let searchTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/search/mobile/search.mustache')), { encoding: 'utf-8' });

let searchListTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/search/desktop/search_list.mustache')), { encoding: 'utf-8' });
let searchListTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/search/mobile/search_list.mustache')), { encoding: 'utf-8' });

class SearchBuilder extends AbstractBuilder {
	header = undefined
	messagesList = undefined
	controlPanel = undefined

	addMessagesList(messagesList: SearchMessagesList) {
		this.messagesList = mustache.render(this.isMobile ? searchListTemplateMobile : searchListTemplate, {
			messagesList: messagesList,
		})
		return this
	}

	addHeader(header: Builder) {
		this.header = header.build()
		return this
	}

	addControlPanel(panel: Builder) {
		this.controlPanel = panel.build()
		return this
	}

	build() {
		return mustache.render(this.isMobile ? searchTemplateMobile : searchTemplate, {}, {
			header: this.header,
			controlPanel: this.controlPanel,
			messagesList: this.messagesList,
		})
	}
}

export default SearchBuilder
