# AJO Living 数据库设计文档

## 文档说明

本文档描述 AJO Living 房产平台的数据库设计，包含所有数据模型的字段定义和关系说明。

生成日期：2025年12月18日

---

## 1. 用户模块

### 1.1 用户表 (users)

用户分为两种类型：普通用户 (individual) 和地产代理公司 (agency)

| 字段名 | 类型 | 必填 | 说明 | 索引 |
|--------|------|------|------|------|
| id | BIGINT UNSIGNED | 是 | 用户ID（主键，自增） | PRIMARY |
| user_type | VARCHAR(20) | 是 | 用户类型：individual=普通用户, agency=地产代理公司 | INDEX |
| email | VARCHAR(255) | 是 | 邮箱地址 | UNIQUE |
| password_hash | VARCHAR(255) | 是 | 密码哈希值 | - |
| name | VARCHAR(100) | 是 | 用户名称/公司名称 | - |
| phone | VARCHAR(20) | 否 | 联系电话 | - |
| status | VARCHAR(20) | 是 | 状态：active=活跃, inactive=停用, suspended=暂停 | INDEX |
| email_verified | BOOLEAN | 是 | 邮箱是否已验证 | - |
| email_verified_at | TIMESTAMP | 否 | 邮箱验证时间 | - |
| last_login_at | TIMESTAMP | 否 | 最后登录时间 | - |
| created_at | TIMESTAMP | 是 | 创建时间 | INDEX |
| updated_at | TIMESTAMP | 是 | 更新时间 | - |
| deleted_at | TIMESTAMP | 否 | 软删除时间 | INDEX |

**说明：**
- `user_type` 区分普通用户和地产代理公司
- 邮箱必须唯一，用于登录
- 使用软删除，保留历史数据

---

## 2. 房产模块

### 2.1 房产表 (properties)

房产信息表，包含买卖和租赁两种类型的房产

| 字段名 | 类型 | 必填 | 说明 | 索引 |
|--------|------|------|------|------|
| id | BIGINT UNSIGNED | 是 | 房产ID（主键，自增） | PRIMARY |
| property_no | VARCHAR(50) | 是 | 物业编号（系统生成） | UNIQUE |
| estate_no | VARCHAR(50) | 否 | 楼盘编号（外部编号） | INDEX |
| listing_type | VARCHAR(20) | 是 | 房源类型：sale=出售, rent=出租 | INDEX |
| title | VARCHAR(255) | 是 | 房产标题 | - |
| description | TEXT | 否 | 房产描述 | - |
| area | DECIMAL(10,2) | 是 | 面积（平方尺） | - |
| price | DECIMAL(15,2) | 是 | 价格（港币）- 售价或月租 | INDEX |
| address | VARCHAR(500) | 是 | 详细地址 | - |
| district_id | BIGINT UNSIGNED | 是 | 所属地区ID | INDEX |
| building_name | VARCHAR(200) | 否 | 大厦/楼宇名称 | INDEX |
| floor | VARCHAR(20) | 否 | 楼层（如：10/F, G/F） | - |
| orientation | VARCHAR(50) | 否 | 座向（如：东南、西北） | - |
| bedrooms | INT | 是 | 房间数 | INDEX |
| bathrooms | INT | 否 | 浴室数 | - |
| primary_school_net | VARCHAR(50) | 否 | 小学校网 | INDEX |
| secondary_school_net | VARCHAR(50) | 否 | 中学校网 | INDEX |
| property_type | VARCHAR(50) | 是 | 物业类型：apartment=公寓, villa=别墅, townhouse=联排别墅等 | INDEX |
| status | VARCHAR(20) | 是 | 状态：available=可用, pending=待定, sold=已售/已租, cancelled=已取消 | INDEX |
| publisher_id | BIGINT UNSIGNED | 是 | 发布者ID（关联users表） | INDEX |
| publisher_type | VARCHAR(20) | 是 | 发布者类型：individual=个人, agency=代理公司 | - |
| agent_id | BIGINT UNSIGNED | 否 | 负责地产代理ID（关联agents表） | INDEX |
| view_count | INT | 是 | 浏览次数 | - |
| favorite_count | INT | 是 | 收藏次数 | - |
| published_at | TIMESTAMP | 否 | 发布时间 | INDEX |
| expired_at | TIMESTAMP | 否 | 过期时间 | INDEX |
| created_at | TIMESTAMP | 是 | 创建时间 | INDEX |
| updated_at | TIMESTAMP | 是 | 更新时间 | - |
| deleted_at | TIMESTAMP | 否 | 软删除时间 | INDEX |

**说明：**
- `listing_type` 区分出售和出租
- `price` 对于出售物业表示总价，对于出租物业表示月租
- `publisher_id` 可以是普通用户或地产代理公司
- `agent_id` 关联到具体的地产代理（可选）
- 支持按地区、校网、房间数等多维度筛选

**外键关系：**
- `publisher_id` → `users.id`
- `agent_id` → `agents.id`
- `district_id` → `districts.id`

---

