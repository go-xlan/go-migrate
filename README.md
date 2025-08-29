# go-migrate

Intelligent database migration toolkit with GORM model integration and automated script generation.

<!-- TEMPLATE (EN) BEGIN: LANGUAGE NAVIGATION -->
## CHINESE README

[中文说明](README.zh.md)
<!-- TEMPLATE (EN) END: LANGUAGE NAVIGATION -->

## ✨ Features

- 🔍 **Smart Schema Analysis**: Auto-compare GORM models with actual database schemas
- 📝 **Automated Script Generation**: Create migration scripts with intelligent version management  
- 🔄 **Flexible Migration Strategies**: Support file-based, embedded, and database-driven approaches
- 🎯 **Comprehensive CLI**: User-friendly Cobra commands for all migration operations
- 🛡️ **Safe Operations**: DryRun mode and interactive confirmation for secure migrations
- 🔗 **Multi-Database Support**: Works with MySQL, PostgreSQL, SQLite through golang-migrate

## 📦 Installation

```bash
go get github.com/go-xlan/go-migrate
```

### Prerequisites
- Go 1.22.8 or later
- Database driver for your target database
- GORM v2 for model definitions

## 🚀 Quick Start

### Basic Usage

```go
package main

import (
    "github.com/go-xlan/go-migrate/checkmigration"
    "github.com/go-xlan/go-migrate/newmigrate"
    "github.com/yyle88/must"
    "gorm.io/gorm"
)

func main() {
    // Initialize GORM database connection
    db := setupDatabase() // Your database setup
    
    // Check what migrations are needed
    migrateSQLs := checkmigration.CheckMigrate(db, []any{&User{}, &Product{}})
    
    // Create migration instance
    migration := must.Nice(newmigrate.NewWithScriptsAndDatabase(&newmigrate.ScriptsAndDatabaseParam{
        ScriptsInRoot:    "./migrations",
        DatabaseName:     "mysql",
        DatabaseInstance: databaseDriver, // Your database driver instance
    }))
    
    // Execute migrations
    must.Done(migration.Up())
}
```

### CLI Integration

```go
package main

import (
    "github.com/go-xlan/go-migrate/cobramigration"
    "github.com/go-xlan/go-migrate/newscripts"
    "github.com/spf13/cobra"
    "github.com/yyle88/must"
)

func main() {
    // Setup database and migration instance (from previous example)
    db := setupDatabase()
    migration := setupMigration()
    
    var rootCmd = &cobra.Command{Use: "app"}
    
    // Add migration commands
    rootCmd.AddCommand(cobramigration.NewMigrateCmd(migration))
    rootCmd.AddCommand(newscripts.NextScriptCmd(&newscripts.Config{
        Migration: migration,
        Options:   newscripts.NewOptions("./scripts"),
        DB:        db,
        Objects:   []any{&User{}, &Product{}},
    }))
    
    must.Done(rootCmd.Execute())
}
```

## 📋 Core API Reference

### Migration Analysis
- `checkmigration.CheckMigrate(db, models)` - Compare schemas and return needed SQL
- `checkmigration.GetMigrateOps(db, models)` - Get detailed migration operations

### Migration Creation
- `newmigrate.NewWithScriptsAndDBSource[T](param)` - Create with connection string
- `newmigrate.NewWithScriptsAndDatabase(param)` - Create with driver instance
- `newmigrate.NewWithEmbedFsAndDatabase(param)` - Create with embedded files

### Script Management
- `newscripts.GetNextScriptInfo(migration, options, naming)` - Analyze next script requirements
- `newscripts.NextScriptCmd(config)` - CLI command for script generation

### CLI Commands
- `migrate` - Display current migration status
- `migrate all` - Execute all pending migrations
- `migrate inc` - Run next migration step
- `migrate dec` - Rollback one migration step

## 📁 Project Structure

```
go-migrate/
├── checkmigration/     # Schema analysis and SQL generation
├── newmigrate/         # Migration instance factory
├── newscripts/         # Script generation and management  
├── cobramigration/     # Cobra CLI integration
└── internal/           # Demos, examples, and utilities
    ├── demos/          # Complete demo applications
    ├── examples/       # Usage examples
    └── sketches/       # Development sketches
```

## 🔧 Configuration Examples

### Database Setup

