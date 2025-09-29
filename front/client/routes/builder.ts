import Config from 'config';
import i18n from '../i18n';
import mustache from 'mustache';
import { readFile } from 'node:fs/promises';
import { resolve, join } from 'node:path';

abstract class Builder {
	isMobile: boolean = false;
	lang:     string  = "en";

	constructor(isMobile: boolean, lang: string = "en") {
		this.isMobile = isMobile
		this.lang = lang
	}

	i18n(key: string): string {
		return i18n(this.lang)(key)
	}

	abstract build();

	footer = undefined;
	async addFooter() {
		const template = await readFile(resolve(join(Config.get("basedir"), 'templates/footer.mustache')), { encoding: 'utf-8' });

		this.footer = mustache.render(template, {
			terms:            this.i18n("terms"),
			contact:          this.i18n("contact"),
			docs:             this.i18n("docs"),
		})
	}
}

export default Builder
