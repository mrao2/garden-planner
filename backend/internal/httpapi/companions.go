package httpapi

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strings"
)

type CompanionRecord struct {
	CommonName     string   `json:"common_name"`
	ScientificName string   `json:"scientific_name"`
	Helps          string   `json:"helps"`
	HelpedBy       string   `json:"helped_by"`
	Attracts       string   `json:"attracts"`
	Repels         string   `json:"repels"`
	Avoid          string   `json:"avoid"`
	Comments       string   `json:"comments"`
	HelpsList      []string `json:"helps_list"`
	HelpedByList   []string `json:"helped_by_list"`
	AttractsList   []string `json:"attracts_list"`
	RepelsList     []string `json:"repels_list"`
	AvoidList      []string `json:"avoid_list"`
}

type CompanionIndex struct {
	records []CompanionRecord
	byName  map[string]CompanionRecord
}

func LoadCompanionData(path string) (*CompanionIndex, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}
	if len(headers) == 0 {
		return nil, errors.New("companion csv has no headers")
	}

	index := &CompanionIndex{
		records: make([]CompanionRecord, 0),
		byName:  make(map[string]CompanionRecord),
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		row := make(map[string]string)
		for i, header := range headers {
			if i < len(record) {
				row[header] = strings.TrimSpace(record[i])
			}
		}

		item := CompanionRecord{
			CommonName:     row["Common name"],
			ScientificName: row["Scientific name"],
			Helps:          row["Helps"],
			HelpedBy:       row["Helped by"],
			Attracts:       row["Attracts"],
			Repels:         row["-Repels/+distracts"],
			Avoid:          row["Avoid"],
			Comments:       row["Comments"],
		}
		item.HelpsList = splitList(item.Helps)
		item.HelpedByList = splitList(item.HelpedBy)
		item.AttractsList = splitList(item.Attracts)
		item.RepelsList = splitList(item.Repels)
		item.AvoidList = splitList(item.Avoid)

		index.records = append(index.records, item)
		nameKey := normalizeForMatch(item.CommonName)
		if nameKey != "" {
			index.byName[nameKey] = item
		}
	}

	return index, nil
}

func (c *CompanionIndex) FindByName(name string) (CompanionRecord, bool) {
	if c == nil {
		return CompanionRecord{}, false
	}

	query := normalizeForMatch(name)
	if query == "" {
		return CompanionRecord{}, false
	}

	if record, ok := c.byName[query]; ok {
		return record, true
	}

	if strings.HasSuffix(query, "s") {
		if record, ok := c.byName[strings.TrimSuffix(query, "s")]; ok {
			return record, true
		}
	}
	if record, ok := c.byName[query+"s"]; ok {
		return record, true
	}
	if record, ok := c.byName[query+"es"]; ok {
		return record, true
	}

	wrappedQuery := " " + query + " "
	for key, record := range c.byName {
		if strings.Contains(" "+key+" ", wrappedQuery) {
			return record, true
		}
	}

	return CompanionRecord{}, false
}

func splitList(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	items := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		items = append(items, trimmed)
	}
	return items
}
