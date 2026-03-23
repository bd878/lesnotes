import Config from 'config';
import mustache from 'mustache';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import AbstractBuilder from './abstractBuilder';

let searchTemplate = readFileSync(resolve(join(Config.get("basedir"), 'templates/header/desktop/search_form.mustache')), { encoding: 'utf-8' })
let searchMobileTemplate = readFileSync(resolve(join(Config.get("basedir"), 'templates/header/mobile/search_form.mustache')), { encoding: 'utf-8' })
let headerTemplate = readFileSync(resolve(join(Config.get("basedir"), 'templates/header/desktop/header.mustache')), { encoding: 'utf-8' })
let headerMobileTemplate = readFileSync(resolve(join(Config.get("basedir"), 'templates/header/mobile/header.mustache')), { encoding: 'utf-8' })

class HeaderBuilder extends AbstractBuilder {
	search = undefined;

	addSearch() {
		this.search = mustache.render(this.isMobile ? searchTemplate : searchMobileTemplate, {
			action:              "/search" + this.search,
			searchPlaceholder:   this.i18n("searchPlaceholder"),
			searchMessages:      this.i18n("search"),
		})
	}

	build() {
		return mustache.render(this.isMobile ? headerMobileTemplate : headerTemplate, {
			mainHref:   "/" + this.search,
		}, {
			searchForm: this.search,
		})
	}
}

export default HeaderBuilder;
