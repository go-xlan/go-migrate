# go-migrate

智能数据库迁移工具包，集成 GORM 模型分析和自动化脚本生成功能。

<!-- TEMPLATE (ZH) BEGIN: LANGUAGE NAVIGATION -->
## 英文文档

[ENGLISH README](README.md)
<!-- TEMPLATE (ZH) END: LANGUAGE NAVIGATION -->

## ✨ 核心特性

- 🔍 **智能结构分析**：自动对比 GORM 模型与实际数据库结构
- 📝 **自动脚本生成**：智能版本管理的迁移脚本创建功能
- 🔄 **灵活迁移策略**：支持基于文件、嵌入式和数据库驱动的方式
- 🎯 **全面 CLI 支持**：用户友好的 Cobra 命令覆盖所有迁移操作
- 🛡️ **安全操作模式**：DryRun 模式和交互式确认保障迁移安全
- 🔍 **迁移预览功能**：事务回滚测试实现零成本错误恢复
- 🔗 **多数据库兼容**：通过 golang-migrate 支持 MySQL、PostgreSQL、SQLite

## 📦 安装

```bash
go get github.com/go-xlan/go-migrate
```

### 前置条件
- Go 1.22.8 或更高版本
- 目标数据库的相应驱动
- GORM v2 用于模型定义

## 🚀 快速开始

### 基础用法

```go
package main

import (
    "github.com/go-xlan/go-migrate/checkmigration"
    "github.com/go-xlan/go-migrate/newmigrate"
    "github.com/yyle88/must"
    "gorm.io/gorm"
)

func main() {
    // 初始化 GORM 数据库连接
    db := setupDatabase() // 你的数据库设置
    
    // 检查需要执行的迁移
    migrateSQLs := checkmigration.CheckMigrate(db, []any{&User{}, &Product{}})
    
    // 创建迁移实例
    migration := must.Nice(newmigrate.NewWithScriptsAndDatabase(&newmigrate.ScriptsAndDatabaseParam{
        ScriptsInRoot:    "./migrations",
        DatabaseName:     "mysql",
        DatabaseInstance: databaseDriver, // 你的数据库驱动实例
    }))
    
    // 执行迁移
    must.Done(migration.Up())
}
```

### CLI 集成

```go
package main

import (
    "github.com/go-xlan/go-migrate/cobramigration"
    "github.com/go-xlan/go-migrate/newscripts"
    "github.com/spf13/cobra"
    "github.com/yyle88/must"
)

func main() {
    // 定义工厂函数用于延迟初始化
    getDB := func() *gorm.DB {
        return setupDatabase()
    }
    getMigration := func(db *gorm.DB) *migrate.Migrate {
        return setupMigration(db)
    }

    var rootCmd = &cobra.Command{Use: "app"}

    // 添加迁移命令
    rootCmd.AddCommand(cobramigration.NewMigrateCmd(getDB, getMigration))
    rootCmd.AddCommand(newscripts.NextScriptCmd(&newscripts.Config{
        GetMigration: getMigration,
        GetDB:        getDB,
        Options:      newscripts.NewOptions("./scripts"),
        Objects:      []any{&User{}, &Product{}},
    }))

    must.Done(rootCmd.Execute())
}
```

## 📋 核心 API 参考

### 迁移分析
- `checkmigration.CheckMigrate(db, models)` - 对比结构并返回所需 SQL
- `checkmigration.GetMigrateOps(db, models)` - 获取详细迁移操作信息

### 迁移创建
- `newmigrate.NewWithScriptsAndDBSource[T](param)` - 使用连接字符串创建
- `newmigrate.NewWithScriptsAndDatabase(param)` - 使用驱动实例创建
- `newmigrate.NewWithEmbedFsAndDatabase(param)` - 使用嵌入文件创建

### 脚本管理
- `newscripts.GetNextScriptInfo(migration, options, naming)` - 分析下一脚本需求
- `newscripts.NextScriptCmd(config)` - 脚本生成的 CLI 命令

### CLI 命令
- `migrate` - 显示当前迁移状态
- `migrate all` - 执行所有待处理迁移
- `migrate inc` - 运行下一个迁移步骤
- `migrate dec` - 回滚一个迁移步骤

## 📁 项目结构

```
go-migrate/
├── checkmigration/     # 结构分析和 SQL 生成
├── newmigrate/         # 迁移实例工厂
├── newscripts/         # 脚本生成和管理
├── cobramigration/     # Cobra CLI 集成
└── internal/           # 演示、示例和工具
    ├── demos/          # 完整演示应用
    ├── examples/       # 使用示例
    └── sketches/       # 开发草图
```

## 🔧 配置示例

### 数据库设置

```go
// MySQL 配置
migration := rese.V1(newmigrate.NewWithScriptsAndDatabase(&newmigrate.ScriptsAndDatabaseParam{
    ScriptsInRoot:    "./migrations",
    DatabaseName:     "mysql",
    DatabaseInstance: mysqlDriver,
}))

// PostgreSQL 配置
migration := rese.V1(newmigrate.NewWithScriptsAndDBSource[*postgres.Postgres](&newmigrate.ScriptsAndDBSourceParam{
    ScriptsInRoot: "./migrations",
    ConnectSource: "postgres://user:pass@localhost/db?sslmode=disable",
}))

// SQLite 配置
migration := rese.V1(newmigrate.NewWithScriptsAndDBSource[*sqlite3.Sqlite](&newmigrate.ScriptsAndDBSourceParam{
    ScriptsInRoot: "./migrations",
    ConnectSource: "sqlite3://./database.db",
}))
```

