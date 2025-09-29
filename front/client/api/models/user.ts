export interface User {
	ID:         number;
	login:      string;
	theme:      string;
	lang:       string;
	fontSize:   number;
}

const empty: User = {
	ID:    0,
	login: "",
	theme: "light",
	lang:  "",
	fontSize: 0,
}

export default function mapUserFromProto(user?: any): User {
	if (!user)
		return empty

	return {
		ID:        user.id,
		login:     user.login,
		theme:     user.theme,
		lang:      user.language,
		fontSize:  user.font_size,
	}
}