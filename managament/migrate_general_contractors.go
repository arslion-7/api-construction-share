// GORM equivalent for a migration script from Django to Go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/arslion-7/api-construction-share/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var PostgresDB *gorm.DB

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	// Load environment variables
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Connect to MySQL
	mysqlDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", "root", "010203", "192.168.0.242", "3306", "paylynew")
	mysqlDB, err := sql.Open("mysql", mysqlDSN)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	defer mysqlDB.Close()

	// Connect to PostgreSQL
	postgresDSN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Ashgabat",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	PostgresDB, err = gorm.Open(postgres.Open(postgresDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	// Ensure PostgreSQL schema is created
	if err := PostgresDB.AutoMigrate(&models.GeneralContractor{}); err != nil {
		log.Fatalf("Failed to migrate PostgreSQL schema: %v", err)
	}

	// Query data from MySQL
	rows, err := mysqlDB.Query(`SELECT 
		t_b, 
		telekeci_rayat, guramaAdy, wezipesi, 
		familiya, ady, otchestvo, 
		shahadatBelgi, shahadatSene, ygtyyarNumber, ygtyyarBeginDate, ygtyyarEndDate 
		FROM potr_bash_compl`)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	defer rows.Close()

	// Process rows and insert into PostgreSQL
	for rows.Next() {
		var (
			t_b                 *int
			telekeci_rayat      *string
			guramaAdy           *string
			wezipesi            *string
			familiya            *string
			ady                 *string
			otchestvo           *string
			shahadatBelgi       sql.NullInt64
			shahadatSeneRaw     []uint8
			ygtyyarNumber       *string
			ygtyyarBeginDateRaw []uint8
			ygtyyarEndDateRaw   []uint8
		)

		if err := rows.Scan(&t_b, &telekeci_rayat, &guramaAdy, &wezipesi, &familiya, &ady, &otchestvo, &shahadatBelgi, &shahadatSeneRaw, &ygtyyarNumber, &ygtyyarBeginDateRaw, &ygtyyarEndDateRaw); err != nil {
			log.Printf("Row scan failed: %v", err)
			continue
		}

		// Parse date fields
		shahadatSene := parseDate(shahadatSeneRaw)
		ygtyyarBeginDate := parseDate(ygtyyarBeginDateRaw)
		ygtyyarEndDate := parseDate(ygtyyarEndDateRaw)

		// Parse certNumber
		var certNumber *int
		if shahadatBelgi.Valid {
			n := int(shahadatBelgi.Int64)
			certNumber = &n
		}

		// Create contractor
		headFullName := fmt.Sprintf("%s %s %s", getNonNilString(familiya), getNonNilString(ady), getNonNilString(otchestvo))
		contractor := models.GeneralContractor{
			Org: models.Org{
				TB:      t_b,
				OrgName: guramaAdy,
				OrgType: telekeci_rayat,
				// OrgAdditionalInfo: ,
				HeadPosition: wezipesi,
				HeadFullName: &headFullName,
			},
			CertNumber:          certNumber,
			CertDate:            shahadatSene,
			ResolutionCode:      ygtyyarNumber,
			ResolutionBeginDate: ygtyyarBeginDate,
			ResolutionEndDate:   ygtyyarEndDate,
		}

		// Insert into PostgreSQL
		if err := PostgresDB.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "t_b"}}, // Conflict target
			DoUpdates: clause.AssignmentColumns([]string{
				"org_type", "org_name", "head_position", "head_full_name",
				"org_additional_info", "updated_at", "cert_number", "cert_date",
				"resolution_code", "resolution_begin_date", "resolution_end_date",
			}),
		}).Create(&contractor).Error; err != nil {
			log.Printf("Failed to insert or update contractor: %v", err)
		} else {
			log.Printf("Successfully processed contractor: %v", contractor.OrgName)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Row iteration error: %v", err)
	}
}

// Utility function to parse date fields
func parseDate(raw []uint8) *time.Time {
	if len(raw) == 0 {
		return nil
	}
	parsedTime, err := time.Parse("2006-01-02", string(raw))
	if err != nil {
		log.Printf("Failed to parse date: %v", err)
		return nil
	}
	return &parsedTime
}

// Utility function to handle nil strings
func getNonNilString(input *string) string {
	if input == nil {
		return ""
	}
	return *input
}
