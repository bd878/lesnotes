const empty = {
	ID: 0,
	login: "",
}

export default function mapUserFromProto(user) {
	if (!user)
		return empty

	return {
		ID:    user.id,
		login: user.login,
	}
}