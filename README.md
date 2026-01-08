## 📚 API 接口文档

## 基础路由

| # | 方法 | 路径 | Handler | 说明 |
|---|------|------|---------|------|
| 1 | GET | `/api/v1/health` | HealthCheck | 健康检查 |
| 2 | GET | `/api/v1/version` | Version | 版本信息 |

## 认证模块 (Auth)

| # | 方法 | 路径 | Handler | 说明 |
|---|------|------|---------|------|
| 3 | POST | `/api/v1/auth/register` | Register | 用户注册 |
| 4 | POST | `/api/v1/auth/login` | Login | 用户登录 |
| 5 | POST | `/api/v1/auth/logout` | Logout | 用户登出 |

## 用户模块 (User)

| # | 方法 | 路径 | Handler | 说明 |
|---|------|------|---------|------|
| 6 | GET | `/api/v1/users/me` | GetCurrentUser | 获取当前用户信息 |
| 7 | PUT | `/api/v1/users/me` | UpdateCurrentUser | 更新当前用户信息 |
| 8 | GET | `/api/v1/users/me/listings` | GetMyListings | 获取我的发布 |

## 房产模块 (Property)

| # | 方法 | 路径 | Handler | 说明 |
|---|------|------|---------|------|
| 9 | GET | `/api/v1/properties` | ListProperties | 房产列表（支持筛选） |
| 10 | GET | `/api/v1/properties/:id` | GetProperty | 房产详情 |
| 11 | POST | `/api/v1/properties` | CreateProperty | 创建房产（需认证） |
| 12 | PUT | `/api/v1/properties/:id` | UpdateProperty | 更新房产（需认证） |
| 13 | DELETE | `/api/v1/properties/:id` | DeleteProperty | 删除房产（需认证） |
| 14 | GET | `/api/v1/properties/:id/similar` | GetSimilarProperties | 相似房源 |
| 15 | GET | `/api/v1/properties/featured` | GetFeaturedProperties | 精选房源 |
| 16 | GET | `/api/v1/properties/hot` | GetHotProperties | 热门房源 |

### 买房 (Buy)

| # | 方法 | 路径 | Handler | 说明 |
|---|------|------|---------|------|
| 17 | GET | `/api/v1/properties/buy` | ListBuyProperties | 买房房源列表 |
| 18 | GET | `/api/v1/properties/buy/new` | ListNewProperties | 新房列表 |
| 19 | GET | `/api/v1/properties/buy/secondhand` | ListSecondhandProperties | 二手房列表 |

### 租房 (Rent)

| # | 方法 | 路径 | Handler | 说明 |
|---|------|------|---------|------|
| 20 | GET | `/api/v1/properties/rent` | ListRentProperties | 租房房源列表 |
| 21 | GET | `/api/v1/properties/rent/short-term` | ListShortTermRent | 短租房源 |
| 22 | GET | `/api/v1/properties/rent/long-term` | ListLongTermRent | 长租房源 |

### 新盘 (New Properties)

| # | 方法 | 路径 | Handler | 说明 |
|---|------|------|---------|------|
| 23 | GET | `/api/v1/new-properties` | ListNewDevelopments | 新楼盘列表 |
| 24 | GET | `/api/v1/new-properties/:id` | GetNewDevelopment | 新楼盘详情 |
| 25 | GET | `/api/v1/new-properties/:id/units` | GetDevelopmentUnits | 楼盘单位列表 |

### 服务式住宅 (Serviced Apartments)

| # | 方法 | 路径 | Handler | 说明 |
|---|------|------|---------|------|
| 26 | GET | `/api/v1/serviced-apartments` | ListServicedApartments | 服务式公寓列表 |
| 27 | GET | `/api/v1/serviced-apartments/:id` | GetServicedApartment | 服务式公寓详情 |
| 28 | GET | `/api/v1/serviced-apartments/:id/units` | GetServicedApartmentUnits | 服务式公寓房型列表 |
| 29 | GET | `/api/v1/serviced-apartments/:id/images` | GetServicedApartmentImages | 服务式公寓图片 |
| 30 | POST | `/api/v1/serviced-apartments` | CreateServicedApartment | 创建服务式公寓（需认证） |
| 31 | PUT | `/api/v1/serviced-apartments/:id` | UpdateServicedApartment | 更新服务式公寓（需认证） |
| 32 | DELETE | `/api/v1/serviced-apartments/:id` | DeleteServicedApartment | 删除服务式公寓（需认证） |

