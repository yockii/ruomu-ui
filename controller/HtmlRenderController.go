package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/ruomu-core/config"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-ui/model"
	"html/template"
	"path/filepath"
)

var HtmlRenderController = new(htmlRenderController)

type htmlRenderController struct {
	templates map[string]*template.Template
}

func JsVar(name string) template.JS {
	return template.JS(name)
}

func parseHtmlTemplate(path string) (*template.Template, error) {
	t := template.New(filepath.Base(path))

	t = t.Funcs(template.FuncMap{
		"jsVar": JsVar,
	})

	var err error
	t, err = t.ParseFiles(path)
	if err != nil {
		logger.Errorln(err)
		return nil, err
	}

	return t, nil
}

func init() {
	HtmlRenderController.templates = make(map[string]*template.Template)
	HtmlRenderController.templates["index"] = template.Must(parseHtmlTemplate("./views/index.html"))
	HtmlRenderController.templates["error"] = template.Must(parseHtmlTemplate("./views/error.html"))
}

func (c *htmlRenderController) Index(value []byte) (any, error) {
	instance := new(model.Project)
	if err := json.Unmarshal(value, instance); err != nil {
		logger.Errorln(err)
		return c.renderError(err)
	}

	if instance.ID == 0 {
		instance.ID = config.GetUint64("project.id")
	}

	// 查询下项目用到的库
	var pmlList []*model.ProjectMaterialLibVersion
	if err := database.DB.Where(&model.ProjectMaterialLibVersion{ProjectID: instance.ID}).Find(&pmlList).Error; err != nil {
		logger.Errorln(err)
		return c.renderError(err)
	}

	var mlvList []*model.MaterialLibVersion

	if len(pmlList) > 0 {
		for _, pml := range pmlList {
			mlv := new(model.MaterialLibVersion)
			if err := database.DB.Where(&model.MaterialLibVersion{ID: pml.LibVersionID}).First(&mlv).Error; err != nil {
				logger.Errorln(err)
				return c.renderError(err)
			}
			mlvList = append(mlvList, mlv)
		}
	}

	buf := new(bytes.Buffer)
	var temp *template.Template
	var err error
	if logger.IsLevelEnabled(logger.DebugLevel) {
		temp, err = parseHtmlTemplate("./views/index.html")
		if err != nil {
			logger.Errorln(err)
			return c.renderError(err)
		}
	} else {
		temp = c.templates["index"]
	}
	err = temp.Execute(buf, map[string]any{
		"instance": instance,
		"libList":  mlvList,
	})
	return buf.Bytes(), err
}
func (c *htmlRenderController) Canvas(value []byte) (any, error) {
	instance := new(model.Project)
	if err := json.Unmarshal(value, instance); err != nil {
		logger.Errorln(err)
		return c.renderError(err)
	}

	if instance.ID == 0 {
		logger.Warn("project id is empty")
		return c.renderError(errors.New("project id is empty"))
	}

	// 查询下项目用到的库
	var pmlList []*model.ProjectMaterialLibVersion
	if err := database.DB.Where(&model.ProjectMaterialLibVersion{ProjectID: instance.ID}).Find(&pmlList).Error; err != nil {
		logger.Errorln(err)
		return c.renderError(err)
	}

	var mlvList []*model.MaterialLibVersion

	if len(pmlList) > 0 {
		for _, pml := range pmlList {
			mlv := new(model.MaterialLibVersion)
			if err := database.DB.Where(&model.MaterialLibVersion{ID: pml.LibVersionID}).First(&mlv).Error; err != nil {
				logger.Errorln(err)
				return c.renderError(err)
			}
			mlvList = append(mlvList, mlv)
		}
	}

	buf := new(bytes.Buffer)
	var temp *template.Template
	var err error
	if logger.IsLevelEnabled(logger.DebugLevel) {
		temp, err = parseHtmlTemplate("./views/canvas.html")
		if err != nil {
			logger.Errorln(err)
			return c.renderError(err)
		}
	} else {
		temp = c.templates["index"]
	}

	err = temp.Execute(buf, map[string]any{
		"instance": instance,
		"libList":  mlvList,
	})
	return buf.Bytes(), err
}

func (c *htmlRenderController) renderError(err error) (any, error) {
	buf := new(bytes.Buffer)

	err = c.templates["error"].Execute(buf, map[string]any{
		"error": err.Error(),
	})
	return buf.Bytes(), err
}
