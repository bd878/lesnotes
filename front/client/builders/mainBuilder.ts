import type {Builder} from './builder';
import Config from 'config';
import i18n from '../i18n';
import mustache from 'mustache';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import AbstractBuilder from './abstractBuilder';

let authorizationTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/main/desktop/authorization.mustache')), { encoding: 'utf-8' });
let authorizationTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/main/mobile/authorization.mustache')), { encoding: 'utf-8' });

let mainTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/main/desktop/main.mustache')), { encoding: 'utf-8' });
let mainTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/main/mobile/main.mustache')), { encoding: 'utf-8' });

class MainBuilder extends AbstractBuilder {
	sidebar = undefined;
	authorization = undefined;

	addAuthorization() {
		this.authorization = mustache.render(this.isMobile ? authorizationTemplateMobile : authorizationTemplate, {
			loginHref: "/login" + this.search,
			signupHref: "/signup" + this.search,
			login:     this.i18n("login"),
			signup:    this.i18n("signup"),
		})
		return this
	}

	addSidebar(sidebar: Builder) {
		this.sidebar = sidebar.build()
		return this
	}

	build() {
		return mustache.render(this.isMobile ? mainTemplateMobile : mainTemplate, {
			botUsername:   `${BOT_USERNAME}`,
			authUrl:       `https://${BACKEND_URL}/tg_auth`,
		}, {
			sidebar:       this.sidebar,
			authorization: this.authorization,
		})
	}
}

export default MainBuilder
