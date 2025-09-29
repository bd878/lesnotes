export interface Error {
	error:   boolean;
	status:  number;
	explain: string;
	code:    number;
	human:   string;
}

const empty: Error = {
	error:   false,
	status:  200,
	explain: "",
	code:    0,
	human:   "",
}

export default function mapErrorFromProto(error?: Error): Error {
	if (!error)
		return empty

	return error
}