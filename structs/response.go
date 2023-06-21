package structs

import "io"

type (
	Paginator struct {
		Data            interface{} `json:"data"`
		Draw            int         `json:"draw"`
		RecordsFiltered int64       `json:"recordsFiltered"`
		RecordsTotal    int64       `json:"recordsTotal"`
		HasMore         bool        `json:"hasMore"`
	}

	Response struct {
		Code        string      `json:"responseCode"`
		Description string      `json:"responseDescription"`
		Data        interface{} `json:"data"`
	}

	HTTPResponse struct {
		Body       Response
		StatusCode int
		Headers    map[string][]string
		Error      error
	}

	SurroundingHTTPResponse struct {
		Body       interface{}
		StatusCode int
		Headers    map[string][]string
		Error      error
	}

	HTTPFileReader struct {
		Param    string
		FileName string
		Reader   io.Reader
	}
)
