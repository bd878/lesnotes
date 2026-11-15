import type {Builder} from './builder'
import Config from 'config';
import mustache from 'mustache';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import ContentBuilder from './contentBuilder';

let template = readFileSync(resolve(join(Config.get('basedir'),'templates/content/desktop/public_thread_content.mustache')), { encoding: 'utf-8' });
let templateMobile = readFileSync(resolve(join(Config.get('basedir'),'templates/content/mobile/public_thread_content.mustache')), { encoding: 'utf-8' });

class PublicThreadContentBuilder extends ContentBuilder {
	build() {
		return mustache.render(this.isMobile ? templateMobile : template, {
			hasContent:    this.threadView != undefined || this.messageView != undefined,
			hasFeatures:   this.messageFeatures != undefined,
		}, {
			header:        this.header,
			messagesTree:  this.messagesTree,
			messageFeatures: this.messageFeatures,
			controlPanel:  this.controlPanel,
			threadView:    this.threadView,
			messageView:   this.messageView,
		})
	}

}

export default PublicThreadContentBuilder