### 2.2 房产图片表 (property_images)

存储房产的图片信息

| 字段名 | 类型 | 必填 | 说明 | 索引 |
|--------|------|------|------|------|
| id | BIGINT UNSIGNED | 是 | 图片ID（主键，自增） | PRIMARY |
| property_id | BIGINT UNSIGNED | 是 | 关联的房产ID | INDEX |
| image_url | VARCHAR(500) | 是 | 图片URL | - |
| image_type | VARCHAR(20) | 是 | 图片类型：cover=封面, interior=室内, exterior=外观, floorplan=户型图 | - |
| sort_order | INT | 是 | 排序顺序 | - |
| created_at | TIMESTAMP | 是 | 创建时间 | - |

**外键关系：**
- `property_id` → `properties.id` (CASCADE DELETE)

---

### 2.3 地区表 (districts)

香港地区分区表

| 字段名 | 类型 | 必填 | 说明 | 索引 |
|--------|------|------|------|------|
| id | BIGINT UNSIGNED | 是 | 地区ID（主键，自增） | PRIMARY |
| name_zh_hant | VARCHAR(100) | 是 | 中文繁体名称 | - |
| name_zh_hans | VARCHAR(100) | 否 | 中文简体名称 | - |
| name_en | VARCHAR(100) | 否 | 英文名称 | - |
| region | VARCHAR(50) | 是 | 区域：HK_ISLAND=港岛, KOWLOON=九龙, NEW_TERRITORIES=新界 | INDEX |
| sort_order | INT | 是 | 排序顺序 | - |
| created_at | TIMESTAMP | 是 | 创建时间 | - |
| updated_at | TIMESTAMP | 是 | 更新时间 | - |

---

### 2.4 房产设施表 (property_facilities)

房产设施多对多关联表

| 字段名 | 类型 | 必填 | 说明 | 索引 |
|--------|------|------|------|------|
| id | BIGINT UNSIGNED | 是 | ID（主键，自增） | PRIMARY |
| property_id | BIGINT UNSIGNED | 是 | 房产ID | INDEX |
| facility_id | BIGINT UNSIGNED | 是 | 设施ID | INDEX |
| created_at | TIMESTAMP | 是 | 创建时间 | - |

**联合唯一索引：**
- UNIQUE(`property_id`, `facility_id`)

**外键关系：**
- `property_id` → `properties.id` (CASCADE DELETE)
- `facility_id` → `facilities.id` (CASCADE DELETE)

---

### 2.5 设施字典表 (facilities)

可选设施列表（如：游泳池、健身房、停车场等）

| 字段名 | 类型 | 必填 | 说明 | 索引 |
|--------|------|------|------|------|
| id | BIGINT UNSIGNED | 是 | 设施ID（主键，自增） | PRIMARY |
| name_zh_hant | VARCHAR(100) | 是 | 中文繁体名称 | - |
| name_zh_hans | VARCHAR(100) | 否 | 中文简体名称 | - |
| name_en | VARCHAR(100) | 否 | 英文名称 | - |
| icon | VARCHAR(100) | 否 | 图标标识 | - |
| category | VARCHAR(50) | 是 | 分类：building=大厦设施, unit=单位设施 | INDEX |
| sort_order | INT | 是 | 排序顺序 | - |
| created_at | TIMESTAMP | 是 | 创建时间 | - |
| updated_at | TIMESTAMP | 是 | 更新时间 | - |

---

## 索引设计说明

### 用户表索引
- PRIMARY KEY: `id`
- UNIQUE KEY: `email`
- INDEX: `user_type`, `status`, `created_at`, `deleted_at`

### 房产表索引
- PRIMARY KEY: `id`
- UNIQUE KEY: `property_no`
- INDEX: `estate_no`, `listing_type`, `price`, `district_id`, `building_name`, `bedrooms`, `primary_school_net`, `secondary_school_net`, `property_type`, `status`, `publisher_id`, `agent_id`, `published_at`, `expired_at`, `created_at`, `deleted_at`
- COMPOSITE INDEX: (`listing_type`, `status`, `created_at`) - 用于列表查询

---

## 数据类型说明

- **BIGINT UNSIGNED**: 用于ID字段，支持大规模数据
- **VARCHAR**: 可变长度字符串
- **TEXT**: 长文本内容
- **DECIMAL**: 精确数值（价格、面积）
- **INT**: 整数（计数器、排序）
- **BOOLEAN**: 布尔值
- **TIMESTAMP**: 时间戳（支持时区）

---

## 3. 新盘模块

### 3.1 新盘表 (new_properties)

新楼盘信息表

