import AbstractBuilder from './abstractBuilder';

abstract class AbstractPublicBuilder extends AbstractBuilder {
	isAuthed: boolean = false;

	constructor(isAuthed: boolean = false, isMobile: boolean, lang: string = "en", theme: string = "light", fontSize: string = "medium", search: string = "", path: string = "") {
		super(isMobile, lang, theme, fontSize, search, path)

		this.isAuthed = isAuthed
	}
}

export default AbstractPublicBuilder
