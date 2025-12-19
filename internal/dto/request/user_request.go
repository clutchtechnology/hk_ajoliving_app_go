package request

// UpdateUserRequest 更新用户信息请求
type UpdateUserRequest struct {
	FullName string `json:"full_name" binding:"omitempty,max=100"`
	Phone    string `json:"phone" binding:"omitempty,min=8,max=20"`
	Avatar   string `json:"avatar" binding:"omitempty,url,max=500"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=6,max=100"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=100"`
}

// UpdateSettingsRequest 更新设置请求
type UpdateSettingsRequest struct {
	Language            string `json:"language" binding:"omitempty,oneof=zh_HK zh_CN en"`
	NotificationEnabled bool   `json:"notification_enabled"`
	EmailNotification   bool   `json:"email_notification"`
	SmsNotification     bool   `json:"sms_notification"`
}
