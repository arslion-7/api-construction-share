// g2_extract reads a HeidiSQL/mysqldump .sql dump of the legacy `paylynew`
// database and extracts all goşmaça şertnama (additional agreement) data into
// a single merged CSV.
//
// Source tables:
//   - g2_mainpayly (primary, most complete): min_hat, paychy,
//     seneSeredHasabaAlnan, sebabi, login
//   - g2_min_hat, g2_goshmacha_maglumat, g2_sebabi, g2_seneler: normalized
//     column groups sharing the same t_b, merged in by id where present
//
// Usage:
//
//	go run main.go -dump /path/to/dump.sql -out g2_mainpayly.csv
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// tableSection returns the dump fragment between DISABLE KEYS and ENABLE KEYS
// for the given table, which contains all its REPLACE INTO statements.
func tableSection(dump, table string) string {
	start := strings.Index(dump, "ALTER TABLE `"+table+"` DISABLE KEYS")
	if start < 0 {
		return ""
	}
	rest := dump[start:]
	end := strings.Index(rest, "ALTER TABLE `"+table+"` ENABLE KEYS")
	if end < 0 {
		return rest
	}
	return rest[:end]
}

var escapes = map[byte]string{
	'n': "\n", 'r': "\r", 't': "\t", '0': "\x00",
}

// parseTuples parses every VALUES (...),(...) tuple in a table section.
// NULL becomes an empty string; MySQL escapes are decoded.
func parseTuples(section string) [][]string {
	var rows [][]string
	i := 0
	n := len(section)
	for {
		j := strings.Index(section[i:], "VALUES")
		if j < 0 {
			break
		}
		i += j + len("VALUES")
		for i < n {
			for i < n && (section[i] == ' ' || section[i] == '\t' || section[i] == '\r' || section[i] == '\n' || section[i] == ',') {
				i++
			}
			if i >= n || section[i] == ';' {
				break
			}
			if section[i] != '(' {
				break
			}
			i++
			var vals []string
			var buf strings.Builder
			isString := false
			inString := false
			for i < n {
				c := section[i]
				if inString {
					if c == '\\' && i+1 < n {
						if repl, ok := escapes[section[i+1]]; ok {
							buf.WriteString(repl)
						} else {
							buf.WriteByte(section[i+1])
						}
						i += 2
						continue
					}
					if c == '\'' {
						if i+1 < n && section[i+1] == '\'' {
							buf.WriteByte('\'')
							i += 2
							continue
						}
						inString = false
						i++
						continue
					}
					buf.WriteByte(c)
					i++
					continue
				}
				switch c {
				case '\'':
					// whitespace before the opening quote is separator junk
					if strings.TrimSpace(buf.String()) == "" {
						buf.Reset()
					}
					inString = true
					isString = true
					i++
				case ',':
					vals = append(vals, finishValue(&buf, &isString))
					i++
				case ')':
					vals = append(vals, finishValue(&buf, &isString))
					i++
				default:
					buf.WriteByte(c)
					i++
				}
				if c == ')' {
					break
				}
			}
			rows = append(rows, vals)
		}
	}
	return rows
}

func finishValue(buf *strings.Builder, isString *bool) string {
	v := buf.String()
	buf.Reset()
	if !*isString {
		v = strings.TrimSpace(v)
		if v == "NULL" {
			v = ""
		}
	}
	*isString = false
	return v
}

// tableByID parses a table section into id -> row values (id column dropped),
// validating the expected column count.
func tableByID(dump, table string, columns int) map[int][]string {
	out := map[int][]string{}
	rows := parseTuples(tableSection(dump, table))
	bad := 0
	for _, row := range rows {
		if len(row) != columns {
			bad++
			continue
		}
		id, err := strconv.Atoi(row[0])
		if err != nil {
			bad++
			continue
		}
		out[id] = row[1:]
	}
	fmt.Printf("%-24s rows=%d skipped=%d\n", table, len(out), bad)
	return out
}

func at(row []string, idx int) string {
	if row == nil || idx >= len(row) {
		return ""
	}
	return row[idx]
}

func main() {
	dumpPath := flag.String("dump", "", "path to the .sql dump file (required)")
	outPath := flag.String("out", "g2_mainpayly.csv", "output CSV path")
	flag.Parse()
	if *dumpPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	raw, err := os.ReadFile(*dumpPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to read dump:", err)
		os.Exit(1)
	}
	dump := string(raw)

	flat := tableByID(dump, "g2_mainpayly", 6)          // min_hat, paychy, seneSeredHasabaAlnan, sebabi, login
	minHat := tableByID(dump, "g2_min_hat", 7)          // hat_sene, hat_belgi, tabsh_sene, sany, tapgyr, minToMud
	person := tableByID(dump, "g2_goshmacha_maglumat", 8) // telekeci_rayat, guramaAdy, wezipe, familiya, ady, otchestvo, gosh1
	sebabi := tableByID(dump, "g2_sebabi", 3)           // sebabi_esasy, sebabi_goshmacha
	seneler := tableByID(dump, "g2_seneler", 6)         // seredilenDate, hasabaAlnanPGGS, hasabaAlnanBelgisi, hasabaAlnanDate, gosh1

	idSet := map[int]bool{}
	for id := range flat {
		idSet[id] = true
	}
	for _, m := range []map[int][]string{minHat, person, sebabi, seneler} {
		for id := range m {
			idSet[id] = true
		}
	}
	ids := make([]int, 0, len(idSet))
	for id := range idSet {
		ids = append(ids, id)
	}
	sort.Ints(ids)

	f, err := os.Create(*outPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to create output:", err)
		os.Exit(1)
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	header := []string{
		"t_b",
		// g2_mainpayly (raw)
		"min_hat", "paychy", "sene_sered_hasaba_alnan", "sebabi", "login",
		// g2_min_hat (typed)
		"hat_sene", "hat_belgi", "tabsh_sene", "sany", "tapgyr", "min_to_mud",
		// g2_goshmacha_maglumat (person)
		"telekeci_rayat", "gurama_ady", "wezipe", "familiya", "ady", "otchestvo", "person_gosh",
		// g2_sebabi (reason)
		"sebabi_esasy", "sebabi_goshmacha",
		// g2_seneler (dates/registration)
		"seredilen_date", "hasaba_alnan_pggs", "hasaba_alnan_belgisi", "hasaba_alnan_date", "seneler_gosh",
	}
	if err := w.Write(header); err != nil {
		fmt.Fprintln(os.Stderr, "write failed:", err)
		os.Exit(1)
	}

	for _, id := range ids {
		fr := flat[id]
		mh := minHat[id]
		pe := person[id]
		se := sebabi[id]
		sn := seneler[id]
		row := []string{
			strconv.Itoa(id),
			at(fr, 0), at(fr, 1), at(fr, 2), at(fr, 3), at(fr, 4),
			at(mh, 0), at(mh, 1), at(mh, 2), at(mh, 3), at(mh, 4), at(mh, 5),
			at(pe, 0), at(pe, 1), at(pe, 2), at(pe, 3), at(pe, 4), at(pe, 5), at(pe, 6),
			at(se, 0), at(se, 1),
			at(sn, 0), at(sn, 1), at(sn, 2), at(sn, 3), at(sn, 4),
		}
		// legacy values carry padded spaces and trailing tabs/newlines
		for k := range row {
			row[k] = strings.TrimSpace(row[k])
		}
		if err := w.Write(row); err != nil {
			fmt.Fprintln(os.Stderr, "write failed:", err)
			os.Exit(1)
		}
	}

	fmt.Printf("\nwrote %d agreements to %s\n", len(ids), *outPath)
}
