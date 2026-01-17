export interface User {
	ID:         number;
	login:      string;
	theme:      string;
	lang:       string;
	fontSize:   string;
}

const empty: User = {
	ID:    0,
	login: "",
	theme: "light",
	lang:  "",
	fontSize: "medium",
}

export default function mapUserFromProto(user?: any): User {
	if (!user) {
		return empty
	}

	let fontSize: string = "medium"
	if (user.font_size < 10) { fontSize = "small"; }
	else if (user.font_size > 14) { fontSize = "large"; }
	else { fontSize = "medium" }

	return {
		ID:        user.id,
		login:     user.login,
		theme:     user.theme,
		lang:      user.language,
		fontSize:  fontSize,
	}
}