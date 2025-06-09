## 🧾 Module Overview (English)

This module demonstrates how to embed SQL migration scripts into a Go program using `embed.FS`, and how to read and apply them using the `iofs` driver from `golang-migrate`.

Benefits of this approach:

* Migration scripts are versioned with the code;
* Easier deployment without managing external `.sql` files;
* Useful for automated testing or CI/CD.

Two embedding strategies are shown:

* Embedding the full `scripts/` directory;
* Embedding only selected `.sql` files.

Each has its own use case — full history vs. lightweight testing.

---

## 🧾 模块简介（中文）

这个模块展示了如何将数据库迁移脚本通过 `embed.FS` 嵌入到 Go 程序中，并结合 `golang-migrate` 的 `iofs` 驱动实现读取和执行。

这样做的优点是：

* 脚本和程序版本一致，不易丢失或出错；
* 部署更简单，无需依赖外部 SQL 文件；
* 更适合在测试、CI/CD 中使用。

模块中分别演示了：

* 嵌入整个 `scripts/` 文件夹；
* 只嵌入部分特定的 `.sql` 文件。

这两种方式适用于不同场景，前者适合完整管理迁移记录，后者适合简化测试或减小二进制体积。
