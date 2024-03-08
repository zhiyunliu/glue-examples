package demos

import (
	"fmt"
	"strconv"
	"strings"

	"192.168.1.27/microservice/common/util"
	"github.com/zhiyunliu/golibs/xtypes"
	"github.com/zhiyunliu/golibs/xtypes/datetime"
)

type Test struct {
	Binary NBinary `json:"binary"`
	B      bool    `json:"bool"`
	//Date         time.Time          `json:"date"`
	//DateTime     *time.Time         `json:"datetime"`
	CDate        datetime.DateTime  `json:"date"`
	CDateTime    *datetime.DateTime `json:"datetime"`
	Decval       xtypes.Decimal     `json:"decval"`
	Floatval     float32            `json:"floatval"`
	Id           int                `json:"id"`
	Money        *xtypes.Decimal    `json:"money"`
	Name         string             `json:"name"`
	Numeric      xtypes.Decimal     `json:"numeric"`
	Nvarchar_max string             `json:"nvarchar_max"`
	Status       int8               `json:"status"`
	Text         string             `json:"text"`
	Xmap         xtypes.XMap        `json:"xmap"`
	Xmaps        xtypes.XMaps       `json:"xmaps"`
	Xml          *string            `json:"xml"`
	StrDate      datetime.DateTime  `json:"str_date"`
	StrDateTime  *datetime.DateTime `json:"str_datetime"`
	Utime        util.DateTime      `json:"utime"`
	PUtime       *util.DateTime     `json:"putime"`
}

type NBinary struct {
	Data []byte `json:"data"`
}

func (nb NBinary) MarshalJSON() ([]byte, error) {
	builder := strings.Builder{}
	for i := range nb.Data {
		builder.WriteString(strconv.Itoa(int(nb.Data[i])) + ",")
	}
	return []byte(fmt.Sprintf(`"%s"`, builder.String())), nil
}

func (nb *NBinary) Scan(data any) error {
	nb.Data = data.([]byte)
	return nil
}
