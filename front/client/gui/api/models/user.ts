export default function mapUserFromProto(user) {
	if (!user) {
		return {
			ID: 0,
			name: "",
		}
	}

	return {
		ID: user.id,
		name: user.name,
	}
}