| 字段名 | 类型 | 必填 | 说明 | 索引 |
|--------|------|------|------|------|
| id | BIGINT UNSIGNED | 是 | 新盘ID（主键，自增） | PRIMARY |
| name | VARCHAR(200) | 是 | 新盘名称 | INDEX |
| name_en | VARCHAR(200) | 否 | 英文名称 | - |
| address | VARCHAR(500) | 是 | 详细地址 | - |
| district_id | BIGINT UNSIGNED | 是 | 所属地区ID | INDEX |
| status | VARCHAR(20) | 是 | 状态：upcoming=即将推出, presale=预售中, selling=销售中, completed=已完成 | INDEX |
| units_for_sale | INT | 否 | 在售单位数 | - |
| units_sold | INT | 否 | 已售单位数 | - |
| developer | VARCHAR(200) | 是 | 开发商名称 | INDEX |
| management_company | VARCHAR(200) | 否 | 管理公司名称 | - |
| total_units | INT | 是 | 物业总伙数 | - |
| total_blocks | INT | 是 | 座数 | - |
| max_floors | INT | 是 | 最高层数 | - |
| primary_school_net | VARCHAR(50) | 否 | 小学校网 | INDEX |
| secondary_school_net | VARCHAR(50) | 否 | 中学校网 | INDEX |
| website_url | VARCHAR(500) | 否 | 官方网页地址 | - |
| sales_office_address | VARCHAR(500) | 否 | 销售处地址 | - |
| sales_phone | VARCHAR(50) | 否 | 销售电话 | - |
| expected_completion | DATE | 否 | 预计落成日期 | - |
| occupation_date | DATE | 否 | 入伙日期 | - |
| description | TEXT | 否 | 项目描述 | - |
| view_count | INT | 是 | 浏览次数 | - |
| favorite_count | INT | 是 | 收藏次数 | - |
| sort_order | INT | 是 | 排序顺序 | - |
| is_featured | BOOLEAN | 是 | 是否精选推荐 | INDEX |
| created_at | TIMESTAMP | 是 | 创建时间 | INDEX |
| updated_at | TIMESTAMP | 是 | 更新时间 | - |
| deleted_at | TIMESTAMP | 否 | 软删除时间 | INDEX |

**说明：**
- `units_for_sale` 和 `units_sold` 动态更新，反映实时销售情况
- 支持按开发商、地区、校网等筛选
- `is_featured` 用于首页推荐展示

**外键关系：**
- `district_id` → `districts.id`

---

### 3.2 新盘图片表 (new_property_images)

新盘项目图片

| 字段名 | 类型 | 必填 | 说明 | 索引 |
|--------|------|------|------|------|
| id | BIGINT UNSIGNED | 是 | 图片ID（主键，自增） | PRIMARY |
| new_property_id | BIGINT UNSIGNED | 是 | 关联的新盘ID | INDEX |
| image_url | VARCHAR(500) | 是 | 图片URL | - |
| image_type | VARCHAR(20) | 是 | 图片类型：exterior=外观, interior=室内示范单位, facilities=设施, floorplan=户型图, location=位置图 | - |
| title | VARCHAR(200) | 否 | 图片标题 | - |
| sort_order | INT | 是 | 排序顺序 | - |
| created_at | TIMESTAMP | 是 | 创建时间 | - |

**外键关系：**
- `new_property_id` → `new_properties.id` (CASCADE DELETE)

---

### 3.3 新盘户型表 (new_property_layouts)

新盘户型/价单信息

| 字段名 | 类型 | 必填 | 说明 | 索引 |
|--------|------|------|------|------|
| id | BIGINT UNSIGNED | 是 | 户型ID（主键，自增） | PRIMARY |
| new_property_id | BIGINT UNSIGNED | 是 | 关联的新盘ID | INDEX |
| unit_type | VARCHAR(50) | 是 | 户型类型（如：1房、2房、3房） | - |
| bedrooms | INT | 是 | 房间数 | - |
| bathrooms | INT | 否 | 浴室数 | - |
| saleable_area | DECIMAL(10,2) | 是 | 实用面积（平方尺） | - |
| gross_area | DECIMAL(10,2) | 否 | 建筑面积（平方尺） | - |
| min_price | DECIMAL(15,2) | 是 | 最低售价（港币） | - |
| max_price | DECIMAL(15,2) | 否 | 最高售价（港币） | - |
| price_per_sqft | DECIMAL(10,2) | 否 | 每平方尺价格 | - |
| available_units | INT | 是 | 可售单位数 | - |
| floorplan_url | VARCHAR(500) | 否 | 户型图URL | - |
| created_at | TIMESTAMP | 是 | 创建时间 | - |
| updated_at | TIMESTAMP | 是 | 更新时间 | - |

**外键关系：**
- `new_property_id` → `new_properties.id` (CASCADE DELETE)

---

## 4. 服务式住宅模块

### 4.1 服务式住宅表 (serviced_apartments)

服务式公寓/住宅信息表

