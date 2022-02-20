package repositories

import (
	"database/sql"
	"github.com/l1nkkk/shopSystem/demo/product/common"
	"github.com/l1nkkk/shopSystem/demo/product/datamodels"
	"strconv"
)

//第一步，定义对应的接口
//第二步，实现定义的接口

// l1nkkk: 感觉这里定义个Close 会好点

type IProduct interface {
	// Conn 连接数据库
	Conn() error
	Insert(*datamodels.Product) (int64, error)
	Delete(int64) bool
	Update(*datamodels.Product) error
	SelectByKey(int64) (*datamodels.Product, error)
	SelectAll() ([]*datamodels.Product, error)
}

type ProductManager struct {
	table     string
	mysqlConn *sql.DB
}

func NewProductManager(table string, db *sql.DB) IProduct {
	return &ProductManager{table: table, mysqlConn: db}
}

// Conn 数据库连接，如果已经连接了则忽略（即 p.mysqlConn != nil）
func (p *ProductManager) Conn() (err error) {
	if p.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		p.mysqlConn = mysql
	}
	if p.table == "" {
		p.table = "product"
	}
	return
}

// Insert 插入商品数据
func (p *ProductManager) Insert(product *datamodels.Product) (productId int64, err error) {
	// 1.判断连接是否存在
	if err = p.Conn(); err != nil {
		return
	}
	// 2.准备sql
	sql := "INSERT product SET productName=?,productNum=?,productImage=?,productUrl=?"
	stmt, errSql := p.mysqlConn.Prepare(sql)
	if errSql != nil {
		return 0, errSql
	}

	// l1nkkk: 这里使用到了 sql 标签
	// 3.传入参数， 输出 Result, error
	result, errStmt := stmt.Exec(product.ProductName, product.ProductNum, product.ProductImage, product.ProductUrl)
	if errStmt != nil {
		return 0, errStmt
	}
	return result.LastInsertId()
}

// Delete 删除商品数据
func (p *ProductManager) Delete(productID int64) bool {
	// 1.判断连接是否存在
	if err := p.Conn(); err != nil {
		return false
	}
	// 2.准备sql
	sql := "delete from product where ID=?"
	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return false
	}

	// 3.传入参数执行sql
	_, err = stmt.Exec(strconv.FormatInt(productID, 10))
	if err != nil {
		return false
	}
	return true
}

// Update 商品的更新
func (p *ProductManager) Update(product *datamodels.Product) error {
	// 1.判断连接是否存在
	if err := p.Conn(); err != nil {
		return err
	}

	// 2.准备sql
	sql := "Update product set productName=?,productNum=?,productImage=?,productUrl=? where ID=" + strconv.FormatInt(product.ID, 10)
	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return err
	}

	// 3.执行sql
	_, err = stmt.Exec(product.ProductName, product.ProductNum, product.ProductImage, product.ProductUrl)
	if err != nil {
		return err
	}
	return nil
}

// SelectByKey 根据商品ID查询商品
func (p *ProductManager) SelectByKey(productID int64) (productResult *datamodels.Product, err error) {
	// 1.判断连接是否存在
	if err = p.Conn(); err != nil {
		return &datamodels.Product{}, err
	}

	// 2.执行sql，返回 Rows
	sql := "Select * from " + p.table + " where ID =" + strconv.FormatInt(productID, 10)
	row, errRow := p.mysqlConn.Query(sql)
	// l1nkkk: 注意这里的Close
	if errRow != nil {
		return &datamodels.Product{}, errRow
	}
	defer row.Close()

	// 3.获取结果
	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.Product{}, nil
	}

	// 4.转化成结构体
	productResult = &datamodels.Product{}
	common.DataToStructByTagSql(result, productResult)
	return

}

// SelectAll 获取所有商品
func (p *ProductManager) SelectAll() (productArray []*datamodels.Product, errProduct error) {
	//1.判断连接是否存在
	if err := p.Conn(); err != nil {
		return nil, err
	}

	// 2.执行sql，返回 Rows
	sql := "Select * from " + p.table
	rows, err := p.mysqlConn.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 3.获取结果
	result := common.GetResultRows(rows)
	if len(result) == 0 {
		return nil, nil
	}

	// 4.转化成结构体
	for _, v := range result {
		product := &datamodels.Product{}
		common.DataToStructByTagSql(v, product)
		productArray = append(productArray, product)
	}
	return
}
