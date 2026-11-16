import type {Builder} from './builder'
import Config from 'config';
import mustache from 'mustache';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import AbstractBuilder from './abstractBuilder'

let usernameTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/input/desktop/username.mustache')), { encoding: 'utf-8' });
let usernameTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/input/mobile/username.mustache')), { encoding: 'utf-8' });

let passwordTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/input/desktop/password.mustache')), { encoding: 'utf-8' });
let passwordTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/input/mobile/password.mustache')), { encoding: 'utf-8' });

let submitTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/input/desktop/submit.mustache')), { encoding: 'utf-8' });
let submitTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/input/mobile/submit.mustache')), { encoding: 'utf-8' });

let signupTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/signup/desktop/signup.mustache')), { encoding: 'utf-8' });
let signupTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/signup/mobile/signup.mustache')), { encoding: 'utf-8' });

class SignupBuilder extends AbstractBuilder {
	username = undefined;
	password = undefined;
	submit = undefined;
	sidebar = undefined;
	error = undefined;

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
			href: "/login" + this.search,
			submitText:   this.i18n("signup"),
			linkText:    this.i18n("login"),
		})
		return this
	}

	addSidebar(sidebar: Builder) {
		this.sidebar = sidebar.build()
		return this
	}

	addError(error: string) {
		this.error = error
		return this
	}

	build() {
		const search = this.search

		return mustache.render(this.isMobile ? signupTemplateMobile : signupTemplate, {
			action:         "/signup" + search,
			error:         this.error,
		}, {
			username:  this.username,
			password:  this.password,
			submit:    this.submit,
			sidebar:   this.sidebar,
		})
	}
}

export default SignupBuilder
