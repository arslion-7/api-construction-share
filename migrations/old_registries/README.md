# Old Registries Migration

This migration script transfers data from the MySQL `mainpayly` table to the PostgreSQL `old_registries` table.

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
2. Navigate to this directory: `cd api-construction-share/migrations/old_registries`
3. Run the migration: `go run main.go`

## What it does

1. Connects to both MySQL and PostgreSQL databases
2. Creates the `old_registries` table in PostgreSQL (if it doesn't exist)
3. Reads all records from the MySQL `mainpayly` table
4. Transfers the data to PostgreSQL, preserving all field names and data types
5. Logs progress every 100 records

## Table Structure

The `old_registries` table contains all the original fields from the MySQL `mainpayly` table:

- `t_b` (Primary Key)
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
