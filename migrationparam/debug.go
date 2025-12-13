package migrationparam

var debugModeOpen = false

// SetDebugMode sets the project-wide debug mode for migration operations
// SetDebugMode 设置项目级别的迁移操作调试模式
func SetDebugMode(enable bool) {
	debugModeOpen = enable
}

// GetDebugMode returns the project-wide debug mode for migration operations
// GetDebugMode 获取项目级别的迁移操作调试模式
func GetDebugMode() bool {
	return debugModeOpen
}
