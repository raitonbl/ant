package golang

import (
	"embed"
	"github.com/raitonbl/ant/internal"
)

var (
	//go:embed app.mod.txt
	resources embed.FS
)

func GenerateProject(context internal.GenerateContext) (string, error) {
	var cfg Configuration
	cliObject, _ := context.GetDocument()
	err := context.BindConfiguration(&cfg)

	if err != nil {
		return "", internal.GetProblemFactory().GetProblem(err)
	}

	err = context.Write("app.mod.txt", makeGoMode(cliObject, cfg, true))

	if err != nil {
		return "", err
	}

	return "", nil
}

type Configuration struct {
	Name string `properties:"name"`
	Sdk  Sdk    `properties:"sdk"`
}

type Sdk struct {
	version string `properties:"version"`
}
