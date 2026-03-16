package daos

import (
	"context"
	"errors"
	"github.com/efucloud/common"
	"github.com/efucloud/eauth/pkg/config"
	"github.com/go-sql-driver/mysql"
)

var mysqlErrorsZh = map[uint16]string{
	1045: "数据库用户权限不足",
	1046: "系统未配置数据库",
	1049: "指定的数据库不存在",
	1054: "数据库查询中引用了不存在的列",
	1062: "试图插入或更新数据时数据重复",
	1064: "数据库语法错误",
	1146: "数据库表不存在",
	1216: "数据库死锁",
	1217: "数据库外键冲突",
}
var mysqlErrorsEn = map[uint16]string{
	1045: "Access denied for user",
	1046: "No database selected",
	1049: "Unknown database",
	1054: "Unknown column",
	1062: "Duplicate entry",
	1064: "Syntax error",
	1146: "Table does not exist",
	1216: "Deadlock found when trying to get lock",
	1217: "Cannot delete or update a parent row in use by a foreign key constraint",
}

func ParserDatabaseError(ctx context.Context, err common.ErrorData) common.ErrorData {
	var e *mysql.MySQLError
	if errors.As(err.Err, &e) {
		switch ctx.Value(config.RequestLanguage) {
		case common.I18nEN:
			err.MsgCode, _ = mysqlErrorsEn[e.Number]
		default:
			err.MsgCode, _ = mysqlErrorsZh[e.Number]
		}
	}
	return err
}
