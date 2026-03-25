import type { Translation } from '../api/models';
import Config from 'config';
import mustache from 'mustache';
import api from '../api';
import * as is from '../third_party/is';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import HomeBuilder from './homeBuilder';

let translationViewTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/home/desktop/translation_view.mustache')), { encoding: 'utf-8' });
let translationViewTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/home/mobile/translation_view.mustache')), { encoding: 'utf-8' });

class TranslationViewBuilder extends HomeBuilder {
	addTranslationView(messageID: number, translation: Translation) {
		const search = this.search

		this.translationView = mustache.render(this.isMobile ? translationViewTemplateMobile : translationViewTemplate, {
			messageID:        messageID,
			translation:      translation,
			editHref:         function() { return `/editor/messages/${messageID}/${this.lang}` + search; },
			deleteAction:     "/translation/delete" + search,
		})
	}
}

export default TranslationViewBuilder
