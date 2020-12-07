package errcode

var (

	// 用户相关
	CreateUserFail  = NewError(30010001, "创建用户失败")
	UpdateUserFail  = NewError(30010002, "更新用户信息失败")
	DeleteUserFail  = NewError(30001003, "删除用户失败")
	CountUserFail   = NewError(30001004, "统计用户数量失败")
	GetUserListFail = NewError(30001005, "获取用户列表失败")
)
