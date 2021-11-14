package data

import "go-advance/week04/pkg/mysql/model"

// 类型，是否为 null
type Category struct {
	model.Base
	Name     string `gorm:"type:varchar(20);not null;unique"`
	ParentID int
	Parent   *Category
	Level    int     `gorm:"type:int;default:1;not null"`
	IsTab    bool    `gorm:"default:false;not null"`
	Brands   []Brand `gorm:"many2many:category_brands;"`
	Goods    []Goods
}

type Brand struct {
	model.Base
	Name string `gorm:"type:varchar(20);not null"`
	Logo string `gorm:"type:varchar(200);default:'';not null"`

	Goods []Goods
}

type Banner struct {
	model.Base
	Image          string `gorm:"type:varchar(200);not null"`
	GoodsDetailUrl string `gorm:"type:varchar(200);not null"`
	Index          int32  `gorm:"type:int;default:1;not null"`
}

type Goods struct {
	model.Base
	Name    string `gorm:"type:varchar(50);not null"`
	GoodsSn string `gorm:"type:varchar(50);not null"`

	CategoryID int `gorm:"type:int;not null"`
	Category   *Category

	BrandID int `gorm:"type:int;not null"`
	Brand   *Brand

	IsOnSale   bool `gorm:"default:false;not null"`
	IsShipFree bool `gorm:"default:false;not null"`
	IsNew      bool `gorm:"default:false;not null"`
	IsHot      bool `gorm:"default:false;not null"`

	Clicked     int           `gorm:"type:int;default 0;not null"`
	Sold        int           `gorm:"type:int;default 0;not null"`
	FavCount    int           `gorm:"type:int;default 0;not null"`
	MarketPrice int           `gorm:"type:int;not null"`
	ShopPrice   int           `gorm:"type:int;not null"`
	Brief       string        `gorm:"type:varchar(100);not null"`
	Image       string        `gorm:"type:varchar(200);not null"`
	StyleImages model.Strings `gorm:"type:varchar(1000);not null"`
	DescImages  model.Strings `gorm:"type:varchar(1000);not null"`
}

type GoodsImage struct {
	GoodsID     int
	StyleImages string
}
