package controller

import (
	"encoding/json"

	logger "github.com/sirupsen/logrus"
	"github.com/yockii/ruomu-core/database"
	"github.com/yockii/ruomu-core/server"
	"github.com/yockii/ruomu-core/util"

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
	if instance.PageName == "" || instance.PageCode == "" {
		return &server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " name / code",
		}, nil
	}

	if c, err := database.DB.Count(&model.Page{PageCode: instance.PageCode}); err != nil {
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
	if instance.Id == 0 {
		return &server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " id",
		}, nil
	}

	if _, err := database.DB.Update(&model.Page{
		PageName:   instance.PageName,
		PageCode:   instance.PageCode,
		ThemeCode:  instance.ThemeCode,
		PageRoute:  instance.PageRoute,
		PageConfig: instance.PageConfig,
	}, &model.Page{Id: instance.Id}); err != nil {
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

func (_ *pageController) List(value []byte) (any, error) {
	instance := new(model.Page)
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

	condition := &model.Page{
		Id: instance.Id,
	}
	if instance.PageName != "" {
		session.Where("page_name like ?", "%"+instance.PageName+"%")
		instance.PageName = ""
	}
	if instance.PageCode != "" {
		session.Where("page_code like ?", "%"+instance.PageCode+"%")
		instance.PageCode = ""
	}
	if instance.ThemeCode != "" {
		session.Where("theme_code like ?", "%"+instance.ThemeCode+"%")
		instance.ThemeCode = ""
	}
	if instance.PageRoute != "" {
		session.Where("page_route like ?", "%"+instance.PageRoute+"%")
		instance.PageRoute = ""
	}

	var list []*model.Page
	total, err := session.Omit("page_config").FindAndCount(&list, condition)
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
	condition := new(model.Page)
	if err := json.Unmarshal(value, condition); err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		}, nil
	}
	instance := new(model.Page)
	if has, err := database.DB.Get(instance); err != nil {
		logger.Errorln(err)
		return &server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		}, nil
	} else if !has {
		return &server.CommonResponse{}, nil
	}
	return &server.CommonResponse{
		Data: instance,
	}, nil
}