| 字段名 | 类型 | 必填 | 说明 | 索引 |
|--------|------|------|------|------|
| id | BIGINT UNSIGNED | 是 | 服务式住宅ID（主键，自增） | PRIMARY |
| name | VARCHAR(200) | 是 | 住宅名称 | INDEX |
| name_en | VARCHAR(200) | 否 | 英文名称 | - |
| address | VARCHAR(500) | 是 | 详细地址 | - |
| district_id | BIGINT UNSIGNED | 是 | 所属地区ID | INDEX |
| description | TEXT | 否 | 详细描述 | - |
| phone | VARCHAR(50) | 是 | 联系电话 | - |
| website_url | VARCHAR(500) | 否 | 官方网站 | - |
| email | VARCHAR(255) | 否 | 联系邮箱 | - |
| company_id | BIGINT UNSIGNED | 是 | 所属公司ID（关联users表，必须是agency类型） | INDEX |
| check_in_time | VARCHAR(50) | 否 | 入住时间 | - |
| check_out_time | VARCHAR(50) | 否 | 退房时间 | - |
| min_stay_days | INT | 否 | 最少入住天数 | - |
| status | VARCHAR(20) | 是 | 状态：active=营业中, inactive=暂停营业, closed=已关闭 | INDEX |
| rating | DECIMAL(3,2) | 否 | 评分（0-5） | - |
| review_count | INT | 是 | 评价数量 | - |
| view_count | INT | 是 | 浏览次数 | - |
| favorite_count | INT | 是 | 收藏次数 | - |
| is_featured | BOOLEAN | 是 | 是否精选推荐 | INDEX |
| created_at | TIMESTAMP | 是 | 创建时间 | INDEX |
| updated_at | TIMESTAMP | 是 | 更新时间 | - |
| deleted_at | TIMESTAMP | 否 | 软删除时间 | INDEX |

**说明：**
- 服务式住宅由对应的地产代理公司发布和管理
- `company_id` 必须关联到 `user_type='agency'` 的用户
- 支持评分和评价功能

**外键关系：**
- `district_id` → `districts.id`
- `company_id` → `users.id` (WHERE user_type='agency')

---

### 4.2 服务式住宅房型表 (serviced_apartment_units)

服务式住宅的具体房型和价格

| 字段名 | 类型 | 必填 | 说明 | 索引 |
|--------|------|------|------|------|
| id | BIGINT UNSIGNED | 是 | 房型ID（主键，自增） | PRIMARY |
| serviced_apartment_id | BIGINT UNSIGNED | 是 | 关联的服务式住宅ID | INDEX |
| unit_type | VARCHAR(50) | 是 | 房型名称（如：标准房、豪华房、套房） | - |
| bedrooms | INT | 是 | 房间数 | - |
| bathrooms | INT | 否 | 浴室数 | - |
| area | DECIMAL(10,2) | 是 | 面积（平方尺） | - |
| max_occupancy | INT | 是 | 最多入住人数 | - |
| daily_price | DECIMAL(10,2) | 否 | 日租价格（港币） | - |
| weekly_price | DECIMAL(10,2) | 否 | 周租价格（港币） | - |
| monthly_price | DECIMAL(10,2) | 是 | 月租价格（港币） | INDEX |
| available_units | INT | 是 | 可用单位数 | - |
| description | TEXT | 否 | 房型描述 | - |
| sort_order | INT | 是 | 排序顺序 | - |
| created_at | TIMESTAMP | 是 | 创建时间 | - |
| updated_at | TIMESTAMP | 是 | 更新时间 | - |

**外键关系：**
- `serviced_apartment_id` → `serviced_apartments.id` (CASCADE DELETE)

---

### 4.3 服务式住宅图片表 (serviced_apartment_images)

服务式住宅图片

| 字段名 | 类型 | 必填 | 说明 | 索引 |
|--------|------|------|------|------|
| id | BIGINT UNSIGNED | 是 | 图片ID（主键，自增） | PRIMARY |
| serviced_apartment_id | BIGINT UNSIGNED | 否 | 关联的服务式住宅ID（整体照片） | INDEX |
| unit_id | BIGINT UNSIGNED | 否 | 关联的房型ID（房型照片） | INDEX |
| image_url | VARCHAR(500) | 是 | 图片URL | - |
| image_type | VARCHAR(20) | 是 | 图片类型：exterior=外观, lobby=大堂, room=房间, bathroom=浴室, facilities=设施 | - |
| title | VARCHAR(200) | 否 | 图片标题 | - |
| sort_order | INT | 是 | 排序顺序 | - |
| created_at | TIMESTAMP | 是 | 创建时间 | - |

**说明：**
- `serviced_apartment_id` 和 `unit_id` 至少有一个不为空
- 支持整体项目照片和具体房型照片

**外键关系：**
- `serviced_apartment_id` → `serviced_apartments.id` (CASCADE DELETE)
- `unit_id` → `serviced_apartment_units.id` (CASCADE DELETE)

---

### 4.4 服务式住宅设施表 (serviced_apartment_facilities)

服务式住宅设施多对多关联表

| 字段名 | 类型 | 必填 | 说明 | 索引 |
|--------|------|------|------|------|
| id | BIGINT UNSIGNED | 是 | ID（主键，自增） | PRIMARY |
| serviced_apartment_id | BIGINT UNSIGNED | 是 | 服务式住宅ID | INDEX |
| facility_id | BIGINT UNSIGNED | 是 | 设施ID | INDEX |
| created_at | TIMESTAMP | 是 | 创建时间 | - |

