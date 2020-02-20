package typredis

import (
	"github.com/typical-go/typical-go/pkg/typbuild"
	"github.com/typical-go/typical-go/pkg/typcfg"
	"github.com/urfave/cli/v2"
)

// Redis of Redis
type Redis struct {
	host        string
	port        string
	password    string
	dockerName  string
	dockerImage string
	prefix      string
}

// New Redis Redis
func New() *Redis {
	return &Redis{
		host:        "localhost",
		port:        "6379",
		password:    "redispass",
		dockerImage: "redis:4.0.5-alpine",
		dockerName:  "redis",
		prefix:      "REDIS",
	}
}

// Withhost to return module with new host
func (m *Redis) Withhost(host string) *Redis {
	m.host = host
	return m
}

// Withport to return module with new port
func (m *Redis) Withport(port string) *Redis {
	m.port = port
	return m
}

// Withpassword to return module with new password
func (m *Redis) Withpassword(password string) *Redis {
	m.password = password
	return m
}

// WithdockerImage to return module with new docker image
func (m *Redis) WithdockerImage(dockerImage string) *Redis {
	m.dockerImage = dockerImage
	return m
}

// WithdockerName to return module with new docker name
func (m *Redis) WithdockerName(dockerName string) *Redis {
	m.dockerName = dockerName
	return m
}

// WithPrefix to return module with new prefix
func (m *Redis) WithPrefix(prefix string) *Redis {
	m.prefix = prefix
	return m
}

// BuildCommands of module
func (m *Redis) BuildCommands(c *typbuild.Context) []*cli.Command {
	return []*cli.Command{
		{
			Name:  "redis",
			Usage: "Redis Tool",
			Before: func(ctx *cli.Context) error {
				return typcfg.LoadEnvFile()
			},
			Subcommands: []*cli.Command{
				{
					Name:    "console",
					Aliases: []string{"c"},
					Usage:   "Redis Interactive",
					Action:  c.ActionFunc(m.console),
				},
			},
		},
	}
}

// Configure Redis
func (m *Redis) Configure(loader typcfg.Loader) *typcfg.Detail {
	return &typcfg.Detail{
		Prefix: m.prefix,
		Spec: &Config{
			Host:     m.host,
			Port:     m.port,
			Password: m.password,
		},
		Constructor: func() (cfg Config, err error) {
			err = loader.Load(m.prefix, &cfg)
			return
		},
	}
}

// Provide dependencies
func (m *Redis) Provide() []interface{} {
	return []interface{}{
		m.connect,
	}
}

// Prepare the module
func (m *Redis) Prepare() []interface{} {
	return []interface{}{
		m.ping,
	}
}

// Destroy dependencies
func (m *Redis) Destroy() []interface{} {
	return []interface{}{
		m.disconnect,
	}
}