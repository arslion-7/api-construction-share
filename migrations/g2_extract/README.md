# G2 (Goşmaça Şertnama) Dump Extractor

Extracts all additional agreement (goşmaça şertnama) data from a legacy
`paylynew` MySQL dump file (HeidiSQL / mysqldump format) into one merged CSV.
Works with any dump that contains the `g2_*` tables, so it can be re-run on
newer dumps as they arrive.

## Source tables

| table | role |
|---|---|
| `g2_mainpayly` | primary ledger (most complete): ministry letter, shareholder, review/registration text, reason, login |
| `g2_min_hat` | ministry letter as typed fields (dates, number, count, tapgyr) |
| `g2_goshmacha_maglumat` | shareholder person/org details |
| `g2_sebabi` | reason category + detail |
| `g2_seneler` | typed review date, PGGŞ registration number and date |

The four normalized tables share one `t_b` key and are merged into the
`g2_mainpayly` row with the same id. Rows present in only one source are kept.
NULL values become empty CSV cells.

Note: the `g2_*` `t_b` is an independent id sequence — it is NOT the main
registry (`mainpayly`) number. Linking an agreement to a registry must be done
during import, by shareholder name.

## Usage

```bash
cd api-construction-share/migrations/g2_extract
go run main.go -dump /path/to/dump.sql -out g2_mainpayly.csv
```

Prints per-table row counts and writes the merged CSV (UTF-8, header row,
multi-line values quoted).

## Output columns

`t_b`,
`min_hat`, `paychy`, `sene_sered_hasaba_alnan`, `sebabi`, `login` (raw ledger),
`hat_sene`, `hat_belgi`, `tabsh_sene`, `sany`, `tapgyr`, `min_to_mud` (letter),
`telekeci_rayat`, `gurama_ady`, `wezipe`, `familiya`, `ady`, `otchestvo`,
`person_gosh` (person),
`sebabi_esasy`, `sebabi_goshmacha` (reason),
`seredilen_date`, `hasaba_alnan_pggs`, `hasaba_alnan_belgisi`,
`hasaba_alnan_date`, `seneler_gosh` (registration).

Do not commit generated CSVs — they contain personal data.