**联合唯一索引：**
- UNIQUE(`serviced_apartment_id`, `facility_id`)

**外键关系：**
- `serviced_apartment_id` → `serviced_apartments.id` (CASCADE DELETE)
- `facility_id` → `facilities.id` (CASCADE DELETE)

---

## 5. 屋苑/小区模块

### 5.1 屋苑表 (estates)

屋苑/小区信息表

| 字段名 | 类型 | 必填 | 说明 | 索引 |
|--------|------|------|------|------|
| id | BIGINT UNSIGNED | 是 | 屋苑ID（主键，自增） | PRIMARY |
| name | VARCHAR(200) | 是 | 屋苑名称 | INDEX |
| name_en | VARCHAR(200) | 否 | 英文名称 | - |
| address | VARCHAR(500) | 是 | 详细地址 | - |
| district_id | BIGINT UNSIGNED | 是 | 所属地区ID | INDEX |
| total_blocks | INT | 否 | 总座数 | - |
| total_units | INT | 否 | 总单位数 | - |
| completion_year | INT | 否 | 落成年份 | - |
| developer | VARCHAR(200) | 否 | 发展商 | - |
| management_company | VARCHAR(200) | 否 | 管理公司 | - |
| primary_school_net | VARCHAR(50) | 否 | 小学校网 | INDEX |
| secondary_school_net | VARCHAR(50) | 否 | 中学校网 | INDEX |
| recent_transactions_count | INT | 是 | 近期成交数量（最近3个月） | - |
| for_sale_count | INT | 是 | 当前放盘数量 | - |
| for_rent_count | INT | 是 | 当前租盘数量 | - |
| avg_transaction_price | DECIMAL(15,2) | 否 | 平均成交价（港币/平方尺） | INDEX |
| avg_transaction_price_updated_at | TIMESTAMP | 否 | 平均成交价更新时间 | - |
| description | TEXT | 否 | 屋苑描述 | - |
| view_count | INT | 是 | 浏览次数 | - |
| favorite_count | INT | 是 | 收藏次数 | - |
| is_featured | BOOLEAN | 是 | 是否精选屋苑 | INDEX |
| created_at | TIMESTAMP | 是 | 创建时间 | INDEX |
| updated_at | TIMESTAMP | 是 | 更新时间 | - |
| deleted_at | TIMESTAMP | 否 | 软删除时间 | INDEX |

**说明：**
- `recent_transactions_count`, `for_sale_count`, `for_rent_count` 通过定时任务统计更新
- `avg_transaction_price` 基于近期成交记录计算
- 房产表中的 `building_name` 可关联到此表的 `name`

**外键关系：**
- `district_id` → `districts.id`

---

### 5.2 屋苑图片表 (estate_images)

屋苑相关图片

| 字段名 | 类型 | 必填 | 说明 | 索引 |
|--------|------|------|------|------|
| id | BIGINT UNSIGNED | 是 | 图片ID（主键，自增） | PRIMARY |
| estate_id | BIGINT UNSIGNED | 是 | 关联的屋苑ID | INDEX |
| image_url | VARCHAR(500) | 是 | 图片URL | - |
| image_type | VARCHAR(20) | 是 | 图片类型：exterior=外观, facilities=设施, environment=环境, aerial=航拍 | - |
| title | VARCHAR(200) | 否 | 图片标题 | - |
| sort_order | INT | 是 | 排序顺序 | - |
| created_at | TIMESTAMP | 是 | 创建时间 | - |

**外键关系：**
- `estate_id` → `estates.id` (CASCADE DELETE)

---

### 5.3 屋苑设施表 (estate_facilities)

屋苑设施多对多关联表

| 字段名 | 类型 | 必填 | 说明 | 索引 |
|--------|------|------|------|------|
| id | BIGINT UNSIGNED | 是 | ID（主键，自增） | PRIMARY |
| estate_id | BIGINT UNSIGNED | 是 | 屋苑ID | INDEX |
| facility_id | BIGINT UNSIGNED | 是 | 设施ID | INDEX |
| created_at | TIMESTAMP | 是 | 创建时间 | - |

**联合唯一索引：**
- UNIQUE(`estate_id`, `facility_id`)

**外键关系：**
- `estate_id` → `estates.id` (CASCADE DELETE)
- `facility_id` → `facilities.id` (CASCADE DELETE)

---

## 6. 家具模块

### 6.1 家具表 (furniture)

二手家具/家居用品交易表

