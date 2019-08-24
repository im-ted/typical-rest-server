package typical

import (
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/app"
	"github.com/typical-go/typical-rest-server/app/config"
	"github.com/typical-go/typical-rest-server/pkg/module/typpostgres"
	"github.com/typical-go/typical-rest-server/pkg/module/typserver"
)

// Context instance of Context
var Context = &typictx.Context{
	Name:        "Typical-RESTful-Server",
	Description: "Example of typical and scalable RESTful API Server for Go",
	Application: typictx.Application{
		StartFunc: app.Start,
		Config: typictx.Config{
			Prefix: "APP",
			Spec:   &config.Config{},
		},
		Initiations: []interface{}{
			app.Middlewares,
			app.Routes,
		},
	},
	Modules: []*typictx.Module{
		typserver.Module(),
		typpostgres.Module(),
	},
	Release: typictx.Release{
		Version: "0.6.6",
		GoOS:    []string{"linux", "darwin"},
		GoArch:  []string{"amd64"},
		Github: &typictx.Github{
			Owner:    "typical-go",
			RepoName: "typical-rest-server",
		},
		Versioning: typictx.Versioning{
			WithGitBranch:       true,
			WithLatestGitCommit: true,
		},
	},
}
