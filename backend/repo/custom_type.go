package repo

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

func createTypeWithRaw(db *gorm.DB, types, def string) error {
	query := fmt.Sprintf(`DO $$ BEGIN IF to_regtype('%v') IS NULL THEN CREATE TYPE  %v AS %v ;END IF;END $$;`, types, types, def)
	return db.Exec(query).Error
}
func CreateCustomTypes(db *gorm.DB) {
	err := createTypeWithRaw(db, "augmont_kyc_status", `ENUM ('approved', 'pending', 'rejected')`)
	if err != nil {
		log.Fatalln(err)
	}
}
