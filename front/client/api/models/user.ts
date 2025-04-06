const empty = {
	ID: 0,
	name: "",
}
export default function mapUserFromProto(user) {
	if (!user)
		return empty

	return {
		ID: user.id,
		name: user.name,
	}
}