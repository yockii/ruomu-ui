package domain

import "github.com/yockii/ruomu-ui/model"

type MaterialLibCondition struct {
	model.MaterialLib
	ProjectID uint64 `json:"projectId,omitempty,string"`
}
