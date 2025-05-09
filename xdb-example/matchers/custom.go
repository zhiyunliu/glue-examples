package matchers

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/zhiyunliu/glue/contrib/xdb/expression"
	"github.com/zhiyunliu/glue/xdb"
)

type CustomMatcher struct {
	regexp    *regexp.Regexp
	symbolMap xdb.SymbolMap
}

func NewCustomMatcher() xdb.ExpressionMatcher {
	m := &CustomMatcher{}
	m.regexp = regexp.MustCompile(m.Pattern())
	m.symbolMap = expression.DefaultSymbols
	m.symbolMap.Regist(&cusSymbol{})
	return m
}

// 表达式名称
func (m *CustomMatcher) Name() string {
	return "custom"
}

// 表达式正则匹配
func (m *CustomMatcher) Pattern() string {
	const pattern = `#(({((\w+\.)?\w+)\s*=\s*(\w+)}))`
	return pattern
}

// 符号列表
func (m *CustomMatcher) GetOperatorMap() xdb.OperatorMap {
	return xdb.NewOperatorMap(
		xdb.NewOperator("=", func(valuer xdb.ExpressionValuer, param xdb.DBParam, phName string, value any) string {
			return fmt.Sprintf("%s=%s", valuer.GetFullfield(), phName)
		}),
	)
}

// 匹配执行
func (m *CustomMatcher) MatchString(expression string) (xdb.ExpressionValuer, bool) {

	parties := m.regexp.FindStringSubmatch(expression)
	if len(parties) <= 0 {
		return nil, false
	}

	expression = strings.Trim(expression, "#{}")

	parties = strings.Split(expression, "=")

	item := &xdb.ExpressionItem{
		FullField: parties[0],
		PropName:  parties[1],
		Oper:      "=",
	}

	item.Symbol, _ = m.symbolMap.Load("#")
	return item, true
}

type cusSymbol struct {
}

func (s *cusSymbol) Name() string {
	return "#"
}
func (s *cusSymbol) Concat() string {
	return "and   "
}
func (s *cusSymbol) DynamicType() xdb.DynamicType {
	return xdb.DynamicNone
}
func (s *cusSymbol) IsDynamic() bool {
	return false
}
