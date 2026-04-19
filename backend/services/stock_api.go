package services

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"backend/models"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// SinaFinanceAPI fetches stock quote from Sina Finance API
func SinaFinanceAPI(code string) (*models.StockQuote, error) {
	// Convert stock code to Sina format
	// A shares: sh600519 (Shanghai), sz000001 (Shenzhen)
	sinaCode := convertToSinaCode(code)

	url := fmt.Sprintf("http://hq.sinajs.cn/list=%s", sinaCode)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Referer", "http://finance.sina.com.cn")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch quote: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return parseSinaResponse(string(body), code)
}

// convertToSinaCode converts stock code to Sina Finance format
func convertToSinaCode(code string) string {
	code = strings.ToLower(code)
	// If already in sh/sz format, return as is
	if strings.HasPrefix(code, "sh") || strings.HasPrefix(code, "sz") {
		return code
	}
	// A share detection - numeric codes
	// Shanghai: 600xxx, 601xxx, 603xxx, 688xxx
	// Shenzhen: 000xxx, 001xxx, 002xxx, 300xxx
	if len(code) == 6 {
		prefix := code[:3]
		switch {
		case prefix >= "600" && prefix <= "603" || prefix == "688":
			return "sh" + code
		case prefix >= "000" && prefix <= "003" || prefix >= "200" && prefix <= "299" || prefix >= "300" && prefix <= "399":
			return "sz" + code
		}
	}
	// Default to Shanghai
	return "sh" + code
}

// parseSinaResponse parses the Sina Finance API response
func parseSinaResponse(body string, code string) (*models.StockQuote, error) {
	// Response format: var hq_str_sh600519="贵州茅台,1700.00,1688.00,1695.00,1705.00,1680.00,1698.00,1699.00,5000000,8500000,...";
	prefix := `var hq_str_`
	idx := strings.Index(body, prefix)
	if idx == -1 {
		return nil, fmt.Errorf("invalid response format")
	}

	// Find the equals sign and quote
	start := strings.Index(body[idx:], "=")
	if start == -1 {
		return nil, fmt.Errorf("invalid response format")
	}
	start = idx + start + 1

	// Find the semicolon
	end := strings.Index(body[start:], ";")
	if end == -1 {
		return nil, fmt.Errorf("invalid response format")
	}

	dataStr := body[start : start+end]
	dataStr = strings.Trim(dataStr, `"`)

	if dataStr == "" || dataStr == "none" {
		return nil, fmt.Errorf("stock not found: %s", code)
	}

	// Parse CSV data
	reader := csv.NewReader(strings.NewReader(dataStr))
	fields, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to parse CSV: %w", err)
	}

	if len(fields) < 33 {
		return nil, fmt.Errorf("insufficient data fields")
	}

	quote := &models.StockQuote{
		Code: code,
	}

	// Parse fields (convert from GBK to UTF-8)
	quote.Name = gbkToUtf8(fields[0])

	quote.Open = parseFloat(fields[1])
	quote.Current = parseFloat(fields[3]) // fields[2] is yesterday close, fields[3] is current price
	quote.High = parseFloat(fields[4])
	quote.Low = parseFloat(fields[5])

	// Volume (shares)
	quote.Volume = parseInt64(fields[8])

	// Amount (currency)
	quote.Amount = parseFloat(fields[9])

	// Update time
	date := fields[30]
	time := fields[31]
	quote.UpdateTime = fmt.Sprintf("%s %s", date, time)

	return quote, nil
}

func parseFloat(s string) float64 {
	s = strings.TrimSpace(s)
	if s == "-" || s == "" {
		return 0
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f
}

func parseInt64(s string) int64 {
	s = strings.TrimSpace(s)
	if s == "-" || s == "" {
		return 0
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return int64(f)
}

// gbkToUtf8 converts GBK encoded string to UTF-8
func gbkToUtf8(s string) string {
	decoder := simplifiedchinese.GBK.NewDecoder()
	result, _, err := transform.String(decoder, s)
	if err != nil {
		return s
	}
	return result
}
