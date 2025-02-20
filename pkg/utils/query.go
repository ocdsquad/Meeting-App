package utils

type QueryPageLimit struct {
	Page    int
	Limit   int
	OrderBy string
	SortBy  string
}

type QueryCount struct {
	TotalData    int
	TotalContent int
}
