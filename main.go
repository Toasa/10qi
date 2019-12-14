package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Reporter struct {
	url   string
	today Weather
	tomrw Weather
}

type Weather struct {
	date      string
	weather   string
	high_temp int
	high_diff string
	low_temp  int
	low_diff  string
	rainyPct  [4]string // 0-6, 6-12, 12-18, 18-24
}

func main() {
	url := "https://weather.yahoo.co.jp/weather/jp/26/6110.html"
	r := Reporter{url: url}

	r.setDate()
	r.scrape()
	r.report()
}

func (r *Reporter) setDate() {
	now := time.Now()
	tomrw := now.AddDate(0, 0, 1)
	nowStr := fmt.Sprintf("%d-%d-%d", now.Year(), now.Month(), now.Day())
	tomrwStr := fmt.Sprintf("%d-%d-%d", tomrw.Year(), tomrw.Month(), tomrw.Day())
	r.today.date = nowStr
	r.tomrw.date = tomrwStr
}

func (r *Reporter) scrape() {
	selection := r.getSelection()
	tableRow := selection.Children().Children().Children()
	todayNode := tableRow.Children()
	tomrwNode := todayNode.Next()

	r.today.parse(todayNode)
	r.tomrw.parse(tomrwNode)
}

func (w *Weather) parse(node *goquery.Selection) {
	weather := node.Children().Children().Next()
	html, _ := weather.Html()
	head := strings.Index(html, ">")
	w.weather = html[head+1:]

	temp := weather.Next()
	highTemp := temp.Children()
	highDiffStr, _ := highTemp.Html()
	head = strings.Index(highDiffStr, "[")
	w.high_diff = highDiffStr[head:]
	highTempStr, _ := highTemp.Children().Html()
	w.high_temp, _ = strconv.Atoi(highTempStr)

	lowTemp := highTemp.Next()
	lowDiffStr, _ := lowTemp.Html()
	head = strings.Index(lowDiffStr, "[")
	w.low_diff = lowDiffStr[head:]
	lowTempStr, _ := lowTemp.Children().Html()
	w.low_temp, _ = strconv.Atoi(lowTempStr)

	r0_6 := temp.Next().Children().Children().Next().Children().Next()
	r6_12 := r0_6.Next()
	r12_18 := r6_12.Next()
	r18_24 := r12_18.Next()

	w.rainyPct[0], _ = r0_6.Html()
	w.rainyPct[1], _ = r6_12.Html()
	w.rainyPct[2], _ = r12_18.Html()
	w.rainyPct[3], _ = r18_24.Html()
}

func (r Reporter) getSelection() *goquery.Selection {
	doc, err := goquery.NewDocument(r.url)
	if err != nil {
		panic(err)
	}

	return doc.Find("div.forecastCity")
}

func (r Reporter) report() {
	fmt.Printf("\n%s\n", r.today.date)
	fmt.Printf("+--------------------+\n")
	r.today.report()
	fmt.Printf("+--------------------+\n")
	fmt.Printf("\n%s\n", r.tomrw.date)
	fmt.Printf("+--------------------+\n")
	r.tomrw.report()
	fmt.Printf("+--------------------+\n\n")
}

func (w Weather) report() {
	fmt.Printf("|  10qi: %s\n", w.weather)
	fmt.Printf("|  max : %d.C %s\n", w.high_temp, w.high_diff)
	fmt.Printf("|  min : %d.C %s\n", w.low_temp, w.low_diff)
	fmt.Printf("|  rain:\n")

	fmt.Printf("|    0-6  : %s\n", w.rainyPct[0])
	fmt.Printf("|    6-12 : %s\n", w.rainyPct[1])
	fmt.Printf("|    12-18: %s\n", w.rainyPct[2])
	fmt.Printf("|    18-24: %s\n", w.rainyPct[3])
}
