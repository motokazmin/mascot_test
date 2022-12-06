package config

type (
	Config struct {
		HttpService HttpService `json:"http"`
	}

	HttpService struct {
		Listen         string   `json:"listen"`
		Path           string   `json:"path"`
		HttpMode       string   `json:"mode"`
		AllowedHeaders []string `json:"allowed_headers,omitempty"`
	}
)
