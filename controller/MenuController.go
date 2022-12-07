package controller

import (
	"encoding/json"

	logger "github.com/sirupsen/logrus"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/server"
	"github.com/yockii/ruomu-core/util"

	"github.com/yockii/ruomu-ui/model"
)

var MenuController = new(menuController)

type menuController struct{}

func (_ *menuController) Add(value []byte) (any, error) {
	instance := new(model.Menu)
	if err := json.Unmarshal(value, instance); err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		}, nil
	}

	// 处理必填
	if instance.Name == "" || instance.Code == "" {
		return &server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " name / code",
		}, nil
	}

	if c, err := database.DB.Count(&model.Menu{Code: instance.Code}); err != nil {
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

	instance.Id = util.SnowflakeId()
	if _, err := database.DB.Insert(instance); err != nil {
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

func (_ *menuController) Update(value []byte) (any, error) {
	instance := new(model.Menu)
	if err := json.Unmarshal(value, instance); err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		}, nil
	}
	// 处理必填
	if instance.Id == 0 {
		return &server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " id",
		}, nil
	}

	if _, err := database.DB.Update(&model.Menu{
		ParentId: instance.ParentId,
		Name:     instance.Name,
		Code:     instance.Code,
		Icon:     instance.Icon,
		PageCode: instance.PageCode,
	}, &model.Menu{Id: instance.Id}); err != nil {
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

func (_ *menuController) Delete(value []byte) (any, error) {
	instance := new(model.Menu)
	if err := json.Unmarshal(value, instance); err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		}, nil
	}
	// 处理必填
	if instance.Id == 0 {
		return &server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " id",
		}, nil
	}

	if _, err := database.DB.Delete(instance); err != nil {
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

func (_ *menuController) List(value []byte) (any, error) {
	instance := new(model.Menu)
	if err := json.Unmarshal(value, instance); err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		}, nil
	}
	paginate := new(server.Paginate)
	if err := json.Unmarshal(value, instance); err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		}, nil
	}
	if paginate.Limit <= 0 {
		paginate.Limit = 10
	}

	session := database.DB.NewSession().Limit(paginate.Limit, paginate.Offset)

	condition := &model.Menu{
		Id:       instance.Id,
		ParentId: instance.ParentId,
		Name:     "",
		Code:     "",
		Icon:     "",
		PageCode: "",
	}
	if instance.Name != "" {
		session.Where("name like ?", "%"+instance.Name+"%")
		instance.Name = ""
	}
	if instance.Code != "" {
		session.Where("code like ?", "%"+instance.Code+"%")
		instance.Code = ""
	}
	if instance.PageCode != "" {
		session.Where("page_code like ?", "%"+instance.PageCode+"%")
		instance.PageCode = ""
	}

	var list []*model.Menu
	total, err := session.FindAndCount(&list, condition)
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
