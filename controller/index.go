package controller

import (
	"encoding/json"

	"github.com/yockii/ruomu-ui/constant"
)

// Dispatch 注入点
func Dispatch(code string, headers map[string][]string, value []byte) ([]byte, error) {
	switch code {
	case constant.InjectCodeProjectAdd:
		return wrapCall(value, ProjectController.Add)
	case constant.InjectCodeProjectUpdate:
		return wrapCall(value, ProjectController.Update)
	case constant.InjectCodeProjectUpdateFrontend:
		return wrapCall(value, ProjectController.UpdateFrontend)
	case constant.InjectCodeProjectDelete:
		return wrapCall(value, ProjectController.Delete)
	case constant.InjectCodeProjectList:
		return wrapCall(value, ProjectController.List)
	case constant.InjectCodeProjectInstance:
		return wrapCall(value, ProjectController.Instance)

	case constant.InjectCodePageAdd:
		return wrapCall(value, PageController.Add)
	case constant.InjectCodePageUpdate:
		return wrapCall(value, PageController.Update)
	case constant.InjectCodePageDelete:
		return wrapCall(value, PageController.Delete)
	case constant.InjectCodePageInstance:
		return wrapCall(value, PageController.Instance)
	case constant.InjectCodePageList:
		return wrapCall(value, PageController.List)
	case constant.InjectCodePageSchema:
		return wrapCall(value, PageController.Schema)

	case constant.InjectCodeMaterialLibAdd:
		return wrapCall(value, MaterialLibController.Add)
	case constant.InjectCodeMaterialLibUpdate:
		return wrapCall(value, MaterialLibController.Update)
	case constant.InjectCodeMaterialLibDelete:
		return wrapCall(value, MaterialLibController.Delete)
	case constant.InjectCodeMaterialLibInstance:
		return wrapCall(value, MaterialLibController.Instance)
	case constant.InjectCodeMaterialLibList:
		return wrapCall(value, MaterialLibController.List)

	case constant.InjectCodeMaterialLibVersionAdd:
		return wrapCall(value, MaterialLibVersionController.Add)
	case constant.InjectCodeMaterialLibVersionUpdate:
		return wrapCall(value, MaterialLibVersionController.Update)
	case constant.InjectCodeMaterialLibVersionDelete:
		return wrapCall(value, MaterialLibVersionController.Delete)
	case constant.InjectCodeMaterialLibVersionInstance:
		return wrapCall(value, MaterialLibVersionController.Instance)
	case constant.InjectCodeMaterialLibVersionList:
		return wrapCall(value, MaterialLibVersionController.List)

	case constant.InjectCodeMaterialComponentGroupAdd:
		return wrapCall(value, MaterialComponentGroupController.Add)
	case constant.InjectCodeMaterialComponentGroupUpdate:
		return wrapCall(value, MaterialComponentGroupController.Update)
	case constant.InjectCodeMaterialComponentGroupDelete:
		return wrapCall(value, MaterialComponentGroupController.Delete)
	case constant.InjectCodeMaterialComponentGroupInstance:
		return wrapCall(value, MaterialComponentGroupController.Instance)
	case constant.InjectCodeMaterialComponentGroupList:
		return wrapCall(value, MaterialComponentGroupController.List)

	case constant.InjectCodeMaterialComponentAdd:
		return wrapCall(value, MaterialComponentController.Add)
	case constant.InjectCodeMaterialComponentUpdate:
		return wrapCall(value, MaterialComponentController.Update)
	case constant.InjectCodeMaterialComponentDelete:
		return wrapCall(value, MaterialComponentController.Delete)
	case constant.InjectCodeMaterialComponentInstance:
		return wrapCall(value, MaterialComponentController.Instance)
	case constant.InjectCodeMaterialComponentList:
		return wrapCall(value, MaterialComponentController.ListWithMetaInfo)

	case constant.InjectCodeProjectMaterialLibVersionAdd:
		return wrapCall(value, ProjectMaterialLibVersionController.Add)
	case constant.InjectCodeProjectMaterialLibVersionUpdate:
		return wrapCall(value, ProjectMaterialLibVersionController.Update)
	case constant.InjectCodeProjectMaterialLibVersionDelete:
		return wrapCall(value, ProjectMaterialLibVersionController.Delete)
	case constant.InjectCodeProjectMaterialLibVersionList:
		return wrapCall(value, ProjectMaterialLibVersionController.List)

	case constant.InjectCodeMaterialVersionComponentAdd:
		return wrapCall(value, MaterialComponentController.VersionAdd)
	case constant.InjectCodeMaterialVersionComponentDelete:
		return wrapCall(value, MaterialComponentController.VersionDelete)

	case constant.InjectCodeIndexHtml:
		return wrapHtmlCall(value, HtmlRenderController.Index)
	case constant.InjectCodeCanvasHtml:
		return wrapHtmlCall(value, HtmlRenderController.Canvas)
	}
	return nil, nil
}

func wrapCall(v []byte, f func([]byte) (any, error)) ([]byte, error) {
	r, err := f(v)
	if err != nil {
		return nil, err
	}
	bs, err := json.Marshal(r)
	return bs, err
}

func wrapHtmlCall(v []byte, f func([]byte) (any, error)) ([]byte, error) {
	r, err := f(v)
	return r.([]byte), err
}