```go
// MySQL configuration
migration := rese.V1(newmigrate.NewWithScriptsAndDatabase(&newmigrate.ScriptsAndDatabaseParam{
    ScriptsInRoot:    "./migrations",
    DatabaseName:     "mysql",
    DatabaseInstance: mysqlDriver,
}))

// PostgreSQL configuration  
migration := rese.V1(newmigrate.NewWithScriptsAndDBSource[*postgres.Postgres](&newmigrate.ScriptsAndDBSourceParam{
    ScriptsInRoot: "./migrations",
    ConnectSource: "postgres://user:pass@localhost/db?sslmode=disable",
}))

// SQLite configuration
migration := rese.V1(newmigrate.NewWithScriptsAndDBSource[*sqlite3.Sqlite](&newmigrate.ScriptsAndDBSourceParam{
    ScriptsInRoot: "./migrations",
    ConnectSource: "sqlite3://./database.db",
}))
```

### Embedded Migrations

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

## 🎯 Advanced Features

### Custom Script Naming

```go
naming := &newscripts.ScriptNaming{
    NewScriptPrefix: func(version uint) string {
        return fmt.Sprintf("%d_%s", version, description)
    },
}
```

### Migration Options

```go
options := newscripts.NewOptions("./scripts").
    WithDryRun(true).
    WithSurveyWritten(true)
```

## 📖 Examples

Check the `internal/demos/` DIR for complete working examples:

- **demo1x/**: MySQL integration with Makefile commands
- **demo2x/**: PostgreSQL integration with Makefile commands
- **examples/**: Focused feature demonstrations
- **sketches/**: Development prototypes

### Demo Commands

**Demo1x - MySQL Integration:**
```bash
# Navigate to demo1x DIR
cd internal/demos/demo1x

# Generate migration scripts
make CREATE-SCRIPT-CREATE-TABLE
make CREATE-SCRIPT-ALTER-SCHEMA

# Execute migrations
make MIGRATE-ALL
make MIGRATE-INC
```

**Demo2x - PostgreSQL Integration:**
```bash
# Navigate to demo2x DIR
cd internal/demos/demo2x

# Generate migration scripts
make CREATE-SCRIPT-CREATE-TABLE
make CREATE-SCRIPT-ALTER-SCHEMA

# Execute migrations
make MIGRATE-ALL
make MIGRATE-INC
```

<!-- TEMPLATE (EN) BEGIN: STANDARD PROJECT FOOTER -->
<!-- VERSION 2025-08-29 08:33:43.829511 +0000 UTC -->

## 📄 License

MIT License. See [LICENSE](LICENSE).

---

## 🤝 Contributing

Contributions are welcome! Report bugs, suggest features, and contribute code:

- 🐛 **Found a bug?** Open an issue on GitHub with reproduction steps
- 💡 **Have a feature idea?** Create an issue to discuss the suggestion
- 📖 **Documentation confusing?** Report it so we can improve
- 🚀 **Need new features?** Share your use cases to help us understand requirements
- ⚡ **Performance issue?** Help us optimize by reporting slow operations
- 🔧 **Configuration problem?** Ask questions about complex setups
- 📢 **Follow project progress?** Watch the repo for new releases and features
- 🌟 **Success stories?** Share how this package improved your workflow
- 💬 **General feedback?** All suggestions and comments are welcome

---

## 🔧 Development

New code contributions, follow this process:

1. **Fork**: Fork the repo on GitHub (using the webpage interface).
2. **Clone**: Clone the forked project (`git clone https://github.com/yourname/go-migrate.git`).
3. **Navigate**: Navigate to the cloned project (`cd go-migrate`)
4. **Branch**: Create a feature branch (`git checkout -b feature/xxx`).
5. **Code**: Implement your changes with comprehensive tests
6. **Testing**: (Golang project) Ensure tests pass (`go test ./...`) and follow Go code style conventions
7. **Documentation**: Update documentation for user-facing changes and use meaningful commit messages
8. **Stage**: Stage changes (`git add .`)
9. **Commit**: Commit changes (`git commit -m "Add feature xxx"`) ensuring backward compatible code
10. **Push**: Push to the branch (`git push origin feature/xxx`).
11. **PR**: Open a pull request on GitHub (on the GitHub webpage) with detailed description.

Please ensure tests pass and include relevant documentation updates.

---

## 🌟 Support

Welcome to contribute to this project by submitting pull requests and reporting issues.

**Project Support:**

- ⭐ **Give GitHub stars** if this project helps you
- 🤝 **Share with teammates** and (golang) programming friends
- 📝 **Write tech blogs** about development tools and workflows - we provide content writing support
- 🌟 **Join the ecosystem** - committed to supporting open source and the (golang) development scene

**Happy Coding with this package!** 🎉

<!-- TEMPLATE (EN) END: STANDARD PROJECT FOOTER -->

---

## GitHub Stars

[![Stargazers](https://starchart.cc/go-xlan/go-migrate.svg?variant=adaptive)](https://starchart.cc/go-xlan/go-migrate)