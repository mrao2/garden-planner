package httpapi

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type seedPacketScanResponse struct {
	Text            string             `json:"text"`
	Name            string             `json:"name"`
	Variety         string             `json:"variety"`
	Season          string             `json:"season"`
	Spacing         string             `json:"spacing"`
	Notes           string             `json:"notes"`
	DaysToMaturity  *int               `json:"days_to_maturity"`
	SowDepth        string             `json:"sow_depth"`
	GerminationTime string             `json:"germination_time"`
	WhenToSow       string             `json:"when_to_sow"`
	DaysToEmerge    string             `json:"days_to_emerge"`
	SeedSpacing     string             `json:"seed_spacing"`
	RowSpacing      string             `json:"row_spacing"`
	Thinning        string             `json:"thinning"`
	Boxes           map[string]scanBox `json:"boxes,omitempty"`
}

type scanBox struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	W float64 `json:"w"`
	H float64 `json:"h"`
}

func (h *Handlers) scanSeedPacket(w http.ResponseWriter, r *http.Request) {
	ocrEndpoint := strings.TrimSpace(os.Getenv("OCR_ENDPOINT"))
	if ocrEndpoint == "" {
		respondError(w, http.StatusServiceUnavailable, "not_ready", "OCR service not configured")
		return
	}

	if err := r.ParseMultipartForm(12 << 20); err != nil {
		respondError(w, http.StatusBadRequest, "invalid_form", "Unable to read upload")
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		respondError(w, http.StatusBadRequest, "missing_image", "Image is required")
		return
	}
	defer file.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("image", header.Filename)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "upload_failed", "Unable to process image")
		return
	}
	if _, err := io.Copy(part, file); err != nil {
		respondError(w, http.StatusInternalServerError, "upload_failed", "Unable to process image")
		return
	}
	company := strings.TrimSpace(r.FormValue("company"))
	if company != "" {
		_ = writer.WriteField("company", company)
	}

	if err := writer.Close(); err != nil {
		respondError(w, http.StatusInternalServerError, "upload_failed", "Unable to process image")
		return
	}

	req, err := http.NewRequestWithContext(r.Context(), http.MethodPost, ocrEndpoint, &body)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "ocr_failed", "Unable to reach OCR service")
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		respondError(w, http.StatusBadGateway, "ocr_failed", "OCR service unavailable")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respondError(w, http.StatusBadGateway, "ocr_failed", "OCR service error")
		return
	}

	var payload struct {
		Text           string `json:"text"`
		TextSpecs      string `json:"text_specs"`
		TextSow        string `json:"text_sow"`
		TextName       string `json:"text_name"`
		Name           string `json:"name"`
		Variety        string `json:"variety"`
		DaysToEmerge   string `json:"days_to_emerge"`
		SowDepth       string `json:"sow_depth"`
		SeedSpacing    string `json:"seed_spacing"`
		RowSpacing     string `json:"row_spacing"`
		Thinning       string `json:"thinning"`
		DaysToMaturity string `json:"days_to_maturity"`
		WhenToSow      string `json:"when_to_sow"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		respondError(w, http.StatusBadGateway, "ocr_failed", "Invalid OCR response")
		return
	}

	result := parseSeedPacketText(payload.TextName, payload.TextSpecs, payload.TextSow, payload.Text)
	if payload.Name != "" {
		result.Name = payload.Name
	}
	if payload.Variety != "" {
		result.Variety = payload.Variety
	}
	if payload.DaysToEmerge != "" {
		result.DaysToEmerge = payload.DaysToEmerge
	}
	if payload.SowDepth != "" {
		result.SowDepth = payload.SowDepth
	}
	if payload.SeedSpacing != "" {
		result.SeedSpacing = payload.SeedSpacing
	}
	if payload.RowSpacing != "" {
		result.RowSpacing = payload.RowSpacing
	}
	if payload.Thinning != "" {
		result.Thinning = payload.Thinning
	}
	if payload.DaysToMaturity != "" && result.DaysToMaturity == nil {
		if value := extractFirstInt(payload.DaysToMaturity); value != nil {
			result.DaysToMaturity = value
		}
	}
	if payload.WhenToSow != "" {
		result.WhenToSow = payload.WhenToSow
	}
	if payload.TextSpecs != "" || payload.TextSow != "" || payload.TextName != "" {
		result.Boxes = seedPacketBoxes()
	}
	respondJSON(w, http.StatusOK, dataEnvelope[seedPacketScanResponse]{Data: result})
}

func parseSeedPacketText(textName string, textSpecs string, textSow string, textFull string) seedPacketScanResponse {
	combined := strings.TrimSpace(strings.Join([]string{textName, textSpecs, textSow, textFull}, "\n"))
	response := seedPacketScanResponse{
		Text:            combined,
		Season:          "",
		Spacing:         "",
		Notes:           "",
		SowDepth:        "",
		GerminationTime: "",
	}

	lines := make([]string, 0)
	for _, line := range strings.Split(textSpecs, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		lines = append(lines, trimmed)
	}

	nameRe := regexp.MustCompile(`(?i)^name\\s*[:\\-]\\s*(.+)$`)
	varietyRe := regexp.MustCompile(`(?i)^variety\\s*[:\\-]\\s*(.+)$`)
	seasonRe := regexp.MustCompile(`(?i)season\\s*[:\\-]\\s*(.+)$`)
	spacingRe := regexp.MustCompile(`(?i)spacing\\s*[:\\-]\\s*(.+)$`)
	dtmRe := regexp.MustCompile(`(?i)(days\\s*to\\s*maturity|maturity|dtm)[^0-9]*(\\d{2,3})`)
	sowDepthRe := regexp.MustCompile(`(?i)(sow\\s*depth|plant\\s*depth|depth)\\s*[:\\-]?\\s*([0-9./ \\-]+\\s*(?:in|inch|inches|cm|mm)?)`)
	germinationRe := regexp.MustCompile(`(?i)(germination|germinates|sprout)\\s*[:\\-]?\\s*([0-9\\-to ]+\\s*(?:days|day))`)
	whenToSowRe := regexp.MustCompile(`(?i)(when\\s*to\\s*sow[^\\n]*\\n?(?:[^\\n]+)?)`)
	daysToEmergeRe := regexp.MustCompile(`(?i)(days\\s*to\\s*emerge)[^0-9]*([0-9\\-to ]+\\s*(?:days|day)?)`)
	seedSpacingRe := regexp.MustCompile(`(?i)(seed\\s*spacing)\\s*[:\\-]?\\s*([-0-9/\"' ]+(?:to\\s*[-0-9/\"' ]+)?)`)
	rowSpacingRe := regexp.MustCompile(`(?i)(row\\s*spacing)\\s*[:\\-]?\\s*([-0-9/\"' ]+(?:to\\s*[-0-9/\"' ]+)?)`)
	thinningRe := regexp.MustCompile(`(?i)(thinning)[^a-z0-9]*([^\\n]+)`)

	if textName != "" {
		linesName := strings.Fields(textName)
		if len(linesName) > 0 {
			response.Name = strings.Join(linesName, " ")
		}
	}

	for _, line := range lines {
		if response.Name == "" {
			if match := nameRe.FindStringSubmatch(line); len(match) > 1 {
				response.Name = strings.TrimSpace(match[1])
				continue
			}
		}
		if response.Variety == "" {
			if match := varietyRe.FindStringSubmatch(line); len(match) > 1 {
				response.Variety = strings.TrimSpace(match[1])
				continue
			}
		}
		if response.Season == "" {
			if match := seasonRe.FindStringSubmatch(line); len(match) > 1 {
				response.Season = strings.TrimSpace(match[1])
				continue
			}
		}
		if response.Spacing == "" {
			if match := spacingRe.FindStringSubmatch(line); len(match) > 1 {
				response.Spacing = strings.TrimSpace(match[1])
				continue
			}
		}
	}

	if response.Name == "" && len(lines) > 0 {
		response.Name = lines[0]
	}

	if match := dtmRe.FindStringSubmatch(textSpecs); len(match) > 2 {
		if value, err := strconv.Atoi(match[2]); err == nil {
			response.DaysToMaturity = &value
		}
	}
	if match := sowDepthRe.FindStringSubmatch(textSpecs); len(match) > 2 {
		response.SowDepth = strings.TrimSpace(match[2])
	}
	if match := germinationRe.FindStringSubmatch(textSpecs); len(match) > 2 {
		response.GerminationTime = strings.TrimSpace(match[2])
	}
	if match := daysToEmergeRe.FindStringSubmatch(textSpecs); len(match) > 2 {
		response.DaysToEmerge = strings.TrimSpace(match[2])
	}
	if match := seedSpacingRe.FindStringSubmatch(textSpecs); len(match) > 2 {
		response.SeedSpacing = strings.TrimSpace(match[2])
	}
	if match := rowSpacingRe.FindStringSubmatch(textSpecs); len(match) > 2 {
		response.RowSpacing = strings.TrimSpace(match[2])
	}
	if match := thinningRe.FindStringSubmatch(textSpecs); len(match) > 2 {
		response.Thinning = strings.TrimSpace(match[2])
	}
	if match := whenToSowRe.FindStringSubmatch(textSow); len(match) > 1 {
		response.WhenToSow = cleanInline(match[1])
	} else if textSow != "" {
		response.WhenToSow = cleanInline(textSow)
	}

	return response
}

func seedPacketBoxes() map[string]scanBox {
	// Botanical Interests relative boxes (approximate).
	return map[string]scanBox{
		"name":             {X: 0.099, Y: 0.091, W: 0.25, H: 0.103},
		"variety":          {X: 0.099, Y: 0.188, W: 0.267, H: 0.099},
		"emerge":           {X: 0.132, Y: 0.558, W: 0.235, H: 0.036},
		"depth":            {X: 0.135, Y: 0.611, W: 0.235, H: 0.039},
		"seed_spacing":     {X: 0.135, Y: 0.674, W: 0.239, H: 0.06},
		"row_spacing":      {X: 0.137, Y: 0.753, W: 0.241, H: 0.037},
		"thinning":         {X: 0.135, Y: 0.812, W: 0.242, H: 0.062},
		"days_to_maturity": {X: 0.137, Y: 0.893, W: 0.246, H: 0.039},
		"sow":              {X: 0.375, Y: 0.505, W: 0.598, H: 0.283},
	}
}

func cleanInline(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return ""
	}
	lines := strings.Fields(trimmed)
	return strings.Join(lines, " ")
}

func extractFirstInt(value string) *int {
	for _, token := range strings.Fields(value) {
		cleaned := strings.Trim(token, "\"'()")
		cleaned = strings.TrimSpace(cleaned)
		cleaned = strings.Trim(cleaned, "-")
		if cleaned == "" {
			continue
		}
		if number, err := strconv.Atoi(cleaned); err == nil {
			return &number
		}
	}
	digits := regexp.MustCompile(`\\d+`).FindString(value)
	if digits == "" {
		return nil
	}
	if number, err := strconv.Atoi(digits); err == nil {
		return &number
	}
	return nil
}
