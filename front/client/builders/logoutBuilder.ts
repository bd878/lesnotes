import type {Builder} from './builder'
import Config from 'config';
import mustache from 'mustache';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import AbstractBuilder from './abstractBuilder';

let logoutTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/logout/desktop/logout.mustache')), { encoding: 'utf-8' });
let logoutTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/logout/mobile/logout.mustache')), { encoding: 'utf-8' });

class LogoutBuilder extends AbstractBuilder {
	build() {
		const search = this.search

		return mustache.render(this.isMobile ? logoutTemplateMobile : logoutTemplate, {
			logout:           this.i18n("logout"),
			logoutHref:       function() {
				const params = new URLSearchParams(search);
				params.delete("cwd");
				params.delete("id"); /* TODO: delete pagination */
				return "/logout?" + params.toString()
			},
		})
	}
}

export default LogoutBuilder