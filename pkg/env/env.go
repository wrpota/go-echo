package env

import (
	"fmt"
	"os"
	"strings"
)

var (
	active Environment
	dev    Environment = &environment{value: "dev"}
	fat    Environment = &environment{value: "fat"}
	pro    Environment = &environment{value: "pro"}
)

var _ Environment = (*environment)(nil)

// Environment 环境配置
type Environment interface {
	Value() string
	IsDev() bool
	IsPro() bool
	i()
}

type environment struct {
	value string
}

func (e *environment) Value() string {
	return e.value
}

func (e *environment) IsDev() bool {
	return e.value == "dev"
}

func (e *environment) IsFat() bool {
	return e.value == "fat"
}

func (e *environment) IsPro() bool {
	return e.value == "pro"
}

func (e *environment) i() {}

func init() {
	env := os.Getenv("GO_ENV")
	switch strings.ToLower(strings.TrimSpace(env)) {
	case "dev":
		active = dev
	case "pro":
		active = pro
	case "fat":
		active = fat
	default:
		active = dev
		fmt.Println("Warning: The environment variable 'GO_ENV' cannot be found, or it is illegal. The default 'dev' will be used.")
	}
}

func Active() Environment {
	return active
}
