import type {Builder} from './builder'
import Config from 'config';
import mustache from 'mustache';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import AbstractBuilder from './abstractBuilder';

let searchTemplate = readFileSync(resolve(join(Config.get("basedir"), 'templates/header/desktop/search_form.mustache')), { encoding: 'utf-8' })
let searchTemplateMobile = readFileSync(resolve(join(Config.get("basedir"), 'templates/header/mobile/search_form.mustache')), { encoding: 'utf-8' })

class SearchFormBuilder extends AbstractBuilder {
	threadID  = 0
	threadName = ""

	addThreadID(threadID: number) {
		this.threadID = threadID
		return this
	}

	addThreadName(name: string) {
		this.threadName = name
		return this
	}

	build() {
		return mustache.render(this.isMobile ? searchTemplateMobile : searchTemplate, {
			action:              "/message/search" + this.search,
			threadID:            this.threadID,
			threadName:          this.threadName,
			searchPlaceholder:   this.i18n("searchPlaceholder"),
			searchMessages:      this.i18n("search"),
		})
	}
}

export default SearchFormBuilder