| 字段名 | 类型 | 必填 | 说明 | 索引 |
|--------|------|------|------|------|
| id | BIGINT UNSIGNED | 是 | 家具ID（主键，自增） | PRIMARY |
| furniture_no | VARCHAR(50) | 是 | 家具编号（系统生成） | UNIQUE |
| title | VARCHAR(255) | 是 | 家具名称/标题 | INDEX |
| description | TEXT | 否 | 详细描述 | - |
| price | DECIMAL(10,2) | 是 | 价格（港币） | INDEX |
| category_id | BIGINT UNSIGNED | 是 | 分类ID | INDEX |
| brand | VARCHAR(100) | 否 | 品牌 | INDEX |
| condition | VARCHAR(20) | 是 | 新旧程度：new=全新, like_new=近全新, good=良好, fair=一般, poor=较差 | INDEX |
| purchase_date | DATE | 否 | 购买日期 | - |
| delivery_district_id | BIGINT UNSIGNED | 是 | 交收地区ID | INDEX |
| delivery_time | VARCHAR(100) | 否 | 交收时间（如：工作日晚上、周末全天） | - |
| delivery_method | VARCHAR(50) | 是 | 交收方法：self_pickup=自取, delivery=送货, negotiable=面议 | INDEX |
| status | VARCHAR(20) | 是 | 状态：available=可用, reserved=已预订, sold=已售出, expired=已过期, cancelled=已取消 | INDEX |
| publisher_id | BIGINT UNSIGNED | 是 | 发布者ID（关联users表） | INDEX |
| publisher_type | VARCHAR(20) | 是 | 发布者类型：individual=个人, agency=代理公司 | - |
| view_count | INT | 是 | 浏览次数 | - |
| favorite_count | INT | 是 | 收藏次数 | - |
| published_at | TIMESTAMP | 是 | 刊登日期 | INDEX |
| updated_at | TIMESTAMP | 是 | 更新日期 | INDEX |
| expires_at | TIMESTAMP | 是 | 到期日期 | INDEX |
| created_at | TIMESTAMP | 是 | 创建时间 | - |
| deleted_at | TIMESTAMP | 否 | 软删除时间 | INDEX |

**说明：**
- 家具可以由普通用户或地产代理公司发布
- `expires_at` 默认为发布后90天，到期后状态自动变为expired
- 支持按分类、品牌、新旧程度、地区筛选

**外键关系：**
- `publisher_id` → `users.id`
- `category_id` → `furniture_categories.id`
- `delivery_district_id` → `districts.id`

---

### 6.2 家具分类表 (furniture_categories)

家具分类字典表

| 字段名 | 类型 | 必填 | 说明 | 索引 |
|--------|------|------|------|------|
| id | BIGINT UNSIGNED | 是 | 分类ID（主键，自增） | PRIMARY |
| parent_id | BIGINT UNSIGNED | 否 | 父分类ID（支持二级分类） | INDEX |
| name_zh_hant | VARCHAR(100) | 是 | 中文繁体名称 | - |
| name_zh_hans | VARCHAR(100) | 否 | 中文简体名称 | - |
| name_en | VARCHAR(100) | 否 | 英文名称 | - |
| icon | VARCHAR(100) | 否 | 图标标识 | - |
| sort_order | INT | 是 | 排序顺序 | - |
| is_active | BOOLEAN | 是 | 是否启用 | - |
| created_at | TIMESTAMP | 是 | 创建时间 | - |
| updated_at | TIMESTAMP | 是 | 更新时间 | - |

**分类示例：**
- 客厅家具（沙发、茶几、电视柜）
- 卧室家具（床、衣柜、床头柜）
- 厨房用品（餐桌、餐椅、厨柜）
- 办公家具（办公桌、办公椅、文件柜）
- 家电（电视、冰箱、洗衣机）
- 装饰品（灯具、窗帘、地毯）

**外键关系：**
- `parent_id` → `furniture_categories.id`

---

### 6.3 家具图片表 (furniture_images)

家具图片表

| 字段名 | 类型 | 必填 | 说明 | 索引 |
|--------|------|------|------|------|
| id | BIGINT UNSIGNED | 是 | 图片ID（主键，自增） | PRIMARY |
| furniture_id | BIGINT UNSIGNED | 是 | 关联的家具ID | INDEX |
| image_url | VARCHAR(500) | 是 | 图片URL | - |
| is_cover | BOOLEAN | 是 | 是否为封面图 | - |
| sort_order | INT | 是 | 排序顺序 | - |
| created_at | TIMESTAMP | 是 | 创建时间 | - |

**说明：**
- 每件家具至少要有1张图片
- 只能有1张封面图（is_cover=true）

**外键关系：**
- `furniture_id` → `furniture.id` (CASCADE DELETE)

---

### 6.4 家具订单表 (furniture_orders)

家具交易订单表

