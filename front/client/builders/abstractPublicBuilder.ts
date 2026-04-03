import AbstractBuilder from './abstractBuilder';

abstract class AbstractPublicBuilder extends AbstractBuilder {
	isAuthed: boolean = false;
	messageName: string = "";
	threadName: string = "";

	constructor(isAuthed: boolean = false, threadName: string = "", messageName: string = "", isMobile: boolean, lang: string = "en",
		theme: string = "light", fontSize: string = "medium", search: string = "", path: string = "") {
		super(isMobile, lang, theme, fontSize, search, path)

		this.isAuthed = isAuthed
		this.messageName = messageName
		this.threadName = threadName
	}
}

export default AbstractPublicBuilder
