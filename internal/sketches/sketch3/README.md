## 🧾 Module Overview (English)

This module demonstrates how to work with `embed.FS` in Go to load and parse migration scripts using `golang-migrate`'s `source.DefaultParse` and build a migration list manually with `source.NewMigrations()`.

Key features:

* Reads embedded migration scripts with `fs.ReadDir`;
* Parses file names to extract migration metadata (version, direction);
* Appends parsed migrations into a `Migrations` structure;
* Exercises `.First()`, `.Up()`, and `.Down()` methods to simulate migration logic.

Can help you:

* Writing custom migration logic;
* Testing file name parsing without external dependencies;
* Understanding how `golang-migrate` handles internal migration registration.

---

## 🧾 模块简介（中文）

这个模块展示了如何结合 Go 的 `embed.FS` 与 `golang-migrate` 提供的 `source.DefaultParse` 方法，解析迁移脚本文件名，构建 `Migrations` 列表，并通过代码逻辑模拟查找和访问各版本的迁移。

主要内容包括：

* 使用 `fs.ReadDir` 遍历嵌入的迁移脚本文件；
* 使用 `source.DefaultParse` 解析迁移文件名（例如提取版本号、方向等）；
* 将解析结果逐一添加进 `source.NewMigrations()` 所管理的迁移对象列表；
* 验证 `.First()`、`.Up()` 和 `.Down()` 等方法能正确返回期望的迁移。

适合用于开发自定义迁移逻辑、测试迁移元数据解析、或在不依赖外部 driver 的情况下处理嵌入式 SQL 文件的场景。