| 字段名 | 类型 | 必填 | 说明 | 索引 |
|--------|------|------|------|------|
| id | BIGINT UNSIGNED | 是 | 订单ID（主键，自增） | PRIMARY |
| order_no | VARCHAR(50) | 是 | 订单号（系统生成） | UNIQUE |
| buyer_id | BIGINT UNSIGNED | 是 | 买家ID | INDEX |
| seller_id | BIGINT UNSIGNED | 是 | 卖家ID | INDEX |
| furniture_id | BIGINT UNSIGNED | 是 | 家具ID | INDEX |
| price | DECIMAL(10,2) | 是 | 成交价格 | - |
| status | VARCHAR(20) | 是 | 订单状态：pending=待确认, confirmed=已确认, completed=已完成, cancelled=已取消 | INDEX |
| delivery_method | VARCHAR(50) | 是 | 交收方式 | - |
| delivery_address | VARCHAR(500) | 否 | 送货地址 | - |
| delivery_date | DATE | 否 | 交收日期 | - |
| buyer_note | TEXT | 否 | 买家备注 | - |
| seller_note | TEXT | 否 | 卖家备注 | - |
| created_at | TIMESTAMP | 是 | 创建时间 | INDEX |
| confirmed_at | TIMESTAMP | 否 | 确认时间 | - |
| completed_at | TIMESTAMP | 否 | 完成时间 | - |
| cancelled_at | TIMESTAMP | 否 | 取消时间 | - |

**外键关系：**
- `buyer_id` → `users.id`
- `seller_id` → `users.id`
- `furniture_id` → `furniture.id`

---

## 7. 地产代理模块

### 7.1 地产代理表 (agents)

个人地产代理信息表

| 字段名 | 类型 | 必填 | 说明 | 索引 |
|--------|------|------|------|------|
| id | BIGINT UNSIGNED | 是 | 代理ID（主键，自增） | PRIMARY |
| user_id | BIGINT UNSIGNED | 是 | 关联的用户ID | UNIQUE |
| agent_name | VARCHAR(100) | 是 | 代理人姓名 | INDEX |
| agent_name_en | VARCHAR(100) | 否 | 英文姓名 | - |
| license_no | VARCHAR(50) | 是 | 地产代理牌照号码 | UNIQUE |
| license_type | VARCHAR(20) | 是 | 牌照类型：individual=个人牌照, salesperson=营业员牌照 | INDEX |
| license_expiry_date | DATE | 否 | 牌照到期日期 | - |
| agency_id | BIGINT UNSIGNED | 否 | 所属代理公司ID（关联users表，user_type='agency'） | INDEX |
| phone | VARCHAR(20) | 是 | 联系电话 | - |
| mobile | VARCHAR(20) | 否 | 手机号码 | - |
| email | VARCHAR(255) | 是 | 电子邮箱 | INDEX |
| wechat_id | VARCHAR(50) | 否 | 微信号 | - |
| whatsapp | VARCHAR(20) | 否 | WhatsApp号码 | - |
| office_address | VARCHAR(500) | 否 | 办公地址 | - |
| specialization | VARCHAR(200) | 否 | 专长领域（如：豪宅、工商铺、新盘） | - |
| years_experience | INT | 否 | 从业年限 | - |
| profile_photo | VARCHAR(500) | 否 | 个人照片URL | - |
| bio | TEXT | 否 | 个人简介 | - |
| rating | DECIMAL(3,2) | 否 | 评分（0-5） | INDEX |
| review_count | INT | 是 | 评价数量 | - |
| properties_sold | INT | 是 | 已售物业数量 | - |
| properties_rented | INT | 是 | 已租物业数量 | - |
| status | VARCHAR(20) | 是 | 状态：active=活跃, inactive=停用, suspended=暂停 | INDEX |
| is_verified | BOOLEAN | 是 | 是否已验证牌照 | INDEX |
| verified_at | TIMESTAMP | 否 | 验证时间 | - |
| created_at | TIMESTAMP | 是 | 创建时间 | INDEX |
| updated_at | TIMESTAMP | 是 | 更新时间 | - |
| deleted_at | TIMESTAMP | 否 | 软删除时间 | INDEX |

**说明：**
- 地产代理必须关联一个用户账号（`user_id`）
- 可选择加入某个代理公司（`agency_id`），或独立执业
- `license_no` 必须唯一，用于验证代理资格
- `is_verified` 表示平台是否已验证其牌照真实性

**外键关系：**
- `user_id` → `users.id`
- `agency_id` → `users.id` (WHERE user_type='agency')

---

### 7.2 代理公司详情表 (agency_details)

地产代理公司的详细信息（扩展users表）

| 字段名 | 类型 | 必填 | 说明 | 索引 |
|--------|------|------|------|------|
| id | BIGINT UNSIGNED | 是 | ID（主键，自增） | PRIMARY |
| user_id | BIGINT UNSIGNED | 是 | 关联的用户ID（user_type='agency'） | UNIQUE |
| company_name | VARCHAR(200) | 是 | 公司名称 | INDEX |
| company_name_en | VARCHAR(200) | 否 | 英文公司名称 | - |
| license_no | VARCHAR(50) | 是 | 公司牌照号码 | UNIQUE |
| business_registration_no | VARCHAR(50) | 否 | 商业登记号码 | - |
| address | VARCHAR(500) | 是 | 公司地址 | - |
| phone | VARCHAR(20) | 是 | 公司电话 | - |
| fax | VARCHAR(20) | 否 | 传真号码 | - |
| email | VARCHAR(255) | 是 | 公司邮箱 | - |
| website_url | VARCHAR(500) | 否 | 公司网站 | - |
| established_year | INT | 否 | 成立年份 | - |
| agent_count | INT | 是 | 旗下代理人数 | - |
| description | TEXT | 否 | 公司简介 | - |
| logo_url | VARCHAR(500) | 否 | 公司Logo URL | - |
| cover_image_url | VARCHAR(500) | 否 | 封面图URL | - |
| rating | DECIMAL(3,2) | 否 | 评分（0-5） | INDEX |
| review_count | INT | 是 | 评价数量 | - |
| is_verified | BOOLEAN | 是 | 是否已验证 | INDEX |
| verified_at | TIMESTAMP | 否 | 验证时间 | - |
| created_at | TIMESTAMP | 是 | 创建时间 | - |
| updated_at | TIMESTAMP | 是 | 更新时间 | - |

