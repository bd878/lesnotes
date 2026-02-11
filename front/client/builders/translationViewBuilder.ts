import type { Translation } from '../api/models';
import Config from 'config';
import mustache from 'mustache';
import api from '../api';
import * as is from '../third_party/is';
import { readFile } from 'node:fs/promises';
import { resolve, join } from 'node:path';
import HomeBuilder from './homeBuilder';

class TranslationViewBuilder extends HomeBuilder {
	async addTranslationView(messageID: number, translation: Translation) {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/home/mobile/translation_view.mustache' : 'templates/home/desktop/translation_view.mustache'
		)), { encoding: 'utf-8' });

		const search = this.search

		this.translationView = mustache.render(template, {
			messageID:        messageID,
			translation:      translation,
			editHref:         function() { return `/editor/messages/${messageID}/${this.lang}` + search; },
			deleteAction:     "/translation/delete" + search,
		})
	}
}

export default TranslationViewBuilder
