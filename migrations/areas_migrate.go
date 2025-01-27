package main

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/arslion-7/api-construction-share/initializers"
	"github.com/arslion-7/api-construction-share/models"
	_ "github.com/go-sql-driver/mysql"
)

var PostgresDB *gorm.DB

func main() {
	// Load environment variables
	initializers.LoadEnvVars()
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	fmt.Println("dbHost", dbHost)

	// Connect to MySQL
	mysqlDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", "root", "010203", "192.168.0.242", "3306", "adresturkmen")
	mysqlDB, err := sql.Open("mysql", mysqlDSN)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	defer mysqlDB.Close()

	// Connect to PostgreSQL
	postgresDSN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	PostgresDB, err = gorm.Open(postgres.Open(postgresDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	// Query MySQL for data
	rows, err := mysqlDB.Query("SELECT kod, etr_shah FROM hemme_etr_shah")
	if err != nil {
		log.Fatalf("Failed to query MySQL: %v", err)
	}
	defer rows.Close()

	// Process rows and migrate to PostgreSQL
	for rows.Next() {
		var kod float64
		var etrShah string

		if err := rows.Scan(&kod, &etrShah); err != nil {
			log.Printf("Row scan failed: %v", err)
			continue
		}

		var parentID *uint
		kodInt := uint(kod)
		if math.Mod(kod, 1_000_000_000) == 0 || kod == 1_000_000_001 {
			// Root level or Ashgabat
			parentID = nil
		} else if int(kod)%1_000_000 == 0 {
			// Welayat level
			parentID = getWelayat(uint(kod))
		} else {
			// Etrap level
			parentID = getEtrap(uint(kod))
			if parentID == nil {
				parentID = getWelayat(uint(kod))
			}
		}

		area := models.Area{
			Code:     kodInt,
			Name:     etrShah,
			ParentID: parentID,
		}

		// Upsert into PostgreSQL
		if err := PostgresDB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "code"}},
			DoUpdates: clause.AssignmentColumns([]string{"name", "parent_id"}),
		}).Create(&area).Error; err != nil {
			log.Printf("Failed to insert or update area: %v", err)
		} else {
			log.Printf("Processed area: %v", area.Name)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Row iteration error: %v", err)
	}
}

// getWelayat finds the parent ID for a welayat based on the code.
func getWelayat(kod uint) *uint {
	welayatCode := kod / 1_000_000_000 * 1_000_000_000
	if welayatCode == 1_000_000_000 {
		welayatCode = 1_000_000_001 // Ashgabat specific
	}

	var welayat models.Area
	if err := PostgresDB.First(&welayat, "code = ?", welayatCode).Error; err != nil {
		return nil
	}
	return &welayat.Code
}

// getEtrap finds the parent ID for an etrap based on the code.
func getEtrap(kod uint) *uint {
	etrapCode := kod / 1_000_000 * 1_000_000

	var etrap models.Area
	if err := PostgresDB.First(&etrap, "code = ?", etrapCode).Error; err != nil {
		return nil
	}
	return &etrap.Code
}
