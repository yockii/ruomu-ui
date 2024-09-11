package controller

import (
	"encoding/json"
	"errors"
	"github.com/tidwall/gjson"

	logger "github.com/sirupsen/logrus"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/server"
	"github.com/yockii/ruomu-core/util"
	"gorm.io/gorm"

	"github.com/yockii/ruomu-ui/model"
)

var PageController = new(pageController)

type pageController struct{}

func (_ *pageController) Add(value []byte) (any, error) {
	instance := new(model.Page)
	if err := json.Unmarshal(value, instance); err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		}, nil
	}

	// 处理必填
	if instance.Name == "" || instance.ProjectID == 0 {
		return &server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " name / project id",
		}, nil
	}

	var c int64
	if err := database.DB.Model(&model.Page{}).Where(&model.Page{ProjectID: instance.ProjectID, Route: instance.Route}).Count(&c).Error; err != nil {
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

func (_ *pageController) Update(value []byte) (any, error) {
	instance := new(model.Page)
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

	if err := database.DB.Model(&model.Page{ID: instance.ID}).Updates(&model.Page{
		ProjectID: instance.ProjectID,
		Name:      instance.Name,
		ParentID:  instance.ParentID,
		Route:     instance.Route,
		Schema:    instance.Schema,
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

func (_ *pageController) Delete(value []byte) (any, error) {
	instance := new(model.Page)
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

func (_ *pageController) List(value []byte) (any, error) {
	instance := new(model.Page)
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

	condition := &model.Page{
		ID:        instance.ID,
		ProjectID: instance.ProjectID,
		ParentID:  instance.ParentID,
	}
	if instance.Name != "" {
		tx.Where("name like ?", "%"+instance.Name+"%")
		instance.Name = ""
	}
	if instance.Route != "" {
		tx.Where("route like ?", "%"+instance.Route+"%")
		instance.Route = ""
	}

	var list []*model.Page
	var total int64
	err := tx.Omit("schema").Find(&list, condition).Offset(-1).Count(&total).Error
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

func (_ *pageController) Instance(value []byte) (any, error) {
	instance := new(model.Page)
	if err := json.Unmarshal(value, instance); err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		}, nil
	}
	if err := database.DB.Where(instance).Take(instance).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &server.CommonResponse{}, nil
		}
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

func (_ *pageController) Schema(value []byte) (any, error) {
	instance := new(model.Page)
	if err := json.Unmarshal(value, instance); err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		}, nil
	}
	if err := database.DB.Where(instance).Take(instance).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &server.CommonResponse{}, nil
		}
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		}, nil
	}

	j := gjson.Parse(instance.Schema)

	return &server.CommonResponse{
		Data: j.Map(),
	}, nil
}
