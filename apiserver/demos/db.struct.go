package demos

import (
	"time"

	"github.com/zhiyunliu/golibs/xtypes"
)

type Test struct {
	Binary       Binary          `json:"binary"`
	B            bool            `json:"bool"`
	Date         time.Time       `json:"date"`
	DateTime     *time.Time      `json:"datetime"`
	Decval       xtypes.Decimal  `json:"decval"`
	Floatval     float32         `json:"floatval"`
	Id           int             `json:"id"`
	Money        *xtypes.Decimal `json:"money"`
	Name         string          `json:"name"`
	Numeric      xtypes.Decimal  `json:"numeric"`
	Nvarchar_max string          `json:"nvarchar_max"`
	Status       int8            `json:"status"`
	Text         string          `json:"text"`
	Xmap         xtypes.XMap     `json:"xmap"`
	Xmaps        xtypes.XMaps    `json:"xmaps"`
	Xml          *string         `json:"xml"`
}
