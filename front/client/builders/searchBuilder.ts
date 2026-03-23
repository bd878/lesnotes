import type { Message } from '../api/models';
import type { File } from '../api/models';
import Config from 'config';
import mustache from 'mustache';
import api from '../api';
import * as is from '../third_party/is';
import i18n from '../i18n';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import AbstractBuilder from './abstractBuilder';

let messagesListTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/search/desktop/messages_list.mustache')), { encoding: 'utf-8' });
let messagesListTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/search/mobile/messages_list.mustache')), { encoding: 'utf-8' });

let filesListTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/search/desktop/files_list.mustache')), { encoding: 'utf-8' });
let filesListTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/search/mobile/files_list.mustache')), { encoding: 'utf-8' });

let logoutTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/sidebar_vertical/desktop/logout.mustache')), { encoding: 'utf-8' });
let logoutTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/sidebar_vertical/mobile/logout.mustache')), { encoding: 'utf-8' });

let searchFormTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/search/desktop/search_form.mustache')), { encoding: 'utf-8' });
let searchFormTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/search/mobile/search_form.mustache')), { encoding: 'utf-8' });

let searchTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/search/desktop/search.mustache')), { encoding: 'utf-8' });
let searchTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/search/mobile/search.mustache')), { encoding: 'utf-8' });

class SearchBuilder extends AbstractBuilder {
	messagesList = undefined;
	filesList = undefined;
	logout = undefined;
	searchForm = undefined;
	sidebar = undefined;

	addMessagesList(list: Message[]) {
		const search = this.search

		this.messagesList = mustache.render(this.isMobile ? messagesListTemplate : messagesListTemplateMobile, {
			list:             list,
			isEmpty:          () => list.length == 0,
			isSingle:         () => list.length == 1,
			emptyListText:    this.i18n("emptyListText"),
			messageHref:      function() { const params = new URLSearchParams(search); params.delete("query"); return `/messages/${this.ID}?` + params.toString(); },
		})
	}

	addFilesList(list?: File[]) {
		const options = {
			filesPlaceholder:   this.i18n("filesPlaceholder"),
			noFiles:            this.i18n("noFiles"),
			files:              undefined,
		}

		if (is.notEmpty(list)) {
			options.files = list
		}

		this.filesList = mustache.render(this.isMobile ? filesListTemplate : filesListTemplateMobile, options)
	}

	addLogout() {
		const search = this.search

		this.logout = mustache.render(this.isMobile ? logoutTemplate : logoutTemplateMobile, {
			logout:           this.i18n("logout"),
			logoutHref:       function() { const params = new URLSearchParams(search); params.delete("cwd"); params.delete("id"); /* TODO: delete pagination */ return "/logout?" + params.toString() },
		})
	}

	addSearch() {
		this.searchForm = mustache.render(this.isMobile ? searchFormTemplate : searchFormTemplateMobile, {
			action:              "/search" + this.search,
			searchPlaceholder:   this.i18n("searchPlaceholder"),
			searchMessages:      this.i18n("search"),
		})
	}

	build() {
		return mustache.render(this.isMobile ? searchTemplate : searchTemplateMobile, {}, {
			messagesList:    this.messagesList,
			filesList:       this.filesList,
			sidebar:         this.sidebar,
		})
	}
}

export default SearchBuilder
