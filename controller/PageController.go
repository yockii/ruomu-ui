package controller

import (
	"encoding/json"

	"github.com/pkg/errors"
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
	if instance.PageName == "" || instance.PageCode == "" {
		return &server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " name / code",
		}, nil
	}

	var c int64
	if err := database.DB.Model(&model.Page{}).Where(&model.Page{PageCode: instance.PageCode}).Count(&c).Error; err != nil {
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
		PageName:   instance.PageName,
		PageCode:   instance.PageCode,
		ThemeCode:  instance.ThemeCode,
		PageRoute:  instance.PageRoute,
		PageConfig: instance.PageConfig,
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

	tx := database.DB.Limit(paginate.Limit).Offset(paginate.Offset)

	condition := &model.Page{
		ID: instance.ID,
	}
	if instance.PageName != "" {
		tx.Where("page_name like ?", "%"+instance.PageName+"%")
		instance.PageName = ""
	}
	if instance.PageCode != "" {
		tx.Where("page_code like ?", "%"+instance.PageCode+"%")
		instance.PageCode = ""
	}
	if instance.ThemeCode != "" {
		tx.Where("theme_code like ?", "%"+instance.ThemeCode+"%")
		instance.ThemeCode = ""
	}
	if instance.PageRoute != "" {
		tx.Where("page_route like ?", "%"+instance.PageRoute+"%")
		instance.PageRoute = ""
	}

	var list []*model.Page
	var total int64
	err := tx.Omit("page_config").Find(&list, condition).Offset(-1).Count(&total).Error
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
