import type {Builder} from './builder'
import Config from 'config';
import mustache from 'mustache';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import AbstractBuilder from './abstractBuilder';

let linkTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/link/desktop/link.mustache')), { encoding: 'utf-8' });
let linkTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/link/mobile/link.mustache')), { encoding: 'utf-8' });

let authTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/auth/desktop/auth.mustache')), { encoding: 'utf-8' });
let authTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/auth/mobile/auth.mustache')), { encoding: 'utf-8' });

class AuthBuilder extends AbstractBuilder {
	logout = undefined
	signup = undefined
	login = undefined

	addLogout() {
		const search = this.search

		this.logout = mustache.render(this.isMobile ? linkTemplateMobile : linkTemplate, {
			text:           this.i18n("logout"),
			href:       "/logout" + search
		})
		return this
	}

	addSignup() {
		const search = this.search

		this.signup = mustache.render(this.isMobile ? linkTemplateMobile : linkTemplate, {
			text:           this.i18n("signup"),
			href:       "/signup" + search
		})
		return this
	}

	addLogin() {
		const search = this.search

		this.login = mustache.render(this.isMobile ? linkTemplateMobile : linkTemplate, {
			text:           this.i18n("login"),
			href:       "/login" + search
		})
		return this
	}

	build() {
		return mustache.render(this.isMobile ? authTemplateMobile : authTemplate, {}, {
			logout: this.logout,
			signup: this.signup,
			login:  this.login,
		})
	}
}

export default AuthBuilder