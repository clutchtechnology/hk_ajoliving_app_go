-- 创建地区表
CREATE TABLE IF NOT EXISTS districts (
    id SERIAL PRIMARY KEY,
    name_zh_hant VARCHAR(100) NOT NULL,
    name_zh_hans VARCHAR(100),
    name_en VARCHAR(100),
    region VARCHAR(50) NOT NULL,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 创建设施表
CREATE TABLE IF NOT EXISTS facilities (
    id SERIAL PRIMARY KEY,
    name_zh_hant VARCHAR(100) NOT NULL,
    name_zh_hans VARCHAR(100),
    name_en VARCHAR(100),
    icon VARCHAR(100),
    category VARCHAR(50) NOT NULL,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 创建代理人表
CREATE TABLE IF NOT EXISTS agents (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    name VARCHAR(100) NOT NULL,
    phone VARCHAR(20),
    email VARCHAR(255),
    avatar VARCHAR(500),
    license_no VARCHAR(50),
    agency_id INT,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- 创建房产表
CREATE TABLE IF NOT EXISTS properties (
    id SERIAL PRIMARY KEY,
    property_no VARCHAR(50) NOT NULL UNIQUE,
    estate_no VARCHAR(50),
    listing_type VARCHAR(20) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    area DECIMAL(10,2) NOT NULL,
    price DECIMAL(15,2) NOT NULL,
    address VARCHAR(500) NOT NULL,
    district_id INT NOT NULL REFERENCES districts(id),
    building_name VARCHAR(200),
    floor VARCHAR(20),
    orientation VARCHAR(50),
    bedrooms INT NOT NULL DEFAULT 0,
    bathrooms INT,
    primary_school_net VARCHAR(50),
    secondary_school_net VARCHAR(50),
    property_type VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'available',
    publisher_id INT NOT NULL REFERENCES users(id),
    publisher_type VARCHAR(20) NOT NULL,
    agent_id INT REFERENCES agents(id),
    view_count INT NOT NULL DEFAULT 0,
    favorite_count INT NOT NULL DEFAULT 0,
    published_at TIMESTAMP WITH TIME ZONE,
    expired_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- 创建房产图片表
CREATE TABLE IF NOT EXISTS property_images (
    id SERIAL PRIMARY KEY,
    property_id INT NOT NULL REFERENCES properties(id) ON DELETE CASCADE,
    url VARCHAR(500) NOT NULL,
    caption VARCHAR(255),
    sort_order INT NOT NULL DEFAULT 0,
    is_cover BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 创建房产设施关联表
CREATE TABLE IF NOT EXISTS property_facilities (
    property_id INT NOT NULL REFERENCES properties(id) ON DELETE CASCADE,
    facility_id INT NOT NULL REFERENCES facilities(id) ON DELETE CASCADE,
    PRIMARY KEY (property_id, facility_id)
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_properties_listing_type ON properties(listing_type);
CREATE INDEX IF NOT EXISTS idx_properties_district_id ON properties(district_id);
CREATE INDEX IF NOT EXISTS idx_properties_price ON properties(price);
CREATE INDEX IF NOT EXISTS idx_properties_bedrooms ON properties(bedrooms);
CREATE INDEX IF NOT EXISTS idx_properties_property_type ON properties(property_type);
CREATE INDEX IF NOT EXISTS idx_properties_status ON properties(status);
CREATE INDEX IF NOT EXISTS idx_properties_publisher_id ON properties(publisher_id);
CREATE INDEX IF NOT EXISTS idx_properties_building_name ON properties(building_name);
CREATE INDEX IF NOT EXISTS idx_properties_primary_school_net ON properties(primary_school_net);
CREATE INDEX IF NOT EXISTS idx_properties_secondary_school_net ON properties(secondary_school_net);
CREATE INDEX IF NOT EXISTS idx_properties_created_at ON properties(created_at);
CREATE INDEX IF NOT EXISTS idx_properties_published_at ON properties(published_at);
CREATE INDEX IF NOT EXISTS idx_properties_deleted_at ON properties(deleted_at);

CREATE INDEX IF NOT EXISTS idx_property_images_property_id ON property_images(property_id);
CREATE INDEX IF NOT EXISTS idx_property_images_is_cover ON property_images(is_cover);

CREATE INDEX IF NOT EXISTS idx_agents_user_id ON agents(user_id);
CREATE INDEX IF NOT EXISTS idx_agents_agency_id ON agents(agency_id);
CREATE INDEX IF NOT EXISTS idx_agents_status ON agents(status);
CREATE INDEX IF NOT EXISTS idx_agents_deleted_at ON agents(deleted_at);

CREATE INDEX IF NOT EXISTS idx_districts_region ON districts(region);
CREATE INDEX IF NOT EXISTS idx_facilities_category ON facilities(category);

-- 插入香港地区数据
INSERT INTO districts (name_zh_hant, name_zh_hans, name_en, region, sort_order) VALUES
-- 港岛
('中西區', '中西区', 'Central and Western', 'HK_ISLAND', 1),
('灣仔區', '湾仔区', 'Wan Chai', 'HK_ISLAND', 2),
('東區', '东区', 'Eastern', 'HK_ISLAND', 3),
('南區', '南区', 'Southern', 'HK_ISLAND', 4),
-- 九龙
('油尖旺區', '油尖旺区', 'Yau Tsim Mong', 'KOWLOON', 5),
('深水埗區', '深水埗区', 'Sham Shui Po', 'KOWLOON', 6),
('九龍城區', '九龙城区', 'Kowloon City', 'KOWLOON', 7),
('黃大仙區', '黄大仙区', 'Wong Tai Sin', 'KOWLOON', 8),
('觀塘區', '观塘区', 'Kwun Tong', 'KOWLOON', 9),
-- 新界
('葵青區', '葵青区', 'Kwai Tsing', 'NEW_TERRITORIES', 10),
('荃灣區', '荃湾区', 'Tsuen Wan', 'NEW_TERRITORIES', 11),
('屯門區', '屯门区', 'Tuen Mun', 'NEW_TERRITORIES', 12),
('元朗區', '元朗区', 'Yuen Long', 'NEW_TERRITORIES', 13),
('北區', '北区', 'North', 'NEW_TERRITORIES', 14),
('大埔區', '大埔区', 'Tai Po', 'NEW_TERRITORIES', 15),
('沙田區', '沙田区', 'Sha Tin', 'NEW_TERRITORIES', 16),
('西貢區', '西贡区', 'Sai Kung', 'NEW_TERRITORIES', 17),
('離島區', '离岛区', 'Islands', 'NEW_TERRITORIES', 18)
ON CONFLICT DO NOTHING;

-- 插入常见设施数据
INSERT INTO facilities (name_zh_hant, name_zh_hans, name_en, icon, category, sort_order) VALUES
-- 基本设施
('冷氣', '空调', 'Air Conditioning', 'ac', 'basic', 1),
('熱水爐', '热水器', 'Water Heater', 'water-heater', 'basic', 2),
('洗衣機', '洗衣机', 'Washing Machine', 'washing-machine', 'basic', 3),
('雪櫃', '冰箱', 'Refrigerator', 'refrigerator', 'basic', 4),
('電視', '电视', 'Television', 'tv', 'basic', 5),
('寬頻上網', '宽带上网', 'Broadband Internet', 'wifi', 'basic', 6),
-- 厨房设施
('煮食爐', '燃气灶', 'Cooking Stove', 'stove', 'kitchen', 7),
('抽油煙機', '抽油烟机', 'Range Hood', 'range-hood', 'kitchen', 8),
('微波爐', '微波炉', 'Microwave', 'microwave', 'kitchen', 9),
-- 安全设施
('保安系統', '安防系统', 'Security System', 'security', 'security', 10),
('煙霧探測器', '烟雾探测器', 'Smoke Detector', 'smoke-detector', 'security', 11),
('閉路電視', '监控摄像头', 'CCTV', 'cctv', 'security', 12),
-- 公共设施
('會所', '会所', 'Clubhouse', 'clubhouse', 'community', 13),
('游泳池', '游泳池', 'Swimming Pool', 'pool', 'community', 14),
('健身室', '健身房', 'Gym', 'gym', 'community', 15),
('兒童遊樂場', '儿童游乐场', 'Playground', 'playground', 'community', 16),
('停車場', '停车场', 'Parking', 'parking', 'community', 17)
ON CONFLICT DO NOTHING;
