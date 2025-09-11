export interface User {
	ID:         number;
	login:      string;
	theme:      string;
}

const empty: User = {
	ID:    0,
	login: "",
	theme: "light",
}

export default function mapUserFromProto(user?: any): User {
	if (!user)
		return empty

	return {
		ID:    user.id,
		login: user.login,
		theme: user.theme,
	}
}