package structs

type (
	Request struct {
		Page    string `json:"page" query:"page"`
		Limit   string `json:"limit" query:"limit"`
		OrderBy string `json:"orderBy" query:"orderBy"`
		SortBy  string `json:"sortBy" query:"sortBy"`
		Draw    string `json:"draw" query:"draw"`
		Search  string `json:"search" query:"search"`
	}

	RequestV2 struct {
		Offset  string `json:"offset" query:"offset"`
		Limit   string `json:"limit" query:"limit"`
		OrderBy string `json:"orderBy" query:"orderBy"`
		SortBy  string `json:"sortBy" query:"sortBy"`
		Draw    string `json:"draw" query:"draw"`
		Search  string `json:"search" query:"search"`
	}
)
