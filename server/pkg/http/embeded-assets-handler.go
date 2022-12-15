package http

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"strings"
	"text/template"
)

// FIXME: What is the path for the embed?

var (
	//go:embed dist/views/*.html
	embededViews embed.FS

	//go:embed dist/assets
	embededStaticAssets embed.FS
)

type EmbedadAssetsHandler struct {
	templateMap map[string]*template.Template
}

func CreateEmbededAssetsHandler() *EmbedadAssetsHandler {
	viewFiles, err := embededViews.ReadDir("/dist/views")
	if err != nil {
		panic(fmt.Errorf("EmbededAssetsHandler: Can not access embeded views: %w", err))
	}

	eah := new(EmbedadAssetsHandler)
	eah.templateMap = make(map[string]*template.Template, len(viewFiles))

	for _, view := range viewFiles {
		viewNameArray := strings.Split(strings.ToLower(view.Name()), ".")
		if len(viewNameArray) != 2 {
			panic("EmbededAssetsHandler: Invalid embeded view name")
		}

		viewName := viewNameArray[0]
		if _, viewKeyExists := eah.templateMap[viewName]; viewKeyExists {
			panic("EmbededAssetsHandler: View associated to the given name already exists")
		}

		templateContent, err := embededViews.ReadFile(view.Name())
		if err != nil {
			panic(fmt.Errorf("EmbededAssetsHandler: Can not access view file content: %w", err))
		}

		eah.templateMap[viewName] = template.Must(template.New(viewName).Parse(string(templateContent)))
	}

	return eah
}

func (eah *EmbedadAssetsHandler) GetEmbededViewTemplate(viewTemplateName string) *template.Template {
	viewTemplate, viewKeyExists := eah.templateMap[viewTemplateName]
	if !viewKeyExists {
		panic("EmbededAssetsHandler: No view found associated to given key")
	}

	return viewTemplate
}

func (eah *EmbedadAssetsHandler) GetEmbededStaticAssetsHandler() http.Handler {
	staticAssetsDir, err := fs.Sub(embededStaticAssets, "dist/assets")
	if err != nil {
		panic(fmt.Errorf("EmbededAssetsHandler: Can not access embeded assets: %w", err))
	}

	return http.FileServer(http.FS(staticAssetsDir))
}
