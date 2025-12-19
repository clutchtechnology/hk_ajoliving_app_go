# 按揭模块 API 实现总结

## 概述

已完成按揭模块的 7 个 API 接口实现，包括按揭计算、利率查询、利率比较和申请管理功能。

## 已实现的 API 接口

### 1. 计算按揭 (POST /api/v1/mortgage/calculate)
- **功能**: 根据物业价格、首付、利率和还款期计算月供及还款计划
- **权限**: 公开（无需认证）
- **请求参数**:
  - `property_price`: 物业价格
  - `down_payment`: 首付金额
  - `interest_rate`: 年利率（百分比）
  - `loan_period`: 还款期（月数）
- **返回数据**:
  - 贷款金额、月供、总还款额、总利息
  - 还款计划（前12个月和最后1个月）
- **特点**: 使用 PMT 公式计算月供

### 2. 获取按揭利率列表 (GET /api/v1/mortgage/rates)
- **功能**: 获取所有有效的按揭利率
- **权限**: 公开（无需认证）
- **返回数据**: 所有银行的有效利率方案列表

### 3. 获取指定银行的按揭利率 (GET /api/v1/mortgage/rates/bank/:bank_id)
- **功能**: 根据银行ID获取该银行的所有按揭利率方案
- **权限**: 公开（无需认证）
- **路径参数**: `bank_id` - 银行ID
- **返回数据**: 指定银行的利率方案列表

### 4. 比较按揭利率 (POST /api/v1/mortgage/rates/compare)
- **功能**: 比较不同银行的按揭利率和月供金额
- **权限**: 公开（无需认证）
- **请求参数**:
  - `property_price`: 物业价格
  - `down_payment`: 首付金额
  - `loan_period`: 还款期（月数）
  - `rate_type`: 利率类型（fixed/floating/hibor）
- **返回数据**:
  - 各银行利率方案的月供、总还款额、总利息
  - 最低利率方案
  - 与最低利率的差额

### 5. 申请按揭 (POST /api/v1/mortgage/apply)
- **功能**: 提交按揭申请
- **权限**: 需要认证（JWT）
- **请求参数**:
  - 物业信息: `property_id`, `property_price`, `down_payment`
  - 贷款信息: `bank_id`, `interest_rate`, `loan_period`
  - 申请人信息: `applicant_name`, `applicant_phone`, `applicant_email`, `applicant_id_card`
  - 财务信息: `monthly_income`, `employment_status`, `employer_name`, `years_employed`
- **返回数据**: 申请详情（包含自动生成的申请编号）

### 6. 获取按揭申请列表 (GET /api/v1/mortgage/applications)
- **功能**: 获取当前用户的按揭申请列表
- **权限**: 需要认证（JWT）
- **查询参数**:
  - `status`: 状态筛选（pending/under_review/approved/rejected/withdrawn）
  - `bank_id`: 银行ID筛选
  - `page`: 页码
  - `page_size`: 每页数量
  - `sort_by`: 排序字段
  - `sort_order`: 排序方向
- **返回数据**: 分页的申请列表

### 7. 获取按揭申请详情 (GET /api/v1/mortgage/applications/:id)
- **功能**: 根据申请ID获取详细信息
- **权限**: 需要认证（JWT），仅限申请人本人
- **路径参数**: `id` - 申请ID
- **返回数据**: 申请的完整详情

## 数据模型

### Bank (银行)
- 基本信息: 中英文名称、银行代码
- 多语言支持: 繁体中文、简体中文、英文
- Logo 和网站 URL
- 排序和激活状态

### MortgageRate (按揭利率)
- 利率信息: 利率类型（定息/浮息/HIBOR）、利率值
- 贷款条件: 最低/最高贷款金额、贷款期限、按揭成数
- 有效期: 生效日期、到期日期
- 促销标识

### MortgageApplication (按揭申请)
- 申请编号: 自动生成（MTG + 日期 + 序号）
- 物业和贷款信息
- 申请人信息
- 财务信息
- 状态跟踪: pending → under_review → approved/rejected/withdrawn
- 时间戳: 提交、审核、批准、拒绝、撤回时间

## 业务逻辑特性

### 按揭计算
- **PMT 公式**: `PMT = P * [r(1+r)^n] / [(1+r)^n - 1]`
  - P = 贷款本金
  - r = 月利率
  - n = 还款月数
- **还款计划**: 显示前12个月和最后1个月的详细还款信息

### 利率比较
- 自动识别最低利率方案
- 计算每个方案与最低利率的差额
- 支持按利率类型筛选

### 申请管理
- 申请编号自动生成
- 状态工作流管理
- 所有权验证

## 文件结构

```
internal/
├── model/
│   └── mortgage.go                    # Bank, MortgageRate, MortgageApplication 模型
├── dto/
│   ├── request/
│   │   └── mortgage_request.go        # 请求 DTO
│   └── response/
│       └── mortgage_response.go       # 响应 DTO
├── repository/
│   └── mortgage_repository.go         # 数据访问层
├── service/
│   └── mortgage_service.go            # 业务逻辑层
└── handler/
    └── mortgage_handler.go            # HTTP 处理器
```

## 路由配置

```
/api/v1/mortgage
├── POST   /calculate                  # 计算按揭（公开）
├── GET    /rates                      # 获取利率列表（公开）
├── GET    /rates/bank/:bank_id        # 获取银行利率（公开）
├── POST   /rates/compare              # 比较利率（公开）
├── POST   /apply                      # 申请按揭（需认证）
├── GET    /applications               # 获取申请列表（需认证）
└── GET    /applications/:id           # 获取申请详情（需认证）
```

## 数据库迁移

已在 `main.go` 的 `autoMigrate` 中添加以下模型:
- `model.Bank`
- `model.MortgageRate`
- `model.MortgageApplication`

## API 文档

所有接口已添加 Swagger 注解，包括:
- 接口描述
- 参数说明
- 响应格式
- 错误代码

## 测试建议

### 1. 按揭计算测试
```bash
curl -X POST http://localhost:8080/api/v1/mortgage/calculate \
  -H "Content-Type: application/json" \
  -d '{
    "property_price": 5000000,
    "down_payment": 1000000,
    "interest_rate": 2.5,
    "loan_period": 360
  }'
```

### 2. 利率比较测试
```bash
curl -X POST http://localhost:8080/api/v1/mortgage/rates/compare \
  -H "Content-Type: application/json" \
  -d '{
    "property_price": 5000000,
    "down_payment": 1000000,
    "loan_period": 360,
    "rate_type": "floating"
  }'
```

### 3. 申请按揭测试（需要 JWT Token）
```bash
curl -X POST http://localhost:8080/api/v1/mortgage/apply \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "property_id": 1,
    "bank_id": 1,
    "property_price": 5000000,
    "down_payment": 1000000,
    "interest_rate": 2.5,
    "loan_period": 360,
    "monthly_income": 50000,
    "employment_status": "employed",
    "employer_name": "ABC Company",
    "years_employed": 5,
    "applicant_name": "张三",
    "applicant_phone": "98765432",
    "applicant_email": "zhang@example.com",
    "applicant_id_card": "A123456(7)"
  }'
```

## 注意事项

1. **生产环境**: 需要创建数据库迁移文件（在 `migrations/` 目录）
2. **银行数据**: 需要预先填充银行信息和利率数据
3. **利率更新**: 建议定期更新利率数据以保持准确性
4. **申请审核**: 目前状态更新需要后台管理接口（未包含在此次实现中）
5. **PMT 计算**: 已实现标准等额本息还款法

## 下一步工作

1. 创建数据库迁移文件
2. 添加银行和利率数据的初始化脚本
3. 实现后台管理接口（用于审核申请、更新利率等）
4. 添加单元测试和集成测试
5. 生成 Swagger 文档

## 编译状态

✅ 无编译错误
✅ 所有依赖已正确导入
✅ 路由已正确注册
✅ 数据库模型已添加到 AutoMigrate
