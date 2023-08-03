package epr

// means "error payload response"

type errorJson struct {
	Message string `json:"message"`
}

func APIError(err string) errorJson {
	return errorJson{Message: err}
}
