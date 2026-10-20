import type {Builder} from './builder'
import Config from 'config';
import mustache from 'mustache';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import AbstractBuilder from './abstractBuilder';

let newNoteTemplate = readFileSync(resolve(join(Config.get("basedir"), 'templates/header/desktop/new_note.mustache')), { encoding: 'utf-8' })
let newNoteTemplateMobile = readFileSync(resolve(join(Config.get("basedir"), 'templates/header/mobile/new_note.mustache')), { encoding: 'utf-8' })

let headerTemplate = readFileSync(resolve(join(Config.get("basedir"), 'templates/header/desktop/header.mustache')), { encoding: 'utf-8' })
let headerTemplateMobile = readFileSync(resolve(join(Config.get("basedir"), 'templates/header/mobile/header.mustache')), { encoding: 'utf-8' })

class HeaderBuilder extends AbstractBuilder {
	searchForm = undefined;
	newNote = undefined;
	auth = undefined;

	addSearch(search: Builder) {
		this.searchForm = search.build()
		return this
	}

	addNewNote(href?: string) {
		this.newNote = mustache.render(this.isMobile ? newNoteTemplateMobile : newNoteTemplate, {
			newNoteButton:       this.i18n("newNote"),
			newNoteHref:         href || "/home" + this.search,
		})
		return this
	}

	addAuth(auth: Builder) {
		this.auth = auth.build()
		return this
	}

	build() {
		return mustache.render(this.isMobile ? headerTemplateMobile : headerTemplate, {
			mainHref:   "/",
		}, {
			searchForm: this.searchForm,
			newNote:    this.newNote,
			auth:       this.auth,
		})
	}
}

export default HeaderBuilder;
