import Config from 'config';
import i18n from '../../i18n';
import mustache from 'mustache';
import { readFile } from 'node:fs/promises';
import { resolve, join } from 'node:path';
import Builder from '../builder';

async function login(ctx) {
	const builder = new LoginBuilder(ctx.userAgent.isMobile, ctx.state.lang)

	await builder.addUsername()
	await builder.addPassword()
	await builder.addSubmit()
	await builder.addFooter()
	await builder.addSidebar()

	ctx.body = await builder.build()
	ctx.status = 200;
}

class LoginBuilder extends Builder {
	username = undefined;
	async addUsername() {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/login/mobile/username.mustache' : 'templates/login/desktop/username.mustache'
		)), { encoding: 'utf-8' });

		this.username = mustache.render(template, {
			usernamePlaceholder: this.i18n("username"),
		})
	}

	password = undefined;
	async addPassword() {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/login/mobile/password.mustache' : 'templates/login/desktop/password.mustache'
		)), { encoding: 'utf-8' });

		this.password = mustache.render(template, {
			passwordPlaceholder: this.i18n("password"),
		})
	}

	submit = undefined;
	async addSubmit() {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/login/mobile/submit.mustache' : 'templates/login/desktop/submit.mustache'
		)), { encoding: 'utf-8' });

		this.submit = mustache.render(template, {
			register: this.i18n("register"),
			login:    this.i18n("login"),
		})
	}

	sidebar = undefined;
	async addSidebar() {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/login/mobile/sidebar.mustache' : 'templates/login/desktop/sidebar.mustache'
		)), { encoding: 'utf-8' });

		this.sidebar = mustache.render(template)
	}

	async build() {
		const styles = await readFile(resolve(join(Config.get('basedir'), 'public/styles/styles.css')), { encoding: 'utf-8' });
		const layout = await readFile(resolve(join(Config.get('basedir'), 'templates/layout.mustache')), { encoding: 'utf-8' });
		const login = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/login/mobile/login.mustache' : 'templates/login/desktop/login.mustache'
		)), { encoding: 'utf-8' });

		return mustache.render(layout, {
			scripts:  ["/public/pages/login/loginScript.js"],
			manifest: "/public/manifest.json",
			styles:   styles,
			lang:     this.lang,
			isMobile: this.isMobile ? "true" : "",
		}, {
			footer:  this.footer,
			content: mustache.render(login, {}, {
				username:  this.username,
				password:  this.password,
				submit:    this.submit,
				sidebar:   this.sidebar,
			}),
		});
	}
}

export default login;
