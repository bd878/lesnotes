import type {Builder} from './builder'
import Config from 'config';
import mustache from 'mustache';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import AbstractBuilder from './abstractBuilder';

let controlPanelTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/home/desktop/control_panel.mustache')), { encoding: 'utf-8' });
let controlPanelTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/home/mobile/control_panel.mustache')), { encoding: 'utf-8' });

class ControlPanelBuilder extends AbstractBuilder {
	auth = undefined

	addAuth(auth: Builder) {
		this.auth = auth.build()
		return this
	}

	build() {
		return mustache.render(this.isMobile ? controlPanelTemplateMobile : controlPanelTemplate, {}, {
			logout:           this.auth,
		})
	}
}

export default ControlPanelBuilder
