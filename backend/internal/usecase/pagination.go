package usecase

// normalizePagination clamps limit and offset to valid ranges.
func normalizePagination(limit, offset int32) (int32, int32) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}

// newPaginationResponse builds a PaginationResponse from total count and normalized limit/offset.
func newPaginationResponse(total int64, limit, offset int32) PaginationResponse {
	return PaginationResponse{
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		HasMore: int64(offset+limit) < total,
	}
}
