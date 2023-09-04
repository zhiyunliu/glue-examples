package dbconn

import (
	"fmt"
	"strings"

	"github.com/zhiyunliu/glue/global"
	"github.com/zhiyunliu/glue/xdb"
	"github.com/zhiyunliu/golibs/xtypes/igcmap"
)

const (
	_dbAppName = "app name"
	_dbSlowSql = "dbslowsql"
)

func Refactor(opts ...Option) error {
	dbOpt := &options{
		slowOpts: &slowsqlOpts{},
		connOpts: make(map[string]*connOpts),
	}

	for i := range opts {
		opts[i](dbOpt)
	}

	if dbOpt.slowOpts.IsValid() {
		xdb.RegistryLogger(&dbLogger{
			name:      _dbSlowSql,
			alarmCode: dbOpt.slowOpts.AlarmCode,
			queueName: dbOpt.slowOpts.QueueName,
		})
	}
	xdb.ConnRefactor = func(connName string, cfg *xdb.Config) (newcfg *xdb.Config, err error) {
		copts, ok := dbOpt.connOpts[connName]
		if ok {
			for i := range copts.DBOpts {
				copts.DBOpts[i](cfg)
			}
		}
		return dbConfigRefactor(connName, cfg)
	}
	return nil
}

func dbConfigRefactor(connName string, cfg *xdb.Config) (newcfg *xdb.Config, err error) {
	//server=127.0.0.1;database=test;user id=admin;password=test
	conn := cfg.Conn
	parties := strings.Split(conn, ";")
	connMap := igcmap.New(map[string]interface{}{})
	keyList := []string{}
	for i := range parties {
		subParties := strings.Split(parties[i], "=")
		if len(subParties) != 2 {
			err = fmt.Errorf("数据库配置错误!%s", conn)
			return
		}
		keyList = append(keyList, subParties[0])
		connMap.Set(subParties[0], subParties[1])
	}
	if _, ok := connMap.Get(_dbAppName); !ok {
		keyList = append(keyList, _dbAppName)
		connMap.Set(_dbAppName, global.AppName)
	}
	connList := make([]string, 0)
	for _, k := range keyList {
		v, _ := connMap.Get(k)
		connList = append(connList, fmt.Sprintf("%s=%s", k, v))
	}

	conn = strings.Join(connList, ";")
	cfg.Conn = conn
	return cfg, nil
}
