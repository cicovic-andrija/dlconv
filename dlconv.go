package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"text/template"
	"time"
)

const (
	TemplateName         = "dive-log"
	OutputFile           = TemplateName + ".md"
	TemplateMarkdownText = `---
tags: [diving, my]
title: {{ .Title }}
comment: 'This document is auto-generated from the master dive log using https://github.com/cicovic-andrija/dlconv.'
created: '2023-04-16T20:36:59.779Z'
modified: '{{ .ModifiedTimeUTC }}'
---

# {{ .Title }}

## Index
{{ range .Dives }}
- **[No. {{ .Cardinal }}: {{ .Site }}, {{ .Date }}.](#no-{{ .Cardinal }})**{{ end }}

## Dive Data
{{ range .Dives }}
### <a id="no-{{ .Cardinal }}"></a>No. {{ .Cardinal }}: {{ .Site }}, {{ .Date }}.

| Parameter | Value |
| --------- | ----- |
| Time in | {{ .Time }} |
| Duration | {{ .Duration }} min |
| Max. depth | {{ .MaxDepth }} m |
| Avg. depth | {{ .AvgDepth }} m |
| Tank pressure | {{ .TankPressureStart }} bar - {{ .TankPressureEnd }} bar |
| Gas | {{ .Gas }} {{ .O2 }} % |
| Decompression | {{ .DecompressionDive }} |
| CNS | {{ .CNS }} % |
| Altitude | {{ .Altitude }} m |
| Entry | {{ .From }} |
| Operator | {{ .Operator }} |
| Suit | {{ .SuitType }}{{ if ne .SuitType "Dry suit" }} {{ .SuitThickness }} mm{{ end }} |
| Weights | {{ .Weights }} kg |
| Tank | {{ .TankType }} {{ .TankVolume }} litres |
| Computer | {{ .Computer }} |
| Weather | {{ .Weather }} |
| Air temp. | {{ .AirTemp }} °C |
| Water | {{ .WaterType }} |
| Water temp. | {{ .WaterMinTemp }} °C |
| Visibility | {{ .WaterVisibility }} |
| Drift | {{ .DriftDive }} |
{{ end }}
`
)

type TemplateData struct {
	Title           string
	ModifiedTimeUTC string
	Dives           []DiveData
}

type DiveData struct {
	Cardinal          string
	Site              string
	Date              string
	Time              string
	Duration          string
	MaxDepth          string
	AvgDepth          string
	TankPressureStart string
	TankPressureEnd   string
	DecompressionDive string
	Gas               string
	O2                string
	CNS               string
	Altitude          string
	From              string
	Operator          string
	SuitType          string
	SuitThickness     string
	Weights           string
	TankType          string
	TankVolume        string
	Computer          string
	DecoAlgPFact      string
	Weather           string
	AirTemp           string
	WaterType         string
	WaterMinTemp      string
	WaterVisibility   string
	DriftDive         string
}

// Convert dive data in CSV to Markdown.
func main() {
	var inputFile string
	flag.Parse()
	if inputFile = flag.Arg(0); inputFile == "" {
		panic("input file not provided")
	}

	fh, err := os.Open(inputFile)
	if err != nil {
		panic(fmt.Errorf("failed to open %s: %v", inputFile, err))
	}
	defer fh.Close()

	csvreader := csv.NewReader(fh)
	records, err := csvreader.ReadAll()
	if err != nil {
		panic(fmt.Errorf("failed to read CSV data from %s: %v", inputFile, err))
	}

	dives := make([]DiveData, 0)

	// Note: This is a custom conversion program, it assumes that input file is
	// valid in every sense of that word.
	i := 0
	for i < len(records) {
		increment := 1

		if records[i][4] == "Site" && records[i][9] == "Date" {
			dives = append(dives, DiveData{
				Cardinal:          records[i][1],
				Site:              records[i+1][4],
				Date:              records[i+1][9],
				Time:              records[i+1][12],
				Duration:          records[i+4][4],
				SuitType:          records[i+4][11],
				MaxDepth:          records[i+5][4],
				SuitThickness:     records[i+5][11],
				AvgDepth:          records[i+6][4],
				Weights:           records[i+6][11],
				TankPressureStart: records[i+7][4],
				TankType:          records[i+7][11],
				TankPressureEnd:   records[i+8][4],
				TankVolume:        records[i+8][11],
				DecompressionDive: records[i+9][4],
				Computer:          records[i+9][11],
				Gas:               records[i+10][4],
				DecoAlgPFact:      records[i+10][11],
				O2:                records[i+11][4],
				CNS:               records[i+12][4],
				Altitude:          records[i+13][4],
				Weather:           records[i+13][11],
				From:              records[i+14][4],
				AirTemp:           records[i+14][11],
				Operator:          records[i+17][1],
				WaterType:         records[i+17][11],
				WaterMinTemp:      records[i+18][11],
				WaterVisibility:   records[i+19][11],
				DriftDive:         records[i+20][11],
			})

			increment = 23
		}

		i += increment
	}

	for i, j := 0, len(dives)-1; i < j; i, j = i+1, j-1 {
		dives[i], dives[j] = dives[j], dives[i]
	}

	t := template.Must(template.New(TemplateName).Parse(TemplateMarkdownText))
	outFd, err := os.Create(OutputFile)
	if err != nil {
		panic(fmt.Errorf("failed to open output file %s for writing: %v", OutputFile, err))
	}
	data := &TemplateData{
		Title:           "Dive Log",
		ModifiedTimeUTC: time.Now().UTC().Format(time.RFC3339),
		Dives:           dives,
	}
	t.Execute(outFd, data)
}
