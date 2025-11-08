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
- 🔍 **Migration Preview**: Zero-cost error recovery with transaction rollback testing
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
    "github.com/go-xlan/go-migrate/previewmigrate"
    "github.com/go-xlan/go-migrate/newscripts"
    "github.com/spf13/cobra"
    "github.com/yyle88/must"
)

func main() {
    // Define functions for on-demand initialization
    getDB := func() *gorm.DB {
        return setupDatabase()
    }
    getMigration := func(db *gorm.DB) *migrate.Migrate {
        return setupMigration(db)
    }

    var rootCmd = &cobra.Command{Use: "app"}

    // Add migration commands
    rootCmd.AddCommand(cobramigration.NewMigrateCmd(getDB, getMigration))
    rootCmd.AddCommand(previewmigrate.NewPreviewCmd(getDB, getMigration, "./scripts"))
    rootCmd.AddCommand(newscripts.NextScriptCmd(&newscripts.Config{
        GetMigration: getMigration,
        GetDB:        getDB,
        Options:      newscripts.NewOptions("./scripts"),
        Objects:      []any{&User{}, &Product{}},
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
- `preview inc` - Preview next migration without changes

## 📁 Project Structure

```
go-migrate/
├── checkmigration/     # Schema analysis and SQL generation
├── newmigrate/         # Migration instance factory
├── newscripts/         # Script generation and management  
├── cobramigration/     # Cobra CLI integration
├── previewmigrate/     # Migration preview and testing
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

# Preview and execute migrations
make MIGRATE-PREVIEW-INC
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

# Preview and execute migrations
make MIGRATE-PREVIEW-INC
make MIGRATE-ALL
make MIGRATE-INC
```

<!-- TEMPLATE (EN) BEGIN: STANDARD PROJECT FOOTER -->
<!-- VERSION 2025-09-26 07:39:27.188023 +0000 UTC -->

## 📄 License

MIT License. See [LICENSE](LICENSE).

---

## 🤝 Contributing

Contributions are welcome! Report bugs, suggest features, and contribute code:

- 🐛 **Found a mistake?** Open an issue on GitHub with reproduction steps
- 💡 **Have a feature idea?** Create an issue to discuss the suggestion
- 📖 **Documentation confusing?** Report it so we can improve
- 🚀 **Need new features?** Share the use cases to help us understand requirements
- ⚡ **Performance issue?** Help us optimize through reporting slow operations
- 🔧 **Configuration problem?** Ask questions about complex setups
- 📢 **Follow project progress?** Watch the repo to get new releases and features
- 🌟 **Success stories?** Share how this package improved the workflow
- 💬 **Feedback?** We welcome suggestions and comments

---

## 🔧 Development

New code contributions, follow this process:

1. **Fork**: Fork the repo on GitHub (using the webpage UI).
2. **Clone**: Clone the forked project (`git clone https://github.com/yourname/repo-name.git`).
3. **Navigate**: Navigate to the cloned project (`cd repo-name`)
4. **Branch**: Create a feature branch (`git checkout -b feature/xxx`).
5. **Code**: Implement the changes with comprehensive tests
6. **Testing**: (Golang project) Ensure tests pass (`go test ./...`) and follow Go code style conventions
7. **Documentation**: Update documentation to support client-facing changes and use significant commit messages
8. **Stage**: Stage changes (`git add .`)
9. **Commit**: Commit changes (`git commit -m "Add feature xxx"`) ensuring backward compatible code
10. **Push**: Push to the branch (`git push origin feature/xxx`).
11. **PR**: Open a merge request on GitHub (on the GitHub webpage) with detailed description.

Please ensure tests pass and include relevant documentation updates.

---

## 🌟 Support

Welcome to contribute to this project via submitting merge requests and reporting issues.

**Project Support:**

- ⭐ **Give GitHub stars** if this project helps you
- 🤝 **Share with teammates** and (golang) programming friends
- 📝 **Write tech blogs** about development tools and workflows - we provide content writing support
- 🌟 **Join the ecosystem** - committed to supporting open source and the (golang) development scene

**Have Fun Coding with this package!** 🎉🎉🎉

<!-- TEMPLATE (EN) END: STANDARD PROJECT FOOTER -->

---

## GitHub Stars

[![Stargazers](https://starchart.cc/go-xlan/go-migrate.svg?variant=adaptive)](https://starchart.cc/go-xlan/go-migrate)