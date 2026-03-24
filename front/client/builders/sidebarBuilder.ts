import type {Builder} from './builder';
import Config from 'config';
import i18n from '../i18n';
import mustache from 'mustache';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import AbstractBuilder from './abstractBuilder';

let sidebarTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/sidebar_horizontal/desktop/sidebar_horizontal.mustache')), { encoding: 'utf-8' });
let sidebarTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/sidebar_horizontal/mobile/sidebar_horizontal.mustache')), { encoding: 'utf-8' });

class SidebarBuilder extends AbstractBuilder {
	settings = undefined;

	addSettings(settings: Builder) {
		this.settings = settings.build()
	}

	build() {
		return mustache.render(this.isMobile ? sidebarTemplate : sidebarTemplateMobile, {
			mainHref:        "/" + this.search,
			settingsHeader:  this.i18n("settingsHeader"),
		}, {
			settings: this.settings,
		})
	}
}

export default SidebarBuilder
