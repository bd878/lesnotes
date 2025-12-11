import Config from 'config';
import i18n from '../../i18n';
import mustache from 'mustache';
import { readFile } from 'node:fs/promises';
import { resolve, join } from 'node:path';
import Builder from '../builder';

async function login(ctx) {
	const { lang, theme, fontSize, query } = ctx.state

	const builder = new LoginBuilder(ctx.userAgent.isMobile, lang)

	await builder.addSettings(undefined, lang, theme, fontSize)
	await builder.addUsername()
	await builder.addPassword()
	await builder.addSubmit(query)
	await builder.addFooter()
	await builder.addSidebar(query)

	ctx.body = await builder.build(theme, fontSize)
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
	async addSubmit(query?: string) {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/login/mobile/submit.mustache' : 'templates/login/desktop/submit.mustache'
		)), { encoding: 'utf-8' });

		this.submit = mustache.render(template, {
			query:    query,
			register: this.i18n("register"),
			login:    this.i18n("login"),
		})
	}

	sidebar = undefined;
	async addSidebar(query?: string) {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/sidebar_horizontal/mobile/sidebar_horizontal.mustache' : 'templates/sidebar_horizontal/desktop/sidebar_horizontal.mustache'
		)), { encoding: 'utf-8' });

		this.sidebar = mustache.render(template, {query: query, settingsHeader: this.i18n("settingsHeader")}, {
			settings:       this.settings,
		})
	}

	async build(theme?: string, fontSize?: string) {
		const styles = await readFile(resolve(join(Config.get('basedir'), 'public/styles/styles.css')), { encoding: 'utf-8' });
		const layout = await readFile(resolve(join(Config.get('basedir'), 'templates/layout.mustache')), { encoding: 'utf-8' });
		const login = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/login/mobile/login.mustache' : 'templates/login/desktop/login.mustache'
		)), { encoding: 'utf-8' });

		return mustache.render(layout, {
			html:     () => (text, render) => {
				let html = "<html"

				if (theme) html += ` class="${theme}"`;
				if (this.lang) html += ` lang="${this.lang}"`;
				if (fontSize) html += ` data-size="${fontSize}"`
				html += ">"

				return html + render(text) + "</html>"
			},
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
