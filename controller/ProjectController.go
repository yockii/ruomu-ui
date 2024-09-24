package controller

import (
	"encoding/json"
	logger "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/server"
	"github.com/yockii/ruomu-core/util"
	"github.com/yockii/ruomu-ui/domain"
	"github.com/yockii/ruomu-ui/model"
)

var ProjectController = new(projectController)

type projectController struct{}

func (_ *projectController) Add(value []byte) (any, error) {
	instance := new(model.Project)
	if err := json.Unmarshal(value, instance); err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		}, nil
	}

	// 处理必填
	if instance.Name == "" {
		return &server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " name",
		}, nil
	}

	var c int64
	if err := database.DB.Model(&model.Project{}).Where(&model.Project{Name: instance.Name}).Count(&c).Error; err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		}, nil
	} else if c > 0 {
		return &server.CommonResponse{
			Code: server.ResponseCodeDuplicated,
			Msg:  server.ResponseMsgDuplicated,
		}, nil
	}

	instance.ID = util.SnowflakeId()
	if err := database.DB.Create(instance).Error; err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		}, nil
	}
	return &server.CommonResponse{
		Data: instance,
	}, nil
}

func (_ *projectController) UpdateFrontend(value []byte) (any, error) {
	instance := new(model.ProjectFrontend)
	if err := json.Unmarshal(value, instance); err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		}, nil
	}
	// 处理必填
	if instance.ID == 0 {
		return &server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " id",
		}, nil
	}

	// 如果没有就新增，如果有就更新
	// 先检查有没有对应的project
	project := new(model.Project)
	if err := database.DB.Where(&model.Project{ID: instance.ID}).First(project).Error; err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		}, nil
	}

	var c int64
	if err := database.DB.Model(&model.ProjectFrontend{}).Where(&model.ProjectFrontend{ID: instance.ID}).Count(&c).Error; err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		}, nil
	}
	if c == 0 {
		if err := database.DB.Create(instance).Error; err != nil {
			logger.Errorln(err)
			return &server.CommonResponse{
				Code: server.ResponseCodeDatabase,
				Msg:  server.ResponseMsgDatabase + err.Error(),
			}, nil
		}
	} else {
		if err := database.DB.Model(&model.ProjectFrontend{ID: instance.ID}).Updates(&model.ProjectFrontend{
			ApiJson:   instance.ApiJson,
			CssJson:   instance.CssJson,
			JsJson:    instance.JsJson,
			StoreJson: instance.StoreJson,
		}).Error; err != nil {
			logger.Errorln(err)
			return &server.CommonResponse{
				Code: server.ResponseCodeDatabase,
				Msg:  server.ResponseMsgDatabase + err.Error(),
			}, nil
		}
	}

	return &server.CommonResponse{
		Data: true,
	}, nil

}

func (_ *projectController) Update(value []byte) (any, error) {
	instance := new(model.Project)
	if err := json.Unmarshal(value, instance); err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		}, nil
	}
	// 处理必填
	if instance.ID == 0 {
		return &server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " id",
		}, nil
	}

	if err := database.DB.Model(&model.Project{ID: instance.ID}).Updates(&model.Project{
		Name:        instance.Name,
		Description: instance.Description,
		HomePageID:  instance.HomePageID,
		Status:      instance.Status,
	}).Error; err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		}, nil
	}
	return &server.CommonResponse{
		Data: true,
	}, nil
}

func (_ *projectController) Delete(value []byte) (any, error) {
	instance := new(model.Project)
	if err := json.Unmarshal(value, instance); err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		}, nil
	}
	// 处理必填
	if instance.ID == 0 {
		return &server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " id",
		}, nil
	}

	if err := database.DB.Where(instance).Delete(instance).Error; err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		}, nil
	}
	return &server.CommonResponse{
		Data: true,
	}, nil
}

func (_ *projectController) List(value []byte) (any, error) {
	instance := new(model.Project)
	if err := json.Unmarshal(value, instance); err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		}, nil
	}

	vj := gjson.ParseBytes(value)
	paginate := new(server.Paginate)
	if vj.Get("limit").Exists() {
		paginate.Limit = int(vj.Get("limit").Int())
	}
	if vj.Get("offset").Exists() {
		paginate.Offset = int(vj.Get("offset").Int())
	}

	if paginate.Limit <= 0 && paginate.Limit != -1 {
		paginate.Limit = 10
	}

	tx := database.DB.Limit(paginate.Limit).Offset(paginate.Offset)

	condition := &model.Project{
		ID:          instance.ID,
		Name:        "",
		Description: "",
		HomePageID:  instance.HomePageID,
		Status:      instance.Status,
	}
	if instance.Name != "" {
		tx.Where("name like ?", "%"+instance.Name+"%")
		instance.Name = ""
	}

	var list []*model.Project
	var total int64
	err := tx.Omit("store_json").Find(&list, condition).Offset(-1).Count(&total).Error
	if err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		}, nil
	}
	return &server.CommonResponse{
		Data: &server.Paginate{
			Total:  total,
			Offset: paginate.Offset,
			Limit:  paginate.Limit,
			Items:  list,
		},
	}, nil
}

func (_ *projectController) Instance(value []byte) (any, error) {
	instance := new(model.Project)
	if err := json.Unmarshal(value, instance); err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		}, nil
	}
	// 处理必填
	if instance.ID == 0 {
		return &server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " id",
		}, nil
	}
	if err := database.DB.Where(instance).Take(instance).Error; err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		}, nil
	}

	result := &domain.Project{
		Project: *instance,
	}

	// 查询pf
	pf := new(model.ProjectFrontend)
	if err := database.DB.Where(&model.ProjectFrontend{ID: instance.ID}).Take(pf).Error; err != nil {
		logger.Errorln(err)
	}

	if pf.StoreJson != "" {
		if err := json.Unmarshal([]byte(pf.StoreJson), &result.Store); err != nil {
			logger.Errorln(err)
		}
	}
	//if pf.JsJson != "" {
	//	if err := json.Unmarshal([]byte(pf.JsJson), &result.Js); err != nil {
	//		logger.Errorln(err)
	//	}
	//}
	//if pf.CssJson != "" {
	//	if err := json.Unmarshal([]byte(pf.CssJson), &result.Css); err != nil {
	//		logger.Errorln(err)
	//	}
	//}
	//if pf.ApiJson != "" {
	//	if err := json.Unmarshal([]byte(pf.ApiJson), &result.Api); err != nil {
	//		logger.Errorln(err)
	//	}
	//}

	return &server.CommonResponse{
		Data: result,
	}, nil
}
