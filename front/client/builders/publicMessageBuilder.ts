import type { Message } from '../api/models';
import type { FileWithMime } from '../types';
import Config from 'config';
import mustache from 'mustache';
import api from '../api';
import * as is from '../third_party/is';
import { readFile } from 'node:fs/promises';
import { resolve, join } from 'node:path';
import AbstractBuilder from './abstractBuilder'

class PublicMessageBuilder extends AbstractBuilder {
	sidebar   = undefined;
	filesView = undefined;

	async addSidebar() {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/sidebar_vertical/mobile/sidebar_vertical.mustache' : 'templates/sidebar_vertical/desktop/sidebar_vertical.mustache'
		)), { encoding: 'utf-8' });

		this.sidebar = mustache.render(template, {
			settingsHeader: this.i18n("settingsHeader")
		}, {
			settings:       this.settings,
		})
	}

	async addFilesView(files: FileWithMime[]) {
		const template = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/message/mobile/files_view.mustache' : 'templates/message/desktop/files_view.mustache'
		)), { encoding: 'utf-8' });

		this.filesView = mustache.render(template, {
			files:    files,
			imgSrc:   function() { return `/files/v1/read/${this.name}` },
			fileHref: function() { return `/files/v1/download?id=${this.ID}` },
		})
	}

	async build(message?: Message) {
		const styles = await readFile(resolve(join(Config.get('basedir'), 'public/styles/styles.css')), { encoding: 'utf-8' });
		const layout = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/layout/mobile/layout.mustache' : 'templates/layout/desktop/layout.mustache'
		)), { encoding: 'utf-8' });
		const content = await readFile(resolve(join(Config.get('basedir'),
			this.isMobile ? 'templates/message/mobile/message.mustache' : 'templates/message/desktop/message.mustache'
		)), { encoding: 'utf-8' });

		const theme = this.theme
		const fontSize = this.fontSize

		return mustache.render(layout, {
			html: () => (text, render) => {
				let html = "<html"

				if (theme) html += ` class="${theme}"`;
				if (this.lang) html += ` lang="${this.lang}"`;
				if (fontSize) html += ` data-size="${fontSize}"`
				html += ">"

				return html + render(text) + "</html>"
			},
			manifest:  "/public/manifest.json",
			styles:    styles,
			lang:      this.lang,
			theme:     theme,
		}, {
			footer:    this.footer,
			content:   mustache.render(content, {
				title:  message.title,
				text:   message.text,
				files:  message.files,
			}, {
				settings:  this.settings,
				sidebar:   this.sidebar,
				filesView: this.filesView,
			})
		})
	}
}

export default PublicMessageBuilder