## 屋苑模块 (Estates)

| # | 方法 | 路径 | Handler | 说明 |
|---|------|------|---------|------|
| 33 | GET | `/api/v1/estates` | ListEstates | 屋苑列表 |
| 34 | GET | `/api/v1/estates/:id` | GetEstate | 屋苑详情 |
| 35 | GET | `/api/v1/estates/:id/properties` | GetEstateProperties | 屋苑内房源列表 |
| 36 | GET | `/api/v1/estates/:id/images` | GetEstateImages | 屋苑图片 |
| 37 | GET | `/api/v1/estates/:id/facilities` | GetEstateFacilities | 屋苑设施 |
| 38 | GET | `/api/v1/estates/:id/transactions` | GetEstateTransactions | 屋苑成交记录 |
| 39 | GET | `/api/v1/estates/:id/statistics` | GetEstateStatistics | 屋苑统计数据 |
| 40 | GET | `/api/v1/estates/featured` | GetFeaturedEstates | 精选屋苑 |
| 41 | POST | `/api/v1/estates` | CreateEstate | 创建屋苑（需认证） |
| 42 | PUT | `/api/v1/estates/:id` | UpdateEstate | 更新屋苑（需认证） |
| 43 | DELETE | `/api/v1/estates/:id` | DeleteEstate | 删除屋苑（需认证） |


## 物业估价模块 (Valuation)

| # | 方法 | 路径 | Handler | 说明 |
|---|------|------|---------|------|
| 44 | GET | `/api/v1/valuation` | ListValuations | 获取屋苑估价列表 |
| 45 | GET | `/api/v1/valuation/:estateId` | GetEstateValuation | 获取指定屋苑估价参考 |
| 46 | GET | `/api/v1/valuation/search` | SearchValuations | 搜索屋苑估价 |
| 47 | GET | `/api/v1/valuation/districts/:districtId` | GetDistrictValuations | 获取地区屋苑估价列表 |

## 家具商城模块 (Furniture)

### 家具商品

| # | 方法 | 路径 | Handler | 说明 |
|---|------|------|---------|------|
| 48 | GET | `/api/v1/furniture` | ListFurniture | 家具列表 |
| 49 | GET | `/api/v1/furniture/categories` | GetFurnitureCategories | 家具分类 |
| 50 | GET | `/api/v1/furniture/:id` | GetFurniture | 家具详情 |
| 51 | POST | `/api/v1/furniture` | CreateFurniture | 发布家具（需认证） |
| 52 | PUT | `/api/v1/furniture/:id` | UpdateFurniture | 更新家具（需认证） |
| 53 | DELETE | `/api/v1/furniture/:id` | DeleteFurniture | 删除家具（需认证） |
| 54 | GET | `/api/v1/furniture/:id/images` | GetFurnitureImages | 家具图片 |
| 55 | PUT | `/api/v1/furniture/:id/status` | UpdateFurnitureStatus | 更新家具状态 |
| 56 | GET | `/api/v1/furniture/featured` | GetFeaturedFurniture | 精选家具 |

### 购物车

| # | 方法 | 路径 | Handler | 说明 |
|---|------|------|---------|------|
| 57 | GET | `/api/v1/cart` | GetCart | 获取购物车 |
| 58 | POST | `/api/v1/cart/items` | AddToCart | 添加到购物车 |
| 59 | PUT | `/api/v1/cart/items/:id` | UpdateCartItem | 更新购物车项 |
| 60 | DELETE | `/api/v1/cart/items/:id` | RemoveFromCart | 移除购物车项 |
| 61 | DELETE | `/api/v1/cart` | ClearCart | 清空购物车 |


## 按揭模块 (Mortgage)

