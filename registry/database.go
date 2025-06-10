package registry

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/gorm"
	"template.go/packages/common"
)

var DBRegistry = &DatabaseRegistry{
	enums: []Enum{
		// enums type
	},
	tables: []string{
		// tables for backup
	},
	models: []any{
		// models for auto-migrate
	},
	extensions: []string{
		// database extensions
	},
	factories: []func(*gorm.DB) error{
		// register factory for seeding, example:
		// func(db *gorm.DB) error {
		// 	return factory.NewUserFactory(db).CreateBatch(1)
		// },
	},
}

type Enum struct {
	Name   string
	Values []string
}

func (e *Enum) CreateQuery() string {
	return fmt.Sprintf("CREATE TYPE IF NOT EXISTS %s AS ENUM (%s)", e.Name, strings.Join(e.Values, ", "))
}

func (e *Enum) DropQuery() string {
	return fmt.Sprintf("DROP TYPE IF EXISTS %s", e.Name)
}

type DatabaseRegistry struct {
	enums      []Enum
	models     []any
	extensions []string
	tables     []string
	factories  []func(*gorm.DB) error
}

func (r *DatabaseRegistry) GetEnums() []Enum {
	return r.enums
}

func (r *DatabaseRegistry) GetTables() []string {
	return r.tables
}

func (r *DatabaseRegistry) GetModels() []any {
	return r.models
}

func (r *DatabaseRegistry) GetExtensions() []string {
	return r.extensions
}

func (r *DatabaseRegistry) GetFactories() []func(*gorm.DB) error {
	return r.factories
}

func (r *DatabaseRegistry) Migrate(db *gorm.DB, fresh bool) error {
	slog.Info("Starting database migration...")

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if fresh {
		if err := dropAll(tx, r); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := createExtensions(tx, r.extensions); err != nil {
		tx.Rollback()
		return err
	}

	if err := createEnums(tx, r.enums); err != nil {
		tx.Rollback()
		return err
	}

	if !fresh {
		for _, enum := range r.enums {
			ensureEnumValues(tx, enum.Name, enum.Values)
		}
	}

	if err := autoMigrateModels(tx, r.models); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	slog.Info("✅ Database migration completed successfully")
	return nil
}

func (r *DatabaseRegistry) Backup(db *gorm.DB, output string) error {
	slog.Info("Starting database backup...")

	if err := os.MkdirAll(output, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	timestamp := time.Now().Format("200601021504") // YYYYMMDDHHMM format

	for _, table := range r.tables {
		var results []map[string]any
		if err := db.Table(table).Find(&results).Error; err != nil {
			return fmt.Errorf("failed to query table %s: %v", table, err)
		}

		for i, row := range results {
			camelCaseRow := make(map[string]any)
			for key, value := range row {
				camelKey := common.ToCamelCase(key)
				camelCaseRow[camelKey] = value
			}
			results[i] = camelCaseRow
		}

		filename := fmt.Sprintf("backup_%s_%s.json", table, timestamp)
		filePath := filepath.Join(output, filename)

		file, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %v", filePath, err)
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(results); err != nil {
			return fmt.Errorf("failed to encode JSON for table %s: %v", table, err)
		}
	}

	slog.Info("✅ Database backup completed successfully")
	return nil
}

// helper

func dropAll(tx *gorm.DB, registry *DatabaseRegistry) error {
	tables, err := tx.Migrator().GetTables()
	if err != nil {
		return err
	}

	if len(tables) > 0 {
		if err := tx.Exec(fmt.Sprintf("TRUNCATE %s CASCADE", strings.Join(tables, ", "))).Error; err != nil {
			return err
		}

		for i, t := range tables {
			tables[i] = fmt.Sprintf(`DROP TABLE IF EXISTS "%s" CASCADE`, t)
		}
		if err := tx.Exec(strings.Join(tables, "; ")).Error; err != nil {
			return err
		}
	}

	enumNames := make([]string, len(registry.enums))
	for i, e := range registry.enums {
		enumNames[i] = e.Name
	}
	if err := dropEnums(tx, enumNames); err != nil {
		return err
	}

	return nil
}

func createExtensions(tx *gorm.DB, exts []string) error {
	for _, q := range exts {
		if err := tx.Exec(q).Error; err != nil {
			return err
		}
	}
	return nil
}

func createEnums(tx *gorm.DB, enums []Enum) error {
	for _, e := range enums {
		if err := tx.Exec(e.CreateQuery()).Error; err != nil {
			return err
		}
	}
	return nil
}

func dropEnums(tx *gorm.DB, names []string) error {
	for _, name := range names {
		if err := tx.Exec("DROP TYPE IF EXISTS " + name).Error; err != nil {
			return err
		}
	}
	return nil
}

func autoMigrateModels(tx *gorm.DB, models []any) error {
	return tx.AutoMigrate(models...)
}

func ensureEnumValues(db *gorm.DB, enumName string, values []string) {
	for _, v := range values {
		safeVal := strings.ReplaceAll(v, `'`, ``)

		var exists bool
		checkSQL := `
			SELECT EXISTS (
				SELECT 1 FROM pg_enum 
				WHERE enumlabel = ? 
				AND enumtypid = ?::regtype
			);`
		if err := db.Raw(checkSQL, safeVal, enumName).Scan(&exists).Error; err != nil {
			slog.Error("Error checking value", slog.Any(enumName, err))
			continue
		}

		if !exists {
			stmt := fmt.Sprintf("ALTER TYPE %s ADD VALUE %s;", enumName, safeVal)
			if err := db.Exec(stmt).Error; err != nil {
				slog.Error("Add enum value", slog.Any(enumName, err))
			}
		}
	}
}
