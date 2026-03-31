import type {Builder} from './builder'
import type { Translation, TranslationPreview } from '../api/models';
import Config from 'config';
import mustache from 'mustache';
import { readFileSync } from 'node:fs';
import { resolve, join } from 'node:path';
import AbstractBuilder from './abstractBuilder'

class PublicTranslationBuilder extends AbstractBuilder {
	build() {
		return ""
	}
}

export default PublicTranslationBuilder