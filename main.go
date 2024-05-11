package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func main() {
	inputDir := os.Args[1]
	rules := parseRules()
	instList := keys(rules)
	pathList := fetchFiles(inputDir, instList)
	for iname, path := range pathList {
		f, _ := os.Open(path)
		fo, _ := os.Create(fmt.Sprintf("%s.csv", iname))
		csvReader := csv.NewReader(f)
		csvWriter := csv.NewWriter(fo)
		recs, _ := csvReader.ReadAll()
		irules := rules[iname]
		for _, rec := range recs {
			ts, _ := strconv.ParseInt(rec[0], 10, 64)
			o, _ := strconv.ParseFloat(rec[1], 64)
			h, _ := strconv.ParseFloat(rec[2], 64)
			l, _ := strconv.ParseFloat(rec[3], 64)
			c, _ := strconv.ParseFloat(rec[4], 64)
			mx := getMultiplier(irules, ts)
			csvWriter.Write([]string{
				rec[0],
				fmt.Sprint(o * mx),
				fmt.Sprint(h * mx),
				fmt.Sprint(l * mx),
				fmt.Sprint(c * mx),
				rec[5],
			})
		}
		f.Close()
		fo.Close()
	}
}

func fetchFiles(srcDir string, instList []string) map[string]string {
	paths := make(map[string]string)
	filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err == nil && in(instList, info.Name()[:len(info.Name())-7]) {
			paths[info.Name()[:len(info.Name())-7]] = path
		}
		return nil
	})
	fmt.Println(paths)
	return paths
}

func getMultiplier(rules []Rule, ts int64) float64 {
	for _, rule := range rules {
		if ts >= rule.From && ts <= rule.To {
			return rule.multiplier
		} else {
			continue
		}
	}
	return 1.0
}

func in[T comparable](arr []T, elem T) bool {
	for _, e := range arr {
		if e == elem {
			return true
		}
	}
	return false
}

func keys[T any](in map[string]T) []string {
	var keys []string
	for k, _ := range in {
		keys = append(keys, k)
	}
	return keys
}

type Rule struct {
	From       int64
	To         int64
	multiplier float64
}

func parseRules() map[string][]Rule {
	ruleFile, err := os.Open("rules.csv")
	if err != nil {
		log.Fatal(err)
	}
	csvReader := csv.NewReader(ruleFile)
	csvReader.Read()
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	rules := make(map[string][]Rule)
	for _, record := range records {
		from, _ := time.Parse("02-01-2006", record[1])
		to, _ := time.Parse("02-01-2006 15:04:05", fmt.Sprintf("%s %s", record[2], "23:59:59"))
		multi, _ := strconv.ParseFloat(record[3], 64)
		record[0] = strings.ToUpper(record[0])
		instrule, ok := rules[record[0]]
		if ok {
			instrule = append(instrule, Rule{
				From:       from.Unix(),
				To:         to.Unix(),
				multiplier: multi,
			})
		} else {
			rules[record[0]] = []Rule{
				{
					From:       from.Unix(),
					To:         to.Unix(),
					multiplier: multi,
				},
			}
		}
	}
	fmt.Println(rules)
	return rules
}