### 嵌入式迁移

```go
//go:embed migrations
var migrationsFS embed.FS

migration := rese.V1(newmigrate.NewWithEmbedFsAndDatabase(&newmigrate.EmbedFsAndDatabaseParam{
    MigrationsFS:     &migrationsFS,
    EmbedDirName:     "migrations",
    DatabaseName:     "mysql",
    DatabaseInstance: driver,
}))
```

## 🎯 高级特性

### 自定义脚本命名

```go
naming := &newscripts.ScriptNaming{
    NewScriptPrefix: func(version uint) string {
        return fmt.Sprintf("%d_%s", version, description)
    },
}
```

### 迁移选项配置

```go
options := newscripts.NewOptions("./scripts").
    WithDryRun(true).
    WithSurveyWritten(true)
```

## 📖 示例

查看 `internal/demos/` DIR 中的完整工作示例：

- **demo1x/**：MySQL 集成和 Makefile 命令
- **demo2x/**：PostgreSQL 集成和 Makefile 命令
- **examples/**：聚焦功能演示
- **sketches/**：开发原型

### 演示命令

**Demo1x - MySQL 集成示例：**
```bash
# 导航到 demo1x DIR
cd internal/demos/demo1x

# 生成迁移脚本
make CREATE-SCRIPT-CREATE-TABLE
make CREATE-SCRIPT-ALTER-SCHEMA

# 执行迁移
make MIGRATE-ALL
make MIGRATE-INC
```

**Demo2x - PostgreSQL 集成示例：**
```bash
# 导航到 demo2x DIR
cd internal/demos/demo2x

# 生成迁移脚本
make CREATE-SCRIPT-CREATE-TABLE
make CREATE-SCRIPT-ALTER-SCHEMA

# 执行迁移
make MIGRATE-ALL
make MIGRATE-INC
```

<!-- TEMPLATE (ZH) BEGIN: STANDARD PROJECT FOOTER -->
<!-- VERSION 2025-09-26 07:39:27.188023 +0000 UTC -->

## 📄 许可证类型

MIT 许可证。详见 [LICENSE](LICENSE)。

---

## 🤝 项目贡献

非常欢迎贡献代码！报告 BUG、建议功能、贡献代码：

- 🐛 **发现问题？** 在 GitHub 上提交问题并附上重现步骤
- 💡 **功能建议？** 创建 issue 讨论您的想法
- 📖 **文档疑惑？** 报告问题，帮助我们改进文档
- 🚀 **需要功能？** 分享使用场景，帮助理解需求
- ⚡ **性能瓶颈？** 报告慢操作，帮助我们优化性能
- 🔧 **配置困扰？** 询问复杂设置的相关问题
- 📢 **关注进展？** 关注仓库以获取新版本和功能
- 🌟 **成功案例？** 分享这个包如何改善工作流程
- 💬 **反馈意见？** 欢迎提出建议和意见

---

## 🔧 代码贡献

新代码贡献，请遵循此流程：

1. **Fork**：在 GitHub 上 Fork 仓库（使用网页界面）
2. **克隆**：克隆 Fork 的项目（`git clone https://github.com/yourname/repo-name.git`）
3. **导航**：进入克隆的项目（`cd repo-name`）
4. **分支**：创建功能分支（`git checkout -b feature/xxx`）
5. **编码**：实现您的更改并编写全面的测试
6. **测试**：（Golang 项目）确保测试通过（`go test ./...`）并遵循 Go 代码风格约定
7. **文档**：为面向用户的更改更新文档，并使用有意义的提交消息
8. **暂存**：暂存更改（`git add .`）
9. **提交**：提交更改（`git commit -m "Add feature xxx"`）确保向后兼容的代码
10. **推送**：推送到分支（`git push origin feature/xxx`）
11. **PR**：在 GitHub 上打开 Merge Request（在 GitHub 网页上）并提供详细描述

请确保测试通过并包含相关的文档更新。

---

## 🌟 项目支持

非常欢迎通过提交 Merge Request 和报告问题来为此项目做出贡献。

**项目支持：**

- ⭐ **给予星标**如果项目对您有帮助
- 🤝 **分享项目**给团队成员和（golang）编程朋友
- 📝 **撰写博客**关于开发工具和工作流程 - 我们提供写作支持
- 🌟 **加入生态** - 致力于支持开源和（golang）开发场景

**祝你用这个包编程愉快！** 🎉🎉🎉

<!-- TEMPLATE (ZH) END: STANDARD PROJECT FOOTER -->

---

## GitHub 标星点赞

[![Stargazers](https://starchart.cc/go-xlan/go-migrate.svg?variant=adaptive)](https://starchart.cc/go-xlan/go-migrate)