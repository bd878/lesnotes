interface Builder {
	build(): string;
}

interface ScriptsBuilder extends Builder {
	scripts: string[]
}

export type { Builder }
export type { ScriptsBuilder }
