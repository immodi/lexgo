package framework

import "fmt"

type AppEnv int

const (
	ENVDev AppEnv = iota
	ENVProd
	ENVNil
)

func ParseAppEnv(s string) (AppEnv, error) {
	switch s {
	case "dev":
		return ENVDev, nil
	case "production":
		return ENVProd, nil
	default:
		return ENVNil, fmt.Errorf(
			"invalid 'config.env' value %q (allowed: dev, production)",
			s,
		)
	}
}

func (e AppEnv) String() string {
	switch e {
	case ENVDev:
		return "dev"
	case ENVProd:
		return "production"
	default:
		return "nil"
	}
}
