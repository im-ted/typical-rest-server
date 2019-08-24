package typicmd

import (
	"io/ioutil"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"gopkg.in/yaml.v2"
)

func dockerCompose(ctx *typictx.ActionContext) (err error) {
	log.Info("Generate docker-compose.yml")
	dockerCompose := ctx.DockerCompose()
	d1, _ := yaml.Marshal(dockerCompose)
	return ioutil.WriteFile("docker-compose.yml", d1, 0644)
}

func dockerUp(ctx *typictx.ActionContext) (err error) {
	if !ctx.Cli.Bool("no-compose") {
		err = dockerCompose(ctx)
		if err != nil {
			return
		}
	}
	cmd := exec.Command("docker-compose", "up", "--remove-orphans", "-d")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func dockerDown(ctx *typictx.ActionContext) (err error) {
	cmd := exec.Command("docker-compose", "down")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
