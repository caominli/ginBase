package models //定义在models包下

//商品数据库模型
type Items struct {
	ID   uint   `gorm:"primaryKey"`        //id
	Name string `gorm:"not null;size:500"` //商品名称
}
