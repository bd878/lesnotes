import Config from 'config';
import mustache from 'mustache';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import AbstractPublicBuilder from './abstractPublicBuilder';

let threadTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/thread/desktop/thread.mustache')), { encoding: 'utf-8' });
let threadTemplateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/thread/mobile/thread.mustache')), { encoding: 'utf-8' });

class PublicThreadBuilder extends AbstractPublicBuilder {

	build() {
		return mustache.render(this.isMobile ? threadTemplate : threadTemplateMobile, {
		}, {
			signup:           this.signup,
			logout:           this.logout,
			sidebar:          this.sidebar,
			searchForm:       this.searchForm,
			messagesList:     this.messagesList,
		})
	}

}

export default PublicThreadBuilder
