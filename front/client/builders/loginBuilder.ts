import type {Builder} from './builder'
import Config from 'config';
import i18n from '../i18n';
import mustache from 'mustache';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import AbstractBuilder from './abstractBuilder';

let usernameTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/login/desktop/username.mustache')), { encoding: 'utf-8' });
let usernameTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/login/mobile/username.mustache')), { encoding: 'utf-8' });

let passwordTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/login/desktop/password.mustache')), { encoding: 'utf-8' });
let passwordTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/login/mobile/password.mustache')), { encoding: 'utf-8' });

let submitTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/login/desktop/submit.mustache')), { encoding: 'utf-8' });
let submitTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/login/mobile/submit.mustache')), { encoding: 'utf-8' });

let loginTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/login/desktop/login.mustache')), { encoding: 'utf-8' });
let loginTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/login/mobile/login.mustache')), { encoding: 'utf-8' });

class LoginBuilder extends AbstractBuilder {
	username = undefined;
	password = undefined;
	submit = undefined;
	sidebar = undefined;

	addUsername() {
		this.username = mustache.render(this.isMobile ? usernameTemplateMobile : usernameTemplate, {
			usernamePlaceholder: this.i18n("username"),
		})
		return this
	}

	addPassword() {
		this.password = mustache.render(this.isMobile ? passwordTemplateMobile : passwordTemplate, {
			passwordPlaceholder: this.i18n("password"),
		})
		return this
	}

	addSubmit() {
		this.submit = mustache.render(this.isMobile ? submitTemplateMobile : submitTemplate, {
			signupHref: "/signup" + this.search,
			signup:   this.i18n("signup"),
			login:    this.i18n("login"),
		})
		return this
	}

	addSidebar(sidebar: Builder) {
		this.sidebar = sidebar.build()
		return this
	}

	build() {
		const search = this.search

		return mustache.render(this.isMobile ? loginTemplateMobile : loginTemplate, {
			action:        function() { return "/login" + search },
			botUsername:   `${BOT_USERNAME}`,
			authUrl:       `https://${BACKEND_URL}/tg_auth`,
		}, {
			username:      this.username,
			password:      this.password,
			submit:        this.submit,
			sidebar:       this.sidebar,
		})
	}
}

export default LoginBuilder
