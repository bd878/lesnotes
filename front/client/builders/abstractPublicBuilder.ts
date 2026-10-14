import AbstractBuilder from './abstractBuilder';

abstract class AbstractPublicBuilder extends AbstractBuilder {
	isAuthed: boolean = false;
	messageName: string = "";
	parentName: string = "";
	messageID: number = 0;

	constructor(isAuthed: boolean = false, parentName: string = "", messageName: string = "", isMobile: boolean, lang: string = "en",
		theme: string = "light", fontSize: string = "medium", search: string = "", path: string = "", messageID: number = 0) {
		super(isMobile, lang, theme, fontSize, search, path)

		this.isAuthed = isAuthed
		this.messageName = messageName
		this.parentName = parentName
		this.messageID = messageID
	}
}

export default AbstractPublicBuilder
