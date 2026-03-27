import type {Builder} from './builder'
import Config from 'config';
import mustache from 'mustache';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import AbstractBuilder from './abstractBuilder';

let searchTemplate = readFileSync(resolve(join(Config.get("basedir"), 'templates/header/desktop/search_form.mustache')), { encoding: 'utf-8' })
let searchTemplateMobile = readFileSync(resolve(join(Config.get("basedir"), 'templates/header/mobile/search_form.mustache')), { encoding: 'utf-8' })

let newNoteTemplate = readFileSync(resolve(join(Config.get("basedir"), 'templates/header/desktop/new_note.mustache')), { encoding: 'utf-8' })
let newNoteTemplateMobile = readFileSync(resolve(join(Config.get("basedir"), 'templates/header/mobile/new_note.mustache')), { encoding: 'utf-8' })

let headerTemplate = readFileSync(resolve(join(Config.get("basedir"), 'templates/header/desktop/header.mustache')), { encoding: 'utf-8' })
let headerTemplateMobile = readFileSync(resolve(join(Config.get("basedir"), 'templates/header/mobile/header.mustache')), { encoding: 'utf-8' })

class HeaderBuilder extends AbstractBuilder {
	searchForm = undefined;
	newNote = undefined;

	addSearch() {
		this.searchForm = mustache.render(this.isMobile ? searchTemplateMobile : searchTemplate, {
			action:              "/search" + this.search,
			searchPlaceholder:   this.i18n("searchPlaceholder"),
			searchMessages:      this.i18n("search"),
		})
	}

	addNewNote() {
		this.newNote = mustache.render(this.isMobile ? newNoteTemplateMobile : newNoteTemplate, {
			newNoteButton:       this.i18n("newNote"),
			newNoteHref:         "/home" + this.search,
		})
	}

	build() {
		return mustache.render(this.isMobile ? headerTemplateMobile : headerTemplate, {
			mainHref:   "/",
		}, {
			searchForm: this.searchForm,
			newNote:    this.newNote,
		})
	}
}

export default HeaderBuilder;
