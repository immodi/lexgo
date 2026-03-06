package framework

type AppData struct {
	Port           int32
	Env            AppEnv
	AllowedOrigins []string
}

func (d *AppData) IsProduction() bool {
	return d.Env == ENVProd
}

func (d *AppData) GetAllowedOrigin(requestOrigin string) string {

	const (
		DEFAULT_ALLOWED_ORIGIN_DEV  = "*"
		DEFAULT_ALLOWED_ORIGIN_PROD = ""
	)

	allowed := DEFAULT_ALLOWED_ORIGIN_DEV
	if d.IsProduction() {
		allowed = DEFAULT_ALLOWED_ORIGIN_PROD
	}

	for _, origin := range d.AllowedOrigins {
		if requestOrigin == origin {
			return origin
		}
	}

	return allowed
}
