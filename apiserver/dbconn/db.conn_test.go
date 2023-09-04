package dbconn

import (
	"testing"

	"github.com/zhiyunliu/glue/global"
	"github.com/zhiyunliu/glue/xdb"
)

func Test_dbConnRefactor(t *testing.T) {

	global.AppName = "test"
	tests := []struct {
		name        string
		conn        string
		wantNewConn string
		wantErr     bool
	}{
		{name: "1.", conn: "aaa=1;bbb=2", wantNewConn: "aaa=1;bbb=2;app name=test", wantErr: false},
		{name: "2.", conn: "aaa=1;bbb=2;app name=demo", wantNewConn: "aaa=1;bbb=2;app name=demo", wantErr: false},
		{name: "3.", conn: "aaa", wantNewConn: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCfg, err := dbConfigRefactor("", &xdb.Config{Conn: tt.conn})
			if (err != nil) != tt.wantErr {
				t.Errorf("dbConnRefactor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			var gotNewConn = ""
			if gotCfg != nil {
				gotNewConn = gotCfg.Conn
			}
			if gotNewConn != tt.wantNewConn {
				t.Errorf("dbConnRefactor() = %v, want %v", gotNewConn, tt.wantNewConn)
			}
		})
	}
}
