import Config from 'config';
import mustache from 'mustache';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import ContentBuilder from './contentBuilder'

let searchTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/content/desktop/search_content.mustache')), { encoding: 'utf-8' });
let searchTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/content/mobile/search_content.mustache')), { encoding: 'utf-8' });

let searchListTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/search_list/desktop/search_list.mustache')), { encoding: 'utf-8' });
let searchListTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/search_list/mobile/search_list.mustache')), { encoding: 'utf-8' });

class SearchContentBuilder extends ContentBuilder {
	messagesList = undefined

	addMessagesList(list) {
		const search = this.search
		const isMobile = this.isMobile

		const limit = parseInt(LIMIT)

		this.messagesList = mustache.render(this.isMobile ? searchListTemplateMobile : searchListTemplate, {
			messages:          list.messages,
			noMessagesText:    this.i18n("noMessagesText"),
			showThreadLink:    true,

			messageHref:       function() {
				const params = new URLSearchParams(search);
				params.delete("nav");
				params.delete("trans");
				return `/search/${this.ID}` + search;
			},
			openThreadHref:    function() {
				const params = new URLSearchParams(search);
				params.delete("nav");
				params.delete("trans");
				params.set(`${this.ID}`, `${limit},0`);
				return isMobile ? `/${this.name}/tree` + search : `/${this.name}` + search;
			},
		})

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

export default SearchContentBuilder