**说明：**
- 此表扩展 `users` 表中 `user_type='agency'` 的详细信息
- `agent_count` 通过统计 `agents` 表中 `agency_id` 自动更新

**外键关系：**
- `user_id` → `users.id` (WHERE user_type='agency')

---

### 7.3 代理服务区域表 (agent_service_areas)

地产代理服务的地区

| 字段名 | 类型 | 必填 | 说明 | 索引 |
|--------|------|------|------|------|
| id | BIGINT UNSIGNED | 是 | ID（主键，自增） | PRIMARY |
| agent_id | BIGINT UNSIGNED | 是 | 代理ID | INDEX |
| district_id | BIGINT UNSIGNED | 是 | 服务地区ID | INDEX |
| created_at | TIMESTAMP | 是 | 创建时间 | - |

**联合唯一索引：**
- UNIQUE(`agent_id`, `district_id`)

**外键关系：**
- `agent_id` → `agents.id` (CASCADE DELETE)
- `district_id` → `districts.id`

---

## 8. 业务关系与权限说明

### 8.1 发布权限矩阵

| 内容类型 | 普通用户 (individual) | 地产代理公司 (agency) | 地产代理 (agent) |
|----------|----------------------|---------------------|-----------------|
| 房源信息 (properties) | ✅ 可发布 | ✅ 可发布 | ✅ 可发布（代表所属公司） |
| 服务式住宅 (serviced_apartments) | ❌ 不可发布 | ✅ 可发布（仅自己的） | ❌ 不可发布 |
| 家具 (furniture) | ✅ 可发布 | ✅ 可发布 | ✅ 可发布（代表个人或公司） |
| 新盘 (new_properties) | ❌ 不可发布 | ✅ 可发布（管理员审核） | ❌ 不可发布 |

### 8.2 数据关联关系

#### 房源发布逻辑
```
properties 表：
- publisher_id → users.id (普通用户或代理公司)
- publisher_type：individual / agency
- agent_id → agents.id (可选，具体负责的代理人)

发布场景：
1. 普通用户发布：publisher_type='individual', agent_id=NULL
2. 代理公司发布：publisher_type='agency', agent_id 可指定具体代理人
3. 代理人代表公司发布：publisher_type='agency', publisher_id=代理公司ID, agent_id=代理人ID
```

#### 服务式住宅发布逻辑
```
serviced_apartments 表：
- company_id → users.id (必须是 user_type='agency')

发布场景：
- 仅代理公司可发布和管理自己的服务式住宅
```

#### 家具发布逻辑
```
furniture 表：
- publisher_id → users.id (普通用户或代理公司)
- publisher_type：individual / agency

发布场景：
1. 普通用户发布：publisher_type='individual'
2. 代理公司发布：publisher_type='agency'
3. 代理人发布：publisher_id=代理人的user_id, publisher_type='individual'
```

### 8.3 用户类型关系图

```
users (user_type='individual')
    ├─ 可注册成为 → agents (地产代理)
    │                      └─ 可加入 → agency_details (代理公司)
    ├─ 可发布 → properties (房源)
    └─ 可发布 → furniture (家具)

users (user_type='agency')
    ├─ 关联 → agency_details (公司详情)
    ├─ 旗下有多个 → agents (地产代理)
    ├─ 可发布 → properties (房源)
    ├─ 可发布 → serviced_apartments (服务式住宅)
    └─ 可发布 → furniture (家具)

agents (地产代理)
    ├─ 关联 → users (user_id)
    ├─ 可选加入 → users[agency] (agency_id)
    ├─ 负责 → properties (agent_id)
    └─ 服务区域 → agent_service_areas → districts
```

---

## 待补充模块

后续将补充以下模块的设计：
- 成交记录模块
- 收藏/浏览历史模块
- 购物车模块
- 评论/评价模块
- 消息通知模块

---

## 更新日志

| 日期 | 版本 | 说明 |
|------|------|------|
| 2025-12-18 | v0.1 | 初始版本 - 用户模块和房产模块 |
| 2025-12-18 | v0.2 | 新增新盘模块和服务式住宅模块 |
| 2025-12-18 | v0.3 | 新增屋苑/小区模块和家具模块 |
| 2025-12-18 | v0.4 | 新增地产代理模块和业务关系说明 |
