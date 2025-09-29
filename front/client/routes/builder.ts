import i18n from '../i18n';

abstract class Builder {
	isMobile: boolean = false;
	lang:     string  = "en";

	constructor(isMobile: boolean, lang: string = "en") {
		this.isMobile = isMobile
		this.lang = lang
	}

	i18n(key: string): string {
		return i18n(this.lang)(key)
	}

	abstract build();
}

export default Builder
