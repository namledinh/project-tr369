package utils

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
	"usp-management-device-api/business/models"
	apperrors "usp-management-device-api/common/app_errors"
	"usp-management-device-api/common/logging"
)

func GenerateUUID() string {
	return uuid.New().String()
}

func MacCleaner(mac string) string {
	lowerMac := strings.ToLower(mac)
	trimSpace := strings.ReplaceAll(lowerMac, " ", "")
	return strings.ReplaceAll(trimSpace, ":", "")
}

func StringCleaner(str string) string {
	lowerStr := strings.ToLower(str)
	return strings.ReplaceAll(lowerStr, " ", "")
}

func RemoveDuplicateString(strSlice []string) []string {
	slices.Sort(strSlice)
	return slices.Compact(strSlice)
}

// GetLastNumber extracts the last number from the input string.
// It returns the number as an integer and an error if no number is found.
func GetLastNumber(input string) (int, error) {
	// Define a regular expression to match one or more digits
	re := regexp.MustCompile(`\d+`)

	// Find all matches of the regex in the input string
	matches := re.FindAllString(input, -1)

	// Check if any matches were found
	if len(matches) == 0 {
		return 0, errors.New("no numbers found in the input string: " + input)
	}

	// Get the last matched number as a string
	lastNumberStr := matches[len(matches)-1]

	// Convert the string to an integer
	lastNumber, err := strconv.Atoi(lastNumberStr)
	if err != nil {
		return 0, fmt.Errorf("error converting '%s' to integer: %v", lastNumberStr, err)
	}

	return lastNumber, nil
}

// SubmitBackgroundJob runs a function in a background goroutine with panic recovery.
func SubmitBackgroundJob(fn func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logging.Errorf("panic occurred: %v", err)
			}
		}()
		fn()
	}()
}

func InsertColonToMacAddress(mac string) string {
	if strings.Contains(mac, ":") || len(mac) < 12 {
		return mac
	}
	var result string
	for i := 0; i < len(mac); i += 2 {
		if i > 0 {
			result += ":"
		}
		result += mac[i : i+2]
	}
	return result
}

func ContainString(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

var alphabet []rune = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

func RandomString(n int) string {

	alphabetSize := len(alphabet)
	var sb strings.Builder

	for i := 0; i < n; i++ {
		ch := alphabet[rand.Intn(alphabetSize)]
		sb.WriteRune(ch)
	}

	s := sb.String()
	return s
}

func stripOuterParens(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 && s[0] == '(' && s[len(s)-1] == ')' {
		depth := 0
		for i := 0; i < len(s); i++ {
			switch s[i] {
			case '\'':
			case '(':
				depth++
			case ')':
				depth--
				if depth == 0 && i != len(s)-1 {
					return s
				}
			}
		}
		return s[1 : len(s)-1]
	}
	return s
}

func splitTopLevel(expr string) (segments []string, joins []string) {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return nil, nil
	}

	var b strings.Builder
	depth := 0
	inQuote := false

	lower := strings.ToLower(expr)
	i := 0
	for i < len(expr) {
		ch := expr[i]

		if ch == '\'' {
			inQuote = !inQuote
			b.WriteByte(ch)
			i++
			continue
		}
		if !inQuote {
			if ch == '(' {
				depth++
				b.WriteByte(ch)
				i++
				continue
			}
			if ch == ')' {
				if depth > 0 {
					depth--
				}
				b.WriteByte(ch)
				i++
				continue
			}
			if depth == 0 {
				if strings.HasPrefix(lower[i:], " and ") {
					segments = append(segments, strings.TrimSpace(b.String()))
					joins = append(joins, "AND")
					b.Reset()
					i += len(" and ")
					continue
				}
				if strings.HasPrefix(lower[i:], " or ") {
					segments = append(segments, strings.TrimSpace(b.String()))
					joins = append(joins, "OR")
					b.Reset()
					i += len(" or ")
					continue
				}
			}
		}

		b.WriteByte(ch)
		i++
	}
	if b.Len() > 0 {
		segments = append(segments, strings.TrimSpace(b.String()))
	}
	return
}

func splitInnerByAndOr(expr string) (parts []string, joins []string) {
	return splitTopLevel(expr)
}

