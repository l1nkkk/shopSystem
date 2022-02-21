package repositories

import (
	"database/sql"
	"strconv"

	"github.com/l1nkkk/shopSystem/demo/productOptim/common"
	"github.com/l1nkkk/shopSystem/demo/productOptim/datamodels"
)

type IOrderRepository interface {
	Conn() error
	Insert(*datamodels.Order) (int64, error)
	Delete(int64) bool
	Update(*datamodels.Order) error
	SelectByKey(int64) (*datamodels.Order, error)
	SelectAll() ([]*datamodels.Order, error)
	SelectAllWithInfo() (map[int]map[string]string, error)
}

func NewOrderMangerRepository(table string, sql *sql.DB) IOrderRepository {
	return &OrderMangerRepository{table: table, mysqlConn: sql}
}

type OrderMangerRepository struct {
	table     string
	mysqlConn *sql.DB
}

// Conn 数据库连接，如果mysqlConn!= nil, 则忽略
func (o *OrderMangerRepository) Conn() error {
	// 1. 判断是否已经连接
	if o.mysqlConn == nil {
		// 2. 进行数据库连接
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		o.mysqlConn = mysql
	}
	if o.table == "" {
		o.table = "orders"
	}
	return nil
}

// Insert 插入订单数据
func (o *OrderMangerRepository) Insert(order *datamodels.Order) (productID int64, err error) {
	// 1. 判断是否已经连接
	if err = o.Conn(); err != nil {
		return
	}

	// 2. 准备sql
	//sql := "INSERT product SET productName=?,productNum=?,productImage=?,productUrl=?"
	sql := "INSERT " + o.table + " SET userID=?,productID=?,orderStatus=?"
	stmt, errStmt := o.mysqlConn.Prepare(sql)
	if errStmt != nil {
		return productID, errStmt
	}

	// 3. 执行sql
	result, errResult := stmt.Exec(order.UserId, order.ProductId, order.OrderStatus)
	if errResult != nil {
		return productID, errResult
	}
	return result.LastInsertId()
}


// Delete 删除订单数据
func (o *OrderMangerRepository) Delete(orderID int64) (isOk bool) {
	// 1. 判断是否已经连接
	if err := o.Conn(); err != nil {
		return
	}

	// 2. 准备sql
	sql := "delete from " + o.table + " where ID =?"
	stmt, errStmt := o.mysqlConn.Prepare(sql)
	if errStmt != nil {
		return
	}

	// 3. 执行sql
	_, err := stmt.Exec(orderID)
	if err != nil {
		return
	}
	return true
}

// Update 更新订单数据
func (o *OrderMangerRepository) Update(order *datamodels.Order) (err error) {
	// 1. 判断是否已经连接
	if errConn := o.Conn(); errConn != nil {
		return errConn
	}

	// 2. 准备sql
	sql := "Update " + o.table + " set userID=?,productID=?,orderStatus=? Where ID=" + strconv.FormatInt(order.ID, 10)
	stmt, errStmt := o.mysqlConn.Prepare(sql)
	if errStmt != nil {
		return errStmt
	}

	// 3. 执行sql
	_, err = stmt.Exec(order.UserId, order.ProductId, order.OrderStatus)
	return
}

// SelectByKey 根据orderID检索单条订单数据
func (o *OrderMangerRepository) SelectByKey(orderID int64) (order *datamodels.Order, err error) {
	// 1. 判断是否已经连接
	if errConn := o.Conn(); errConn != nil {
		return &datamodels.Order{}, errConn
	}

	// 2. 构造并执行query
	sql := "Select * From " + o.table + " where ID=" + strconv.FormatInt(orderID, 10)
	row, errRow := o.mysqlConn.Query(sql)
	if errRow != nil {
		return &datamodels.Order{}, errRow
	}
	defer row.Close()

	// 3. 获取行数据  map[string]string
	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.Order{}, err
	}

	// 4.  map[string]string ===> order struct
	order = &datamodels.Order{}
	common.DataToStructByTagSql(result, order)
	return
}

// SelectAll 查询所有的order
func (o *OrderMangerRepository) SelectAll() (orderArray []*datamodels.Order, err error) {
	// 1. 判断是否已经连接
	if errConn := o.Conn(); errConn != nil {
		return nil, errConn
	}
	// 2. 构造并执行query
	sql := "Select * from " + o.table
	rows, errRows := o.mysqlConn.Query(sql)
	if errRows != nil {
		return nil, errRows
	}
	defer rows.Close()


	// 3. 获取行数据   map[int]map[string]string
	result := common.GetResultRows(rows)
	if len(result) == 0 {
		return nil, err
	}

	// 4.  map[int]map[string]string ==> []order struct
	for _, v := range result {
		order := &datamodels.Order{}
		common.DataToStructByTagSql(v, order)
		orderArray = append(orderArray, order)
	}
	return
}

// SelectAllWithInfo 查询所有的order（输出以下信息 {orderID, 商品名, order 状态}）
func (o *OrderMangerRepository) SelectAllWithInfo() (OrderMap map[int]map[string]string, err error) {
	// 1. 判断是否已经连接
	if errConn := o.Conn(); errConn != nil {
		return nil, errConn
	}

	// 2. 构造并执行query
	// orderID, 商品名, order 状态
	sql := "Select o.ID,p.productName,o.orderStatus From imooc.order as o left join product as p on o.productID=p.ID"
	rows, errRows := o.mysqlConn.Query(sql)
	if errRows != nil {
		return nil, errRows
	}
	return common.GetResultRows(rows), err
}
