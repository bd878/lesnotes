import Config from 'config';
import i18n from '../../i18n';
import mustache from 'mustache';
import { readFile } from 'node:fs/promises';
import { resolve, join } from 'node:path';
import Builder from '../builder';

async function main(ctx) {
	const builder = new MainBuilder(ctx.userAgent.isMobile, ctx.state.lang)

	await builder.addSettings(undefined, ctx.state.lang, ctx.state.theme, ctx.state.fontSize)
	await builder.addFooter()
	await builder.addSidebar()
	await builder.addAuthorization()

	ctx.body = await builder.build(ctx.state.theme)
	ctx.status = 200;
}

class MainBuilder extends Builder {
	sidebar = undefined;
	async addSidebar() {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/main/mobile/sidebar.mustache' : 'templates/main/desktop/sidebar.mustache'
		)), { encoding: 'utf-8' });

		this.sidebar = mustache.render(template)
	}

	authorization = undefined;
	async addAuthorization() {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/main/mobile/authorization.mustache' : 'templates/main/desktop/authorization.mustache'
		)), { encoding: 'utf-8' });

		this.authorization = mustache.render(template, {
			login:     this.i18n("login"),
			register:  this.i18n("register"),
		})
	}

	async build(theme?: string) {
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
				html += ">"

				return html + render(text) + "</html>"
			},
			scripts:  ["/public/pages/main/mainScript.js"],
			manifest: "/public/manifest.json",
			styles:   styles,
			lang:     this.lang,
			isMobile: this.isMobile ? "true" : "",
		}, {
			footer:  this.footer,
			content: mustache.render(main, {
				settingsHeader: this.i18n("settingsHeader"),
			}, {
				settings:  this.settings,
				sidebar:       this.sidebar,
				authorization: this.authorization,
			}),
		});
	}
}

export default main;