func cleanValue(v string) string {
	v = strings.TrimSpace(v)
	v = stripOuterParens(v)
	if len(v) >= 2 && ((v[0] == '\'' && v[len(v)-1] == '\'') || (v[0] == '"' && v[len(v)-1] == '"')) {
		v = v[1 : len(v)-1]
	}
	return strings.TrimSpace(v)
}

func parseCondition(token string) (field, op, value string, err error) {
	re := regexp.MustCompile(`(?i)^\s*([a-zA-Z0-9_\.]+)\s+(eq|ne|lt|gt|lte|gte|like)\s+(.+?)\s*$`)
	m := re.FindStringSubmatch(token)
	if len(m) != 4 {
		return "", "", "", fmt.Errorf("invalid condition: %s", token)
	}
	field = strings.TrimSpace(m[1])
	op = strings.ToLower(strings.TrimSpace(m[2]))
	value = cleanValue(m[3])
	return
}

func ParseFilterExpr(raw string) ([]models.FilterExpr, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}

	raw = stripOuterParens(raw)

	segments, topJoins := splitTopLevel(raw)
	if len(segments) == 0 {
		return nil, apperrors.NewDBError(nil, "invalid filter expression")
	}

	var out []models.FilterExpr

	for si, seg := range segments {
		seg = stripOuterParens(seg) // "(a or b)" -> "a or b"

		innerParts, innerJoins := splitInnerByAndOr(seg)
		if len(innerParts) == 0 {
			return nil, apperrors.NewDBError(nil, "invalid segment: "+seg)
		}

		for pi, part := range innerParts {
			field, op, val, err := parseCondition(part)
			if err != nil {
				return nil, apperrors.NewDBError(err, "invalid filter condition")
			}

			join := "AND" //default is AND
			if si == 0 && pi == 0 {
				join = "AND"
			} else if pi > 0 {
				join = strings.ToUpper(innerJoins[pi-1])
			} else {
				join = strings.ToUpper(topJoins[si-1])
			}

			out = append(out, models.FilterExpr{
				Filter: field,
				Op:     op,
				Value:  val,
				Join:   join,
			})
			logging.Infof("Parsed filter: %+v", out[len(out)-1])
		}
	}

	logging.Infof("Final parsed filters: %+v", out)
	return out, nil
}

func sqlOp(op string) string {
	switch strings.ToLower(op) {
	case "eq":
		return "="
	case "ne":
		return "!="
	case "lt":
		return "<"
	case "gt":
		return ">"
	case "lte":
		return "<="
	case "gte":
		return ">="
	case "like":
		return "LIKE"
	default:
		return "="
	}
}

func BuildCond(f models.FilterExpr) (string, []any) {
	return fmt.Sprintf("%s %s ?", f.Filter, sqlOp(f.Op)), []any{f.Value}
}

func ParseOrderExpr(raw string) ([]models.OrderExpr, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}

	parts := strings.Split(raw, ",")
	var out []models.OrderExpr

	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}

		tokens := strings.Fields(p)
		if len(tokens) == 0 {
			continue
		}

		field := tokens[0]
		direction := "DESC"
		if len(tokens) > 1 {
			dir := strings.ToUpper(tokens[1])
			if dir == "ASC" || dir == "DESC" {
				direction = dir
			} else {
				return nil, fmt.Errorf("invalid order direction: %s", tokens[1])
			}
		}

		out = append(out, models.OrderExpr{
			Field:     field,
			Direction: direction,
		})
	}

	return out, nil
}

func ConvertToCSV(data [][]string) []byte {
	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)

	_ = writer.WriteAll(data)

	return buf.Bytes()
}

// FormatTimeWithTimezone formats time from UTC to specified timezone
func FormatTimeWithTimezone(t *time.Time, layout string, timezone string) string {
	if t == nil {
		return ""
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		// If timezone loading fails, use UTC
		return t.Format(layout)
	}

	return t.In(loc).Format(layout)
}

// FormatTimeGMT7 formats time from UTC to GMT+7 (Asia/Ho_Chi_Minh)
func FormatTimeGMT7(t *time.Time, layout string) string {
	return FormatTimeWithTimezone(t, layout, "Asia/Ho_Chi_Minh")
}
