import type {Builder} from './builder'
import i18n from '../i18n';

abstract class AbstractBuilder implements Builder {
	isMobile:      boolean = false;
	lang:          string  = "en";
	search:        string = "";
	path:          string = "";
	theme:         string = "";
	fontSize:      string = "";

	constructor(isMobile: boolean, lang: string = "en", theme: string = "light", fontSize: string = "medium", search: string = "", path: string = "") {
		this.search = search
		this.isMobile = isMobile
		this.lang = lang
		this.path = path
		this.theme = theme
		this.fontSize = fontSize
	}

	i18n(key: string): string {
		return i18n(this.lang)(key)
	}

	abstract build();
}

export default AbstractBuilder
