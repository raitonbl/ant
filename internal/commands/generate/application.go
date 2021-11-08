package generate

import (
	"github.com/raitonbl/ant/internal"
	golang "github.com/raitonbl/ant/internal/commands/generate/sdk/golang"
	python3 "github.com/raitonbl/ant/internal/commands/generate/sdk/python3"
)

func doGenerateCLIProject(context internal.GenerateContext) (string, error) {

	if context.GetProjectTargetLanguage() == internal.GoLang {
		return golang.GenerateProject(context)
	} else if context.GetProjectTargetLanguage() == internal.Python3 {
		return python3.GenerateProject(context)
	} else {
		return "", internal.GetProblemFactory().GetUnsupportedLanguage(string(internal.ApplicationType), context.GetTargetProjectLocation())
	}

}
