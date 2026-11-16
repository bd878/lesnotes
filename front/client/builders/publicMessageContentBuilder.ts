import Config from 'config';
import mustache from 'mustache';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import ContentBuilder from './contentBuilder'

let template = readFileSync(resolve(join(Config.get('basedir'),'templates/content/desktop/public_message_content.mustache')), { encoding: 'utf-8' });
let mobileTemplate = readFileSync(resolve(join(Config.get('basedir'),'templates/content/mobile/public_message_content.mustache')), { encoding: 'utf-8' });

class PublicMessageContentBuilder extends ContentBuilder {
	build() {
		return mustache.render(this.isMobile ? mobileTemplate : template, {
			hasFeatures: this.messageFeatures != undefined,
		}, {
			messageView:       this.messageView,
			messageFeatures:   this.messageFeatures,
			header:            this.header,
			messageHeader:     this.messageHeader,
			threadPath:        this.threadPath,
		})
	}

}

export default PublicMessageContentBuilder