| # | 方法 | 路径 | Handler | 说明 |
|---|------|------|---------|------|
| 62 | POST | `/api/v1/mortgage/calculate` | CalculateMortgage | 按揭计算 |
| 63 | GET | `/api/v1/mortgage/rates` | GetMortgageRates | 银行利率列表 |
| 64 | GET | `/api/v1/mortgage/rates/bank/:bank_id` | GetBankMortgageRate | 获取指定银行利率 |
| 65 | POST | `/api/v1/mortgage/rates/compare` | CompareMortgageRates | 比较银行利率 |
| 66 | POST | `/api/v1/mortgage/apply` | ApplyMortgage | 按揭申请 |
| 67 | GET | `/api/v1/mortgage/applications` | GetMortgageApplications | 获取按揭申请列表 |
| 68 | GET | `/api/v1/mortgage/applications/:id` | GetMortgageApplication | 获取按揭申请详情 |

## 新闻资讯模块 (News)

> 注：新闻内容通过爬虫自动获取，不提供手动创建/编辑功能

| # | 方法 | 路径 | Handler | 说明 |
|---|------|------|---------|------|
| 69 | GET | `/api/v1/news` | ListNews | 新闻列表 |
| 70 | GET | `/api/v1/news/categories` | GetNewsCategories | 新闻分类 |
| 71 | GET | `/api/v1/news/:id` | GetNews | 新闻详情 |
| 72 | GET | `/api/v1/news/hot` | GetHotNews | 热门新闻 |
| 73 | GET | `/api/v1/news/featured` | GetFeaturedNews | 精选新闻 |
| 74 | GET | `/api/v1/news/latest` | GetLatestNews | 最新新闻 |
| 75 | GET | `/api/v1/news/:id/related` | GetRelatedNews | 相关新闻 |

## 校网模块 (School Net)

| # | 方法 | 路径 | Handler | 说明 |
|---|------|------|---------|------|
| 76 | GET | `/api/v1/school-nets` | ListSchoolNets | 校网列表 |
| 77 | GET | `/api/v1/school-nets/:id` | GetSchoolNet | 校网详情 |
| 78 | GET | `/api/v1/school-nets/:id/schools` | GetSchoolsInNet | 校网内学校 |
| 79 | GET | `/api/v1/school-nets/:id/properties` | GetPropertiesInNet | 校网内房源 |
| 80 | GET | `/api/v1/school-nets/:id/estates` | GetEstatesInNet | 校网内屋苑 |
| 81 | GET | `/api/v1/school-nets/search` | SearchSchoolNets | 搜索校网 |
| 82 | GET | `/api/v1/schools` | ListSchools | 学校列表 |
| 83 | GET | `/api/v1/schools/:id/school-net` | GetSchoolNet | 获取学校所属校网 |
| 84 | GET | `/api/v1/schools/search` | SearchSchools | 搜索学校 |

## 地产代理模块 (Agents)

### 代理人

| # | 方法 | 路径 | Handler | 说明 |
|---|------|------|---------|------|
| 85 | GET | `/api/v1/agents` | ListAgents | 代理人列表 |
| 86 | GET | `/api/v1/agents/:id` | GetAgent | 代理人详情 |
| 87 | GET | `/api/v1/agents/:id/properties` | GetAgentProperties | 代理人房源列表 |
| 88 | POST | `/api/v1/agents/:id/contact` | ContactAgent | 联系代理人 |

### 代理公司

| # | 方法 | 路径 | Handler | 说明 |
|---|------|------|---------|------|
| 89 | GET | `/api/v1/agencies` | ListAgencies | 代理公司列表 |
| 90 | GET | `/api/v1/agencies/:id` | GetAgency | 代理公司详情 |
| 91 | GET | `/api/v1/agencies/:id/properties` | GetAgencyProperties | 代理公司房源列表 |
| 92 | POST | `/api/v1/agencies/:id/contact` | ContactAgency | 联系代理公司 |
| 93 | GET | `/api/v1/agencies/search` | SearchAgencies | 搜索代理公司 |

## 楼价指数模块 (Price Index)

