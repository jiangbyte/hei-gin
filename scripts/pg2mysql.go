//go:build ignore

package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"
)

func main() {
	inPath := "scripts/db.sql"
	outPath := "scripts/db.mysql.sql"
	if len(os.Args) >= 2 {
		inPath = os.Args[1]
	}
	if len(os.Args) >= 3 {
		outPath = os.Args[2]
	}
	raw, err := os.ReadFile(inPath)
	if err != nil {
		fatal(err)
	}
	if !utf8.Valid(raw) {
		fatal(fmt.Errorf("input is not valid UTF-8"))
	}
	out := convert(string(raw))
	if err := os.WriteFile(outPath, []byte(out), 0o644); err != nil {
		fatal(err)
	}
	fmt.Println("wrote", outPath)
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

var (
	rePublicTable = regexp.MustCompile(`(?i)"public"\.`)
	reCollate     = regexp.MustCompile(`\s+COLLATE\s+"pg_catalog"\."default"`)
	reTimestamptz = regexp.MustCompile(`(?i)timestamptz\(\d+\)`)
	reInt4        = regexp.MustCompile(`(?i)\bint4\b`)
	reInt8        = regexp.MustCompile(`(?i)\bint8\b`)
	reFloat8      = regexp.MustCompile(`(?i)\bfloat8\b`)
	reBool        = regexp.MustCompile(`(?i)\bbool\b`)
	reDefaultNow  = regexp.MustCompile(`(?i)DEFAULT\s+now\(\)`)
	reJSONCast    = regexp.MustCompile(`'::json\b`)
	reJSONBCast   = regexp.MustCompile(`'::jsonb\b`)
	reTSOffset    = regexp.MustCompile(`'(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}(?:\.\d+)?)\+\d{2}'`)
	reBoolTrue    = regexp.MustCompile(`([^A-Za-z0-9_])'t'([^A-Za-z0-9_])`)
	reBoolFalse   = regexp.MustCompile(`([^A-Za-z0-9_])'f'([^A-Za-z0-9_])`)
	reCommentCol  = regexp.MustCompile(`(?i)^COMMENT ON (COLUMN|TABLE)\s+.+$`)
	reIndexOps    = regexp.MustCompile(`\s+"pg_catalog"\."[a-z0-9_]+"`)
	reUsingBtree  = regexp.MustCompile(`(?i)USING\s+btree`)
	reNullsOrder  = regexp.MustCompile(`(?i)\s+NULLS\s+(FIRST|LAST)`)
	reNumeric     = regexp.MustCompile(`(?i)\bnumeric\b`)
)

func convert(src string) string {
	var b strings.Builder
	b.WriteString("/*\n MySQL schema converted from scripts/db.sql (PostgreSQL dump).\n Charset: utf8mb4\n\n Usage:\n   mysql -u root -p -e \"CREATE DATABASE hei_gin DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;\"\n   mysql -u root -p hei_gin < scripts/db.mysql.sql\n*/\n\nSET NAMES utf8mb4;\nSET FOREIGN_KEY_CHECKS = 0;\n\n")

	sc := bufio.NewScanner(strings.NewReader(src))
	buf := make([]byte, 0, 1024*1024)
	sc.Buffer(buf, 8*1024*1024)

	inHeader := true
	skipBlock := false
	inCreateTable := false
	inComment := false
	var indexBuf strings.Builder
	inIndex := false

	flushIndex := func() {
		if !inIndex {
			return
		}
		block := indexBuf.String()
		indexBuf.Reset()
		inIndex = false
		lower := strings.ToLower(block)
		if strings.Contains(lower, " where ") || strings.Contains(lower, "using gin") || strings.Contains(block, "::") {
			return
		}
		b.WriteString(block)
	}

	for sc.Scan() {
		line := sc.Text()
		trim := strings.TrimSpace(line)

		if inHeader {
			if strings.HasPrefix(trim, "-- ----") {
				inHeader = false
			} else {
				continue
			}
		}

		// Skip COMMENT ON ... (may span multiple lines until ';')
		if inComment {
			if strings.HasSuffix(trim, ";") {
				inComment = false
			}
			continue
		}
		if strings.HasPrefix(strings.ToUpper(trim), "COMMENT ON") {
			if !strings.HasSuffix(trim, ";") {
				inComment = true
			}
			continue
		}

		lower := strings.ToLower(trim)
		if strings.Contains(lower, "using gin") ||
			strings.Contains(trim, "jsonb_ops") ||
			strings.Contains(trim, "::jsonb") {
			if inIndex {
				indexBuf.Reset()
				inIndex = false
			}
			skipBlock = true
		}
		if skipBlock {
			if strings.HasSuffix(trim, ";") {
				skipBlock = false
			}
			continue
		}

		if strings.HasPrefix(strings.ToUpper(trim), "CREATE INDEX") {
			flushIndex()
			inIndex = true
		}

		outLine := transformLine(line)
		outTrim := strings.TrimSpace(outLine)

		if inIndex {
			indexBuf.WriteString(outLine)
			indexBuf.WriteByte('\n')
			if strings.Contains(outTrim, ";") || outTrim == ";" {
				flushIndex()
			}
			continue
		}

		if strings.HasPrefix(strings.ToUpper(trim), "CREATE TABLE") {
			inCreateTable = true
		}

		if inCreateTable && (outTrim == ")" || outTrim == ");") {
			b.WriteString(") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;\n")
			inCreateTable = false
			continue
		}
		if outTrim == ";" {
			continue
		}
		if strings.HasSuffix(outTrim, ";") {
			inCreateTable = false
		}

		b.WriteString(outLine)
		b.WriteByte('\n')
	}
	flushIndex()

	b.WriteString("\nSET FOREIGN_KEY_CHECKS = 1;\n")
	return b.String()
}

func transformLine(line string) string {
	line = rePublicTable.ReplaceAllString(line, "")
	line = reCollate.ReplaceAllString(line, "")
	line = reTimestamptz.ReplaceAllString(line, "datetime(6)")
	line = reInt4.ReplaceAllString(line, "int")
	line = reInt8.ReplaceAllString(line, "bigint")
	line = reFloat8.ReplaceAllString(line, "double")
	line = reBool.ReplaceAllString(line, "tinyint(1)")
	line = reDefaultNow.ReplaceAllString(line, "DEFAULT CURRENT_TIMESTAMP(6)")
	line = reJSONCast.ReplaceAllString(line, "'")
	line = reJSONBCast.ReplaceAllString(line, "'")
	// MySQL JSON/BLOB/TEXT cannot have DEFAULT (except expression defaults in 8.0.13+; keep simple: drop)
	line = regexp.MustCompile(`(?i)(\bjson\b(?:\s+NOT\s+NULL)?)\s+DEFAULT\s+'[^']*'`).ReplaceAllString(line, "$1")
	line = regexp.MustCompile(`(?i)(\bjson\b(?:\s+NOT\s+NULL)?)\s+DEFAULT\s+\([^)]*\)`).ReplaceAllString(line, "$1")
	line = reTSOffset.ReplaceAllString(line, "'$1'")
	line = reUsingBtree.ReplaceAllString(line, "")
	line = reIndexOps.ReplaceAllString(line, "")
	line = reNullsOrder.ReplaceAllString(line, "")
	line = reNumeric.ReplaceAllString(line, "decimal(20,6)")
	line = strings.ReplaceAll(line, " ASC,", ",")
	line = strings.ReplaceAll(line, " DESC,", " DESC,")
	line = regexp.MustCompile(`\s+ASC\s*$`).ReplaceAllString(line, "")
	line = regexp.MustCompile(`\s+DESC\s*$`).ReplaceAllString(line, " DESC")
	line = strings.ReplaceAll(line, "::json", "")
	line = strings.ReplaceAll(line, "::jsonb", "")
	// double spaces from removed USING btree
	line = regexp.MustCompile(`[ \t]{2,}`).ReplaceAllString(line, " ")
	line = strings.ReplaceAll(line, "( ", "(")
	upper := strings.TrimSpace(strings.ToUpper(line))
	if strings.HasPrefix(upper, "INSERT ") {
		line = reBoolTrue.ReplaceAllString(line, "${1}1${2}")
		line = reBoolFalse.ReplaceAllString(line, "${1}0${2}")
		if idx := strings.Index(strings.ToUpper(line), " VALUES"); idx > 0 {
			head := quoteIdentsOutsideStrings(line[:idx])
			line = head + line[idx:]
		}
		return line
	}

	return quoteIdentsOutsideStrings(line)
}

func quoteIdentsOutsideStrings(s string) string {
	var b strings.Builder
	inSingle := false
	for i := 0; i < len(s); {
		ch := s[i]
		if inSingle {
			b.WriteByte(ch)
			if ch == '\'' {
				if i+1 < len(s) && s[i+1] == '\'' {
					b.WriteByte('\'')
					i += 2
					continue
				}
				inSingle = false
			}
			i++
			continue
		}
		if ch == '\'' {
			inSingle = true
			b.WriteByte(ch)
			i++
			continue
		}
		if ch == '"' {
			j := i + 1
			for j < len(s) && s[j] != '"' {
				j++
			}
			if j < len(s) {
				b.WriteByte('`')
				b.WriteString(s[i+1 : j])
				b.WriteByte('`')
				i = j + 1
				continue
			}
		}
		b.WriteByte(ch)
		i++
	}
	return b.String()
}
