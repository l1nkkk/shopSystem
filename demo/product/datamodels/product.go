package datamodels

// l1nkkk: sql 和 imooc 标签分别用在哪
type Product struct {
	ID           int64  `json:"id" sql:"ID" imooc:"ID"`
	ProductName  string `json:"ProductName" sql:"productName" imooc:"ProductName"`
	ProductNum   int64  `json:"ProductNum" sql:"productNum" imooc:"ProductNum"`
	ProductImage string `json:"ProductImage" sql:"productImage" imooc:"ProductImage"`
	ProductUrl   string `json:"ProductUrl" sql:"productUrl" imooc:"ProductUrl"`
}
