package errcode

var (

	// 用户相关
	CreateUserFail           = NewError(30010001, "创建用户失败")
	CreateUserFailNameExist  = NewError(30010002, "创建用户失败，用户名已被占用")
	CreateUserFailPhoneExist = NewError(30010003, "创建用户失败，手机号已被占用")
	UpdateUserFail           = NewError(30010004, "更新用户信息失败")
	UpdateUserFailNameExist  = NewError(30010005, "更新用户信息失败，用户名已被占用")
	UpdateUserFailPhoneExist = NewError(30010006, "更新用户信息失败，手机号已被占用")
	DeleteUserFail           = NewError(30010007, "删除用户失败")
	CountUserFail            = NewError(30010008, "统计用户数量失败")
	GetUserListFail          = NewError(30010009, "获取用户列表失败")
	GetUserFail              = NewError(300100010, "获取用户信息失败")
)
