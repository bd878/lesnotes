export interface User {
	ID:         number;
	login:      string;
}

const empty: User = {
	ID:    0,
	login: "",
}

export default function mapUserFromProto(user): User {
	if (!user)
		return empty

	return {
		ID:    user.id,
		login: user.login,
	}
}