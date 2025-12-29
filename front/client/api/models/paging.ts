export interface ExternalPaging {
	is_last_page:    boolean;
	is_first_page:   boolean;
	total:           number;
	count:           number;
	offset:          number;
}

export interface Paging {
	isLastPage:   boolean;
	isFirstPage:  boolean;
	total:        number;
	count:        number;
	offset:       number;
}

const EmptyPaging: Paging = Object.freeze({
	isLastPage:   true,
	isFirstPage:  true,
	total:        0,
	count:        0,
	offset:       0,
})

export default function mapPagingFromProto(paging?: ExternalPaging): Paging {
	if (!paging) {
		return EmptyPaging
	}

	const res = {
		isLastPage:  paging.is_last_page,
		isFirstPage: paging.is_first_page,
		total:       paging.total,
		count:       paging.count,
		offset:      paging.offset,
	}

	return res
}

export { EmptyPaging }
