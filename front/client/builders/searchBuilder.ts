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
	messageView = undefined
	messagesList = undefined
	messageFeatures = undefined;
	messageHeader = undefined;
	controlPanel = undefined
	threadPath = undefined;

	addMessagesList(list: SearchMessagesList) {
		const search = this.search
		const isMobile = this.isMobile

		const limit = parseInt(LIMIT)

		this.messagesList = mustache.render(this.isMobile ? searchListTemplateMobile : searchListTemplate, {
			messages:          list.messages,
			noMessagesText:    this.i18n("noMessagesText"),
			showThreadLink:    true,

			messageHref:       function() { const params = new URLSearchParams(search); params.delete("nav"); params.delete("trans"); return `/search/${this.ID}` + search; },
			openThreadHref:    function() { const params = new URLSearchParams(search); params.delete("nav"); params.delete("trans"); params.set(`${this.ID}`, `${limit},0`); return isMobile ? `/${this.name}/tree` + search : `/${this.name}` + search; },
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

	addMessageHeader(header: Builder) {
		this.messageHeader = header.build()
		return this
	}

	addMessageFeatures(features: Builder) {
		this.messageFeatures = features.build()
		return this
	}

	addMessageView(view: Builder) {
		this.messageView = view.build()
		return this
	}

	addThreadPath(threadPath: Builder) {
		this.threadPath = threadPath.build()
		return this
	}

	build() {
		return mustache.render(this.isMobile ? searchTemplateMobile : searchTemplate, {
			hasFeatures: this.messageFeatures != undefined,
			hasMessage: this.messageView != undefined,
		}, {
			header: this.header,
			threadPath: this.threadPath,
			controlPanel: this.controlPanel,
			messageView: this.messageView,
			messagesList: this.messagesList,
			messageHeader: this.messageHeader,
			messageFeatures: this.messageFeatures,
		})
	}
}

export default SearchBuilder
