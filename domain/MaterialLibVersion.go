package domain

import "github.com/yockii/ruomu-ui/model"

type MaterialLibVersionCondition struct {
	model.MaterialLibVersion
}

type MaterialLibVersionImportExport struct {
	Lib        *model.MaterialLib              `json:"lib"`
	Version    *model.MaterialLibVersion       `json:"version"`
	Groups     []*model.MaterialComponentGroup `json:"groups"`
	Components []*model.MaterialComponent      `json:"components"`
}
