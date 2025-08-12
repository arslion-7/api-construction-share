package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/arslion-7/api-construction-share/initializers"
	"github.com/arslion-7/api-construction-share/models"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Load environment variables
	initializers.LoadEnvVars()
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	fmt.Println("Starting migration of onki_bd_2019 table from MySQL to PostgreSQL...")
	fmt.Println("PostgreSQL Host:", dbHost)

	// Connect to MySQL
	mysqlDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", "root", "010203", "192.168.0.242", "3306", "paylynew")
	mysqlDB, err := sql.Open("mysql", mysqlDSN)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	defer mysqlDB.Close()

	// Test MySQL connection
	if err := mysqlDB.Ping(); err != nil {
		log.Fatalf("Failed to ping MySQL: %v", err)
	}
	fmt.Println("Successfully connected to MySQL")

	// Connect to PostgreSQL
	postgresDSN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	PostgresDB, err := gorm.Open(postgres.Open(postgresDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	// Auto migrate the OldRegistry table
	if err := PostgresDB.AutoMigrate(&models.OldRegistry{}); err != nil {
		log.Fatalf("Failed to auto migrate OldRegistry table: %v", err)
	}
	fmt.Println("Successfully created old_registries table in PostgreSQL")

	// Query MySQL for data from onki_bd_2019 table
	rows, err := mysqlDB.Query("SELECT * FROM onki_bd_2019")
	if err != nil {
		log.Fatalf("Failed to query MySQL table 'onki_bd_2019': %v", err)
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		log.Fatalf("Failed to get columns: %v", err)
	}

	// Create a slice of interface{} to hold the values
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	// Process rows and migrate to PostgreSQL
	processedCount := 0
	insertedCount := 0
	for rows.Next() {
		if err := rows.Scan(valuePtrs...); err != nil {
			log.Printf("Row scan failed: %v", err)
			continue
		}

		// Create a map to hold the column values
		rowData := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			rowData[col] = val
		}

		// Convert row data to OldRegistry struct
		oldRegistry := convertRowToOldRegistry(rowData)

		// Always create new record (allow duplicates)
		if err := PostgresDB.Create(&oldRegistry).Error; err != nil {
			log.Printf("Failed to insert old registry (TB: %d): %v", oldRegistry.TB, err)
		} else {
			processedCount++
			insertedCount++
			if processedCount%100 == 0 {
				log.Printf("Processed %d records (inserted: %d)...", processedCount, insertedCount)
			}
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Row iteration error: %v", err)
	}

	log.Printf("Migration completed successfully! Processed %d records (inserted: %d).", processedCount, insertedCount)
}

// convertRowToOldRegistry converts a row from MySQL to OldRegistry struct
func convertRowToOldRegistry(rowData map[string]interface{}) models.OldRegistry {
	oldRegistry := models.OldRegistry{}

	// Convert TB (primary key)
	if tb, ok := rowData["t_b"]; ok && tb != nil {
		switch v := tb.(type) {
		case int64:
			oldRegistry.TB = uint(v)
		case int32:
			oldRegistry.TB = uint(v)
		case int:
			oldRegistry.TB = uint(v)
		case uint64:
			oldRegistry.TB = uint(v)
		case uint32:
			oldRegistry.TB = uint(v)
		case uint:
			oldRegistry.TB = v
		case float64:
			oldRegistry.TB = uint(v)
		case float32:
			oldRegistry.TB = uint(v)
		case []byte:
			// Try to convert string to int
			if strVal := string(v); strVal != "" {
				if intVal, err := strconv.ParseUint(strVal, 10, 64); err == nil {
					oldRegistry.TB = uint(intVal)
				}
			}
		case string:
			if intVal, err := strconv.ParseUint(v, 10, 64); err == nil {
				oldRegistry.TB = uint(intVal)
			}
		}
	}

	// Convert string fields
	stringFields := []string{
		"min_hat", "gurujy", "paychy", "sertnama_gurujy_paychy", "desga", "baha_umumy",
		"meydan_umumy", "kep_resminama", "emlak_paychy", "baha_paychy", "baha_1m2_paychy",
		"salgy_desga", "salgy_gurujy", "salgy_paychy", "bash_potr", "sertnama_gur_potr",
		"potratchy_komek", "shahadatnama", "ygtyyarnama", "patent_pasport", "sene_bashy_songy",
		"sene_hasaba_alnan", "wezipe_alan_adam", "ady_alan_adam", "sene_san_sertnama",
		"ady_paychy_alan", "sene_paychy_alan", "login",
	}

	for _, field := range stringFields {
		if val, ok := rowData[field]; ok && val != nil {
			if str, ok := val.([]byte); ok {
				strVal := string(str)
				switch field {
				case "min_hat":
					oldRegistry.MinHat = &strVal
				case "gurujy":
					oldRegistry.Gurujy = &strVal
				case "paychy":
					oldRegistry.Paychy = &strVal
				case "sertnama_gurujy_paychy":
					oldRegistry.SertnamaGurujyPaychy = &strVal
				case "desga":
					oldRegistry.Desga = &strVal
				case "baha_umumy":
					oldRegistry.BahaUmumy = &strVal
				case "meydan_umumy":
					oldRegistry.MeydanUmumy = &strVal
				case "kep_resminama":
					oldRegistry.KepResminama = &strVal
				case "emlak_paychy":
					oldRegistry.EmlakPaychy = &strVal
				case "baha_paychy":
					oldRegistry.BahaPaychy = &strVal
				case "baha_1m2_paychy":
					oldRegistry.Baha1m2Paychy = &strVal
				case "salgy_desga":
					oldRegistry.SalgyDesga = &strVal
				case "salgy_gurujy":
					oldRegistry.SalgyGurujy = &strVal
				case "salgy_paychy":
					oldRegistry.SalgyPaychy = &strVal
				case "bash_potr":
					oldRegistry.BashPotr = &strVal
				case "sertnama_gur_potr":
					oldRegistry.SertnamaGurPotr = &strVal
				case "potratchy_komek":
					oldRegistry.PotratchyKomek = &strVal
				case "shahadatnama":
					oldRegistry.Shahadatnama = &strVal
				case "ygtyyarnama":
					oldRegistry.Ygtyyarnama = &strVal
				case "patent_pasport":
					oldRegistry.PatentPasport = &strVal
				case "sene_bashy_songy":
					oldRegistry.SeneBashySongy = &strVal
				case "sene_hasaba_alnan":
					oldRegistry.SeneHasabaAlnan = &strVal
				case "wezipe_alan_adam":
					oldRegistry.WezipeAlanAdam = &strVal
				case "ady_alan_adam":
					oldRegistry.AdyAlanAdam = &strVal
				case "sene_san_sertnama":
					oldRegistry.SeneSanSertnama = &strVal
				case "ady_paychy_alan":
					oldRegistry.AdyPaychyAlan = &strVal
				case "sene_paychy_alan":
					oldRegistry.SenePaychyAlan = &strVal
				case "login":
					oldRegistry.Login = &strVal
				}
			}
		}
	}

	// Convert date fields
	dateFields := []string{"sene_hat_min_to_mud", "sene_seredilen"}
	for _, field := range dateFields {
		if val, ok := rowData[field]; ok && val != nil {
			fmt.Printf("DEBUG: %s type: %T, value: %v\n", field, val, val)

			switch v := val.(type) {
			case time.Time:
				switch field {
				case "sene_hat_min_to_mud":
					oldRegistry.SeneHatMinToMud = &v
					fmt.Printf("DEBUG: Converted time.Time for %s: %v\n", field, v)
				case "sene_seredilen":
					oldRegistry.SeneSeredilen = &v
					fmt.Printf("DEBUG: Converted time.Time for %s: %v\n", field, v)
				}
			case []byte:
				// Try to parse date string
				if strVal := string(v); strVal != "" {
					if timeVal, err := time.Parse("2006-01-02", strVal); err == nil {
						switch field {
						case "sene_hat_min_to_mud":
							oldRegistry.SeneHatMinToMud = &timeVal
							fmt.Printf("DEBUG: Converted []byte to date for %s: %s -> %v\n", field, strVal, timeVal)
						case "sene_seredilen":
							oldRegistry.SeneSeredilen = &timeVal
							fmt.Printf("DEBUG: Converted []byte to date for %s: %s -> %v\n", field, strVal, timeVal)
						}
					} else {
						fmt.Printf("DEBUG: Failed to parse date from []byte: %s, error: %v\n", strVal, err)
					}
				}
			case string:
				// Try to parse date string
				if timeVal, err := time.Parse("2006-01-02", v); err == nil {
					switch field {
					case "sene_hat_min_to_mud":
						oldRegistry.SeneHatMinToMud = &timeVal
						fmt.Printf("DEBUG: Converted string to date for %s: %s -> %v\n", field, v, timeVal)
					case "sene_seredilen":
						oldRegistry.SeneSeredilen = &timeVal
						fmt.Printf("DEBUG: Converted string to date for %s: %s -> %v\n", field, v, timeVal)
					}
				} else {
					fmt.Printf("DEBUG: Failed to parse date from string: %s, error: %v\n", v, err)
				}
			default:
				fmt.Printf("DEBUG: Unknown type for date field %s: %T, value: %v\n", field, v, v)
			}
		} else {
			fmt.Printf("DEBUG: %s is nil or not found\n", field)
		}
	}

	return oldRegistry
}