| # | 方法 | 路径 | Handler | 说明 |
|---|------|------|---------|------|
| 94 | GET | `/api/v1/price-index` | GetPriceIndex | 楼价指数 |
| 95 | GET | `/api/v1/price-index/latest` | GetLatestPriceIndex | 最新楼价指数 |
| 96 | GET | `/api/v1/price-index/districts/:districtId` | GetDistrictPriceIndex | 地区楼价指数 |
| 97 | GET | `/api/v1/price-index/estates/:estateId` | GetEstatePriceIndex | 屋苑楼价指数 |
| 98 | GET | `/api/v1/price-index/trends` | GetPriceTrends | 价格走势 |
| 99 | GET | `/api/v1/price-index/compare` | ComparePriceIndex | 对比楼价指数 |
| 100 | GET | `/api/v1/price-index/export` | ExportPriceData | 数据导出 |
| 101 | GET | `/api/v1/price-index/history` | GetPriceIndexHistory | 历史楼价指数 |
| 102 | POST | `/api/v1/price-index` | CreatePriceIndex | 创建楼价指数（需认证） |
| 103 | PUT | `/api/v1/price-index/:id` | UpdatePriceIndex | 更新楼价指数（需认证） |

## 通用模块

### 地区

| # | 方法 | 路径 | Handler | 说明 |
|---|------|------|---------|------|
| 104 | GET | `/api/v1/districts` | ListDistricts | 地区列表 |
| 105 | GET | `/api/v1/districts/:id` | GetDistrict | 地区详情 |
| 106 | GET | `/api/v1/districts/:id/properties` | GetDistrictProperties | 地区内房源 |
| 107 | GET | `/api/v1/districts/:id/estates` | GetDistrictEstates | 地区内屋苑 |
| 108 | GET | `/api/v1/districts/:id/statistics` | GetDistrictStatistics | 地区统计数据 |

### 设施

| # | 方法 | 路径 | Handler | 说明 |
|---|------|------|---------|------|
| 109 | GET | `/api/v1/facilities` | ListFacilities | 设施列表 |
| 110 | GET | `/api/v1/facilities/:id` | GetFacility | 设施详情 |
| 111 | POST | `/api/v1/facilities` | CreateFacility | 创建设施（需认证） |
| 112 | PUT | `/api/v1/facilities/:id` | UpdateFacility | 更新设施（需认证） |
| 113 | DELETE | `/api/v1/facilities/:id` | DeleteFacility | 删除设施（需认证） |


### 文件上传

| # | 方法 | 路径 | Handler | 说明 |
|---|------|------|---------|------|
| 114 | POST | `/api/v1/upload` | UploadFile | 文件上传 |
| 115 | POST | `/api/v1/upload/multiple` | UploadMultipleFiles | 批量上传 |
| 116 | POST | `/api/v1/upload/image` | UploadImage | 图片上传 |
| 117 | DELETE | `/api/v1/upload/:id` | DeleteFile | 删除文件 |

### 搜索

| # | 方法 | 路径 | Handler | 说明 |
|---|------|------|---------|------|
| 118 | GET | `/api/v1/search` | GlobalSearch | 全局搜索 |
| 119 | GET | `/api/v1/search/properties` | SearchProperties | 搜索房产 |
| 120 | GET | `/api/v1/search/estates` | SearchEstates | 搜索屋苑 |
| 121 | GET | `/api/v1/search/agents` | SearchAgents | 搜索代理人 |
| 122 | GET | `/api/v1/search/suggestions` | GetSearchSuggestions | 搜索建议 |
| 123 | GET | `/api/v1/search/history` | GetSearchHistory | 搜索历史 |

### 统计分析

| # | 方法 | 路径 | Handler | 说明 |
|---|------|------|---------|------|
| 124 | GET | `/api/v1/statistics/overview` | GetOverviewStatistics | 总览统计 |
| 125 | GET | `/api/v1/statistics/properties` | GetPropertyStatistics | 房产统计 |
| 126 | GET | `/api/v1/statistics/transactions` | GetTransactionStatistics | 成交统计 |
| 127 | GET | `/api/v1/statistics/users` | GetUserStatistics | 用户统计 |

### 系统配置

| # | 方法 | 路径 | Handler | 说明 |
|---|------|------|---------|------|
| 128 | GET | `/api/v1/config` | GetConfig | 获取系统配置 |
| 129 | GET | `/api/v1/config/regions` | GetRegions | 获取区域配置 |
| 130 | GET | `/api/v1/config/property-types` | GetPropertyTypes | 获取房产类型配置 |