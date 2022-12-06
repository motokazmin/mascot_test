package config

type (
	Config struct {
		HttpService     HttpService     `json:"http"`
		PostgresService PostgresService `json:"db"`
	}

	HttpService struct {
		Listen         string   `json:"listen"`
		Path           string   `json:"path"`
		HttpMode       string   `json:"mode"`
		AllowedHeaders []string `json:"allowed_headers,omitempty"`
	}

	PostgresService struct {
		User     string `json:"user"`
		Dbname   string `json:"dbname"`
		Password string `json:"password"`
		SslMode  string `json:"sslmode"`
	}
)
