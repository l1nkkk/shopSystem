package datamodels

type Order struct {
	ID          int64 `sql:"ID"`          // 订单ID
	UserId      int64 `sql:"userID"`      // 用户ID
	ProductId   int64 `sql:"productID"`   // 商品ID
	OrderStatus int64 `sql:"orderStatus"` // 订单状态
}

const (
	OrderWait    = iota
	OrderSuccess //1
	OrderFailed  //2
)
