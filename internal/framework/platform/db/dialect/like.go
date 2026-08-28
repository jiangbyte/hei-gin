package dialect

import (
	"fmt"
	"regexp"
	"strings"

	"gorm.io/gorm"
)

var identRe = regexp.MustCompile(`(?i)^[a-z_][a-z0-9_]*(\.[a-z_][a-z0-9_]*)?$`)

// MustIdent 校验 SQL 标识符（列名/表名），拒绝注入字符。仅允许字母数字下划线与一层限定。
func MustIdent(name string) string {
	n := strings.TrimSpace(name)
	if !identRe.MatchString(n) {
		panic(fmt.Sprintf("dialect: unsafe sql identifier: %q", name))
	}
	return n
}

// EscapeLike 转义 LIKE 通配符，使 % / _ / \ 按字面匹配。
func EscapeLike(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `%`, `\%`)
	s = strings.ReplaceAll(s, `_`, `\_`)
	return s
}

// Contains 构造已转义的「包含」模糊模式（两侧 %）。
func Contains(s string) string {
	return "%" + EscapeLike(s) + "%"
}

// ILike 返回不区分大小写的模糊匹配子句（含一个 ? 占位符 + ESCAPE）。
// column 必须为安全标识符；参数请传 Contains(keyword)。
func ILike(db *gorm.DB, column string) string {
	col := MustIdent(column)
	if IsMySQL(db) {
		return fmt.Sprintf("LOWER(%s) LIKE LOWER(?) ESCAPE '\\\\'", col)
	}
	return fmt.Sprintf("%s ILIKE ? ESCAPE '\\\\'", col)
}
