package repositories

import (
	"database/sql"
	"errors"
	"github.com/l1nkkk/shopSystem/demo/product/common"
	"github.com/l1nkkk/shopSystem/demo/product/datamodels"
	"strconv"
)

type IUserRepository interface {
	Conn() error

	// Select 根据userName，返回查询到的User
	Select(userName string) (user *datamodels.User, err error)

	// Insert 插入User数据
	Insert(user *datamodels.User) (userId int64, err error)
}

func NewUserRepository(table string, db *sql.DB) IUserRepository {
	return &UserManagerRepository{table, db}
}

type UserManagerRepository struct {
	table     string
	mysqlConn *sql.DB
}

func (u *UserManagerRepository) Conn() (err error) {
	if u.mysqlConn == nil {
		mysql, errMysql := common.NewMysqlConn()
		if errMysql != nil {
			return errMysql
		}
		u.mysqlConn = mysql
	}
	if u.table == "" {
		u.table = "user"
	}
	return
}

func (u *UserManagerRepository) Select(userName string) (user *datamodels.User, err error) {
	// 1. 判断特殊情况
	if userName == "" {
		return &datamodels.User{}, errors.New("条件不能为空！")
	}
	if err = u.Conn(); err != nil {
		return &datamodels.User{}, err
	}

	// 2. 执行query，?为占位符，避免sql注入攻击
	sql := "Select * from " + u.table + " where userName=?"
	rows, errRows := u.mysqlConn.Query(sql, userName)
	defer rows.Close()
	if errRows != nil {
		return &datamodels.User{}, errRows
	}

	// 3. 获取行数据 rows ===> map[string]string
	result := common.GetResultRow(rows)
	if len(result) == 0 {
		return &datamodels.User{}, errors.New("用户不存在！")
	}

	// 4. map[string]string ==> user
	user = &datamodels.User{}
	common.DataToStructByTagSql(result, user)
	return
}

func (u *UserManagerRepository) Insert(user *datamodels.User) (userId int64, err error) {
	// 1. 判断连接是否已经建立
	if err = u.Conn(); err != nil {
		return
	}

	// 2. 准备sql
	sql := "INSERT " + u.table + " SET nickName=?,userName=?,passWord=?"
	stmt, errStmt := u.mysqlConn.Prepare(sql)

	// 3. 执行sql
	if errStmt != nil {
		return userId, errStmt
	}
	result, errResult := stmt.Exec(user.NickName, user.UserName, user.HashPassword)
	if errResult != nil {
		return userId, errResult
	}
	return result.LastInsertId()
}

// SelectByID 通过ID获得用户信息
func (u *UserManagerRepository) SelectByID(userId int64) (user *datamodels.User, err error) {
	// 1. 判断连接是否已经建立
	if err = u.Conn(); err != nil {
		return &datamodels.User{}, err
	}

	// 2. 执行query
	sql := "select * from " + u.table + " where ID=" + strconv.FormatInt(userId, 10)
	row, errRow := u.mysqlConn.Query(sql)
	if errRow != nil {
		return &datamodels.User{}, errRow
	}

	// 3. 获取行数据 rows ===> map[string]string
	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.User{}, errors.New("用户不存在！")
	}

	// 4. map[string]string ==> user
	user = &datamodels.User{}
	common.DataToStructByTagSql(result, user)
	return
}
