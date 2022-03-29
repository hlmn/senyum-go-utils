package structs

type (
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
)
