package dbconn

import "github.com/zhiyunliu/glue/xdb"

type options struct {
	connOpts map[string]*connOpts
	slowOpts *slowsqlOpts
}
type slowsqlOpts struct {
	AlarmCode string
	QueueName string
}

func (o *slowsqlOpts) IsValid() bool {
	return o.AlarmCode != "" && o.QueueName != ""
}

type connOpts struct {
	ConnName string
	DBOpts   []xdb.Option
}

type Option func(opts *options)

func WithDefault(dbOpts ...xdb.Option) Option {
	return func(opts *options) {
		for i := range dbOpts {
			dbOpts[i](xdb.Default)
		}
	}
}

// 配置慢sql
func WithSlowsql(alarmCode, queueName string) Option {
	return func(opts *options) {
		opts.slowOpts.AlarmCode = alarmCode
		opts.slowOpts.QueueName = queueName
	}
}

// 配置慢sql
func WithConfig(connName string, dbOpts ...xdb.Option) Option {
	return func(opts *options) {
		connOpt, ok := opts.connOpts[connName]
		if !ok {
			connOpt = &connOpts{}
			opts.connOpts[connName] = connOpt
		}
		connOpt.ConnName = connName
		connOpt.DBOpts = dbOpts
	}
}
