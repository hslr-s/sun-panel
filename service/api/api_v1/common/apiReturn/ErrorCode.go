package apiReturn

var ErrorCodeMap = map[int]string{
	// -1:操作失败(前端会自动弹窗)
	// 100: "operation failed",

	1000: "Not logged in yet",                   // 还未登录
	1003: "Incorrect username or password",      // 用户名或密码错误
	1004: "Account disabled or not activated",   // 账号已停用或未激活
	1005: "No current permission for operation", // 当前无权限操作
	1006: "Account does not exist",              // 账号不存在
	1007: "Old password error",                  // 旧密码不正确

	// 数据类
	1200: "Database error",           // 数据库错误
	1201: "Please keep at least one", // 请至少保留一个
	1202: "No data record found",     // 未找到数据记录

	1300: "Upload failed",           // 上传失败
	1301: "Unsupported file format", // 不被支持的格式文件

	1400: "Parameter format error", // 参数格式错误

}
