import type {Builder} from './builder'
import Config from 'config';
import mustache from 'mustache';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import AbstractBuilder from './abstractBuilder';

let logoutTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/auth/desktop/logout.mustache')), { encoding: 'utf-8' });
let logoutTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/auth/mobile/logout.mustache')), { encoding: 'utf-8' });

let signupTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/auth/desktop/signup.mustache')), { encoding: 'utf-8' });
let signupTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/auth/mobile/signup.mustache')), { encoding: 'utf-8' });

let loginTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/auth/desktop/login.mustache')), { encoding: 'utf-8' });
let loginTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/auth/mobile/login.mustache')), { encoding: 'utf-8' });

let authTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/auth/desktop/auth.mustache')), { encoding: 'utf-8' });
let authTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/auth/mobile/auth.mustache')), { encoding: 'utf-8' });

class AuthBuilder extends AbstractBuilder {
	logout = undefined
	signup = undefined
	login = undefined

	addLogout() {
		const search = this.search

		this.logout = mustache.render(this.isMobile ? logoutTemplateMobile : logoutTemplate, {
			logout:           this.i18n("logout"),
			logoutHref:       "/logout" + search
		})
		return this
	}

	addSignup() {
		const search = this.search

		this.signup = mustache.render(this.isMobile ? signupTemplateMobile : signupTemplate, {
			signup:           this.i18n("signup"),
			signupHref:       "/signup" + search
		})
		return this
	}

	addLogin() {
		const search = this.search

		this.login = mustache.render(this.isMobile ? loginTemplateMobile : loginTemplate, {
			login:           this.i18n("login"),
			loginHref:       "/login" + search
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