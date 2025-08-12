# Onki BD 2019 Migration

This migration script transfers data from the MySQL `onki_bd_2019` table to the PostgreSQL `old_registries` table.

## Prerequisites

1. MySQL server running on `192.168.0.242` with:

   - Database: `paylynew`
   - User: `root`
   - Password: `010203`

2. PostgreSQL database configured with environment variables:
   - `DB_HOST`
   - `DB_PORT`
   - `DB_USER`
   - `DB_PASSWORD`
   - `DB_NAME`

## Running the Migration

1. Make sure your PostgreSQL database is running and accessible
2. Navigate to this directory: `cd api-construction-share/migrations/onki_bd_2019_migration`
3. Run the migration: `go run main.go`

## What it does

1. Connects to both MySQL and PostgreSQL databases
2. Creates the `old_registries` table in PostgreSQL (if it doesn't exist)
3. Reads all records from the MySQL `onki_bd_2019` table
4. Transfers the data to PostgreSQL, preserving all field names and data types
5. Logs progress every 100 records

## Source Table Structure

The migration reads from the MySQL `onki_bd_2019` table with the following structure:

```sql
CREATE TABLE IF NOT EXISTS `onki_bd_2019` (
  `t_b` int(11) NOT NULL,
  `min_hat` varchar(150) DEFAULT NULL,
  `sene_hat_min_to_mud` date DEFAULT NULL,
  `gurujy` varchar(155) DEFAULT NULL,
  `paychy` varchar(100) DEFAULT NULL,
  `sertnama_gurujy_paychy` varchar(100) DEFAULT NULL,
  `desga` varchar(275) DEFAULT NULL,
  `baha_umumy` varchar(20) DEFAULT '0.00',
  `meydan_umumy` varchar(125) DEFAULT NULL,
  `kep_resminama` varchar(225) DEFAULT NULL,
  `emlak_paychy` varchar(325) DEFAULT NULL,
  `baha_paychy` varchar(20) DEFAULT '0.00',
  `baha_1m2_paychy` decimal(10,2) DEFAULT '0.00',
  `salgy_desga` varchar(255) DEFAULT NULL,
  `salgy_gurujy` varchar(150) DEFAULT NULL,
  `salgy_paychy` varchar(150) DEFAULT NULL,
  `bash_potr` varchar(100) DEFAULT NULL,
  `sertnama_gur_potr` varchar(155) DEFAULT NULL,
  `potratchy_komek` varchar(125) DEFAULT NULL,
  `shahadatnama` varchar(120) DEFAULT NULL,
  `ygtyyarnama` varchar(255) DEFAULT NULL,
  `patent_pasport` varchar(255) DEFAULT NULL,
  `sene_bashy_songy` varchar(100) DEFAULT NULL,
  `sene_seredilen` date DEFAULT NULL,
  `sene_hasaba_alnan` varchar(125) DEFAULT NULL,
  `wezipe_alan_adam` varchar(125) DEFAULT NULL,
  `ady_alan_adam` varchar(125) DEFAULT NULL,
  `sene_san_sertnama` varchar(55) DEFAULT NULL,
  `ady_paychy_alan` varchar(125) DEFAULT NULL,
  `sene_paychy_alan` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`t_b`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Payly gurlushyk mudirliginin reyestrini yoretmek ucin programmasy';
```

## Target Table Structure

The data is migrated to the PostgreSQL `old_registries` table which contains all the original fields plus standard GORM fields:

- `t_b` (Indexed field for quick search)
- `min_hat`
- `sene_hat_min_to_mud`
- `gurujy`
- `paychy`
- `sertnama_gurujy_paychy`
- `desga`
- `baha_umumy`
- `meydan_umumy`
- `kep_resminama`
- `emlak_paychy`
- `baha_paychy`
- `baha_1m2_paychy`
- `salgy_desga`
- `salgy_gurujy`
- `salgy_paychy`
- `bash_potr`
- `sertnama_gur_potr`
- `potratchy_komek`
- `shahadatnama`
- `ygtyyarnama`
- `patent_pasport`
- `sene_bashy_songy`
- `sene_seredilen`
- `sene_hasaba_alnan`
- `wezipe_alan_adam`
- `ady_alan_adam`
- `sene_san_sertnama`
- `ady_paychy_alan`
- `sene_paychy_alan`
- `login`

Plus standard GORM fields:

- `id`
- `created_at`
- `updated_at`
- `deleted_at`

## Migration Strategy

The migration uses a simple insert strategy:

- Creates a new record for every row in the source table
- Allows duplicate `t_b` values (multiple records can have the same `t_b`)
- The `t_b` field is indexed for quick searching but is not the primary key

This ensures that:

- All source data is preserved
- Duplicate `t_b` values are allowed
- Fast queries can be performed using the `t_b` field
- Simple and straightforward migration process

## Example Output

```
Starting migration of onki_bd_2019 table from MySQL to PostgreSQL...
PostgreSQL Host: localhost
Successfully connected to MySQL
Successfully created old_registries table in PostgreSQL
Processed 100 records (inserted: 100)...
Processed 200 records (inserted: 200)...
Migration completed successfully! Processed 250 records (inserted: 250).
```
