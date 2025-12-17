import Config from 'config';
import i18n from '../../i18n';
import mustache from 'mustache';
import { readFile } from 'node:fs/promises';
import { resolve, join } from 'node:path';
import Builder from '../builder';

async function main(ctx) {
	console.log("--> main")

	const builder = new MainBuilder(ctx.userAgent.isMobile, ctx.state.lang, ctx.search)

	await builder.addSettings(ctx.state.lang, ctx.state.theme, ctx.state.fontSize)
	await builder.addFooter()
	await builder.addSidebar()
	await builder.addAuthorization()

	ctx.body = await builder.build(ctx.state.theme, ctx.state.fontSize)
	ctx.status = 200;

	console.log("<-- main")
}

class MainBuilder extends Builder {
	sidebar = undefined;
	async addSidebar() {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/sidebar_horizontal/mobile/sidebar_horizontal.mustache' : 'templates/sidebar_horizontal/desktop/sidebar_horizontal.mustache'
		)), { encoding: 'utf-8' });

		this.sidebar = mustache.render(template, {mainHref: "/" + this.search, settingsHeader: this.i18n("settingsHeader")}, {settings: this.settings})
	}

	authorization = undefined;
	async addAuthorization() {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/main/mobile/authorization.mustache' : 'templates/main/desktop/authorization.mustache'
		)), { encoding: 'utf-8' });

		this.authorization = mustache.render(template, {
			loginHref: "/login" + this.search,
			signupHref: "/signup" + this.search,
			login:     this.i18n("login"),
			signup:    this.i18n("signup"),
		})
	}

	async build(theme?: string, fontSize?: string) {
		const styles = await readFile(resolve(join(Config.get('basedir'), 'public/styles/styles.css')), { encoding: 'utf-8' });
		const layout = await readFile(resolve(join(Config.get('basedir'), 'templates/layout.mustache')), { encoding: 'utf-8' });
		const main = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/main/mobile/main.mustache' : 'templates/main/desktop/main.mustache'
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
			content: mustache.render(main, {
				settingsHeader: this.i18n("settingsHeader"),
				botUsername:   `${BOT_USERNAME}`,
				authUrl:       `https://${BACKEND_URL}/tg_auth`,
			}, {
				sidebar:       this.sidebar,
				authorization: this.authorization,
			}),
		});
	}
}

export default main;
