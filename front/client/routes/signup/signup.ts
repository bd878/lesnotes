import Config from 'config';
import i18n from '../../i18n';
import mustache from 'mustache';
import { readFile } from 'node:fs/promises';
import { resolve, join } from 'node:path';
import Builder from '../builder';

async function signup(ctx) {
	console.log("--> signup")

	const { lang, theme, fontSize, error } = ctx.state

	const builder = new SignupBuilder(ctx.userAgent.isMobile, lang, ctx.search)

	await builder.addSettings(undefined, lang, theme, fontSize)
	await builder.addUsername()
	await builder.addPassword()
	await builder.addSubmit(ctx.search)
	await builder.addFooter()
	await builder.addSidebar(ctx.search)

	ctx.body = await builder.build(theme, fontSize, ctx.search, error)
	ctx.status = 200;

	console.log("<-- signup")
}

class SignupBuilder extends Builder {
	username = undefined;
	async addUsername() {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/signup/mobile/username.mustache' : 'templates/signup/desktop/username.mustache'
		)), { encoding: 'utf-8' });

		this.username = mustache.render(template, {
			usernamePlaceholder: this.i18n("username"),
		})
	}

	password = undefined;
	async addPassword() {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/signup/mobile/password.mustache' : 'templates/signup/desktop/password.mustache'
		)), { encoding: 'utf-8' });

		this.password = mustache.render(template, {
			passwordPlaceholder: this.i18n("password"),
		})
	}

	submit = undefined;
	async addSubmit(query?: string) {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/signup/mobile/submit.mustache' : 'templates/signup/desktop/submit.mustache'
		)), { encoding: 'utf-8' });

		this.submit = mustache.render(template, {
			query:    query,
			signup:   this.i18n("signup"),
			login:    this.i18n("login"),
		})
	}

	sidebar = undefined;
	async addSidebar(query?: string) {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/sidebar_horizontal/mobile/sidebar_horizontal.mustache' : 'templates/sidebar_horizontal/desktop/sidebar_horizontal.mustache'
		)), { encoding: 'utf-8' });

		this.sidebar = mustache.render(template, {query: query, settingsHeader: this.i18n("settingsHeader")}, {settings: this.settings})
	}

	footer = undefined;
	async addFooter() {
		const template = await readFile(resolve(join(Config.get("basedir"), 'templates/footer.mustache')), { encoding: 'utf-8' });

		this.footer = mustache.render(template, {
			terms:            this.i18n("terms"),
			contact:          this.i18n("contact"),
			docs:             this.i18n("docs"),
		})
	}

	async build(theme?: string, fontSize?: string, search?: string, error?: string) {
		const styles = await readFile(resolve(join(Config.get('basedir'), 'public/styles/styles.css')), { encoding: 'utf-8' });
		const layout = await readFile(resolve(join(Config.get('basedir'), 'templates/layout.mustache')), { encoding: 'utf-8' });
		const signup = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/signup/mobile/signup.mustache' : 'templates/signup/desktop/signup.mustache'
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
			manifest: "/public/manifest.json",
			styles:   styles,
			lang:     this.lang,
			isMobile: this.isMobile ? "true" : "",
		}, {
			footer:  this.footer,
			content: mustache.render(signup, {
				action:         function() { return "/signup" + search },
				error:          error,
				settingsHeader: this.i18n("settingsHeader"),
				botUsername:    `${BOT_USERNAME}`,
				authUrl:        `https://${BACKEND_URL}/tg_auth`,
			}, {
				username:  this.username,
				password:  this.password,
				submit:    this.submit,
				sidebar:   this.sidebar,
			}),
		});
	}
}

export default signup;
