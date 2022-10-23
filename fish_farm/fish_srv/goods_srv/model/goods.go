package model

// Category (分类表)
// 生产环境下，建议设置 not null
// proto 内无 int 类型，但是有 int32、int64
type Category struct {
	BaseModel
	Name             string      `gorm:"type:varchar(20);not null" json:"name"`
	ParentCategoryID int32       `json:"parent"`
	ParentCategory   *Category   `json:"-"`
	SubCategory      []*Category `gorm:"foreignKey:ParentCategoryID;references:ID" json:"sub_category"`
	Level            int32       `gorm:"type:int;not null;default:1; comment '产品分类级别'" json:"level"`
	IsTab            bool        `gorm:"default:false;not null" json:"is_tab""`
}

// Brands (品牌表)
type Brands struct {
	BaseModel
	Name string `gorm:"type:varchar(20);not null; comment '品牌名称'"`
	Logo string `gorm:"type:varchar(20);not null;default ''; comment '品牌logo'"`
}

// GoodsCategoryBrand (连接表)
//创建 品牌ID与分类ID的 联合唯一索引
type GoodsCategoryBrand struct {
	BaseModel
	CategoryID int32 `gorm:"type:int;index:idx_category_brand,unique comment '联合唯一索引'"`
	Category   Category
	BrandsID   int32 `gorm:"type:int;index:idx_category_brand,unique comment '联合唯一索引'"`
	Brands     Brands
}

// TableName
//gorm默认创建表名时会默认将，大写转小写使用下划线来进行连接
//自定义表名（重载表名）
func (GoodsCategoryBrand) TableName() string {
	return "goodscategorybrand"
}

// Banner (轮播图)
type Banner struct {
	BaseModel
	Image string `gorm:"type:varchar(200);not null"`
	Url   string `gorm:"type:varchar(200);not null"`
	Index int32  `gorm:"type:int;default:1;not null"`
}

// Goods (商品表)
//数据库内无法存储数组类型，只有 json 类型
//1.使用 gorm 自定义一个类型
//2.另建一张表用来存储
type Goods struct {
	BaseModel
	CategoryID int32 `gorm:"type:int;not null"`
	Category   Category
	BrandsID   int32 `gorm:"type:int;not null"`
	Brands     Brands
	OnSale     bool `gorm:"default:false;not null comment '是否上架'"`
	ShipFree   bool `gorm:"default:false;not null comment '是否免运费'"`
	IsNew      bool `gorm:"default:false;not null comment '是否为热销商品'"`
	IsHot      bool `gorm:"default:false;not null comment '是否放置于主页'"`

	Name            string   `gorm:"type:varchar(50);not null comment '商品名'"`
	GoodsSn         string   `gorm:"type:varchar(50);not null comment '商品内部编号'"`
	ClickNum        int32    `gorm:"type:int;default:0;not null comment '点击量'"`
	SoldNum         int32    `gorm:"type:int;default:0;not null comment '销量'"`
	FavNum          int32    `gorm:"type:int;default:0;not null comment '收藏量'"`
	MarketPrice     float32  `gorm:"not null comment '市场价格'"`
	ShopPrice       float32  `gorm:"not null comment '销售价格'"`
	GoodsBrief      string   `gorm:"type:varchar(100);not null comment '商品名'"`
	Images          GormList `gorm:"type:varchar(1000);not null comment '商品页展示预览图'"`
	DescImages      GormList `gorm:"type:varchar(1000);not null comment '商品详情页展示图'"`
	GoodsFrontImage string   `gorm:"type:varchar(200);not null comment '预览图'"`
}
