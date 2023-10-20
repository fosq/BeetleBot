package tft

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var (
	activePatchNotes    string
	patchNotes          string
	currentPatchVersion string
)

type PatchNote struct {
	ID           int    `json:"id"`
	PatchVersion string `json:"patchVersion"`
	RegisteredAt int64  `json:"registeredAt"`
	Season       string `json:"season"`
}

func createPatchNotes() bool {
	var (
		titleContentMap = make(map[string][]string, 0)
		headings        []string
		setVersion      string
	)

	count := -1
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Create and modify HTTP request before sending
	request, err := http.NewRequest("GET", "https://lolchess.gg/guide/patch-notes", nil)
	if err != nil {
		fmt.Println(err)
		return false
	}
	request.Header.Set("Accept-Language", "en")
	request.Header.Set("User-Agent", "BeetleBot - TFT Patch notes scraper for a private Discord server")

	// Make request
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		fmt.Println(err)
		return false
	}

	// Find patch notes' headings and contents from <section>
	document.Find("section[class*='ep0afc42']").Each(func(index int, items *goquery.Selection) {
		items.Find("div[class='css-1qhzwuq e18ae7l60']").Each(func(index int, item *goquery.Selection) {

			// Headings
			item.Find("div[class='css-kvkqpd e18ae7l61']").Each(func(index int, element *goquery.Selection) {
				content := strings.TrimSpace(element.Text())
				headings = append(headings, content)
			})

			// Headings' content (list)
			item.Find("li[class='css-ouvf7b e18ae7l65']").Each(func(index int, element *goquery.Selection) {
				if index == 0 {
					count++
				}
				titleContentMap[headings[count]] = append(titleContentMap[headings[count]], element.Text())
			})
		})
	})

	var jsonData string
	document.Find("script").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Text(), "patchNotes") {
			jsonData = s.Text()
		}
	})

	// Regex for extracting JSON data
	r := regexp.MustCompile(`patchNotes":(\[.*?\])`)
	matches := r.FindStringSubmatch(jsonData)
	if len(matches) < 2 {
		log.Fatal("Could not find JSON data")
	}
	extractedJSON := matches[1]

	// Parse the extracted JSON data
	var patchNotesObj []PatchNote
	err = json.Unmarshal([]byte(extractedJSON), &patchNotesObj)
	if err != nil {
		log.Fatal(err)
	}

	latestPatchNote := patchNotesObj[0]
	currentPatchVersion = latestPatchNote.PatchVersion

	//// Write formatted patch notes to file
	// Write set version
	parts := regexp.MustCompile(`([a-zA-Z]+)(\d+\.\d+)`).FindStringSubmatch(latestPatchNote.Season)
	setVersion = strings.ToUpper(strings.Join(parts[1:], " "))

	patchNotes += ("# " + setVersion + "\n")

	// Write patch date and version
	fullDate := time.Unix(latestPatchNote.RegisteredAt/1000, 0)
	printDate := fullDate.Month().String()[:3] + " " + strconv.Itoa(fullDate.Day())

	patchNotes += ("## " + printDate + " - __" + latestPatchNote.PatchVersion + "__\n\n")

	for i := range headings {
		patchNotes += ("__**" + headings[i] + "**__\n")

		// Write headings' content
		for j := range titleContentMap[headings[i]] {
			patchNotes += ("- " + titleContentMap[headings[i]][j] + "\n")
		}
		patchNotes += ("\n")
	}

	patchNotes = strings.TrimRight(patchNotes, "\n")

	return true
}

func comparePatchNotes() bool {
	if activePatchNotes == patchNotes {
		patchNotes = ""
		return false
	}

	activePatchNotes = patchNotes
	patchNotes = ""
	fmt.Printf("Patches updated to version %v\n", currentPatchVersion)

	return true
}

// Compares previously saved patch notes from newly fetched ones, returns true if change was found
func CheckForUpdate() bool {
	if !createPatchNotes() {
		return false
	}
	fmt.Printf("%v - Current patch version: %v\n", time.Now().Format(time.RFC3339), currentPatchVersion)

	return comparePatchNotes()
}

func GetPatchNotes() string {
	return activePatchNotes
}
