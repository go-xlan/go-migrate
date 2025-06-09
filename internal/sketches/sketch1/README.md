## 🧾 Module Overview (English)

This module demonstrates how to use the [`golang-migrate`](https://github.com/golang-migrate/migrate) library's `file` source to load SQL migration scripts directly from a local `scripts/` folder, suitable for local development and testing.

### Key Features:

* Uses the `file://` scheme to read migration scripts from the local filesystem;
* Supports iterating through all migration versions and reading both `up` and `down` SQL files;
* Useful for debugging or verifying the content of migration scripts during development;

The `TestFileOpen` function demonstrates how to use the `source.Driver` interface to traverse migration versions and load their SQL contents in sequence.

---

## 🧾 模块简介（中文）

演示如何使用 [`golang-migrate`](https://github.com/golang-migrate/migrate) 提供的 `file` 驱动，从本地文件系统中的 `scripts/` 文件夹读取 SQL 迁移脚本，适用于本地开发和测试场景。

### 核心特点：

* 使用 `file://` 路径从本地加载迁移脚本；
* 支持遍历所有迁移版本，并读取每个版本对应的 `up` 和 `down` SQL 文件；
* 便于在开发过程中调试或验证迁移文件内容；

测试函数 `TestFileOpen` 展示了如何通过 `source.Driver` 接口，顺序读取迁移版本及其 SQL 内容的完整流程。
