import AbstractBuilder from './abstractBuilder';

abstract class AbstractPublicBuilder extends AbstractBuilder {
	isAuthed: boolean = false;
	messageName: string = "";
	parentName: string = "";

	constructor(isAuthed: boolean = false, parentName: string = "", messageName: string = "", isMobile: boolean, lang: string = "en",
		theme: string = "light", fontSize: string = "medium", search: string = "", path: string = "") {
		super(isMobile, lang, theme, fontSize, search, path)

		this.isAuthed = isAuthed
		this.messageName = messageName
		this.parentName = parentName
	}
}

export default AbstractPublicBuilder
