package tft

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var (
	patchNotes   []PatchNote // Latest 5 patches with parsed messages
	newPatchNote PatchNote   // Latest fetched patch note for comparison
	AllPatchInfo []PatchNote // All patches without parsed message
)

type PatchNote struct {
	ID           int    `json:"id"`
	PatchVersion string `json:"patchVersion"`
	RegisteredAt int64  `json:"registeredAt"`
	Season       string `json:"season"`
	Message      string `json:"message"`
}

func fetchPatchNotes(id int) (*http.Response, error) {
	var patchId string

	if id != 0 {
		patchId = strconv.Itoa(id)
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Create and modify HTTP request before sending
	request, err := http.NewRequest("GET", "https://lolchess.gg/guide/patch-notes/"+patchId, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept-Language", "en")
	request.Header.Set("User-Agent", "BeetleBot - TFT Patch notes scraper for a private Discord server")

	// Make request
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func parsePatchNotes(response *http.Response) error {
	var (
		titleContentMap = make(map[string][]string, 0)
		headings        []string
		setVersion      string
	)
	count := -1

	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return err
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
		return fmt.Errorf("could not find JSON data")
	}
	extractedJSON := matches[1]

	// Parse the extracted JSON data
	err = json.Unmarshal([]byte(extractedJSON), &AllPatchInfo)
	if err != nil {
		return err
	}

	newPatchNote = AllPatchInfo[0]

	//// Write formatted patch notes to file
	// Write set version
	parts := regexp.MustCompile(`([a-zA-Z]+)(\d+\.\d+)`).FindStringSubmatch(newPatchNote.Season)
	setVersion = strings.ToUpper(strings.Join(parts[1:], " "))

	newPatchNote.Message += ("# " + setVersion + "\n")

	// Write patch date and version
	fullDate := time.Unix(newPatchNote.RegisteredAt/1000, 0)
	printDate := fullDate.Month().String()[:3] + " " + strconv.Itoa(fullDate.Day())

	newPatchNote.Message += ("## " + printDate + " - __" + newPatchNote.PatchVersion + "__\n\n")

	for i := range headings {
		newPatchNote.Message += ("__**" + headings[i] + "**__\n")

		// Write headings' content
		for j := range titleContentMap[headings[i]] {
			newPatchNote.Message += ("- " + titleContentMap[headings[i]][j] + "\n")
		}
		newPatchNote.Message += ("\n")
	}

	newPatchNote.Message = strings.TrimRight(newPatchNote.Message, "\n")

	return nil
}

// Returns true if latest patchNotes patch version differs from newly fetched patch version
func comparePatchNotes() bool {
	return patchNotes[len(patchNotes)-1].PatchVersion != newPatchNote.PatchVersion
}

// Compares previously saved patch notes from newly fetched ones, returns true if change was found
func UpdatePatches() bool {
	response, err := fetchPatchNotes(0)
	if err != nil {
		fmt.Println(err)
		return false
	}

	if err := parsePatchNotes(response); err != nil {
		fmt.Println(err)
		return false
	}

	if len(patchNotes) == 0 {
		fmt.Println("Initializing: Setting up patch notes list")
		patchNotes = AllPatchInfo[:5] // Fetch 5 latest patch notes
		sortPatchNotes()
		patchNotes[len(patchNotes)-1].Message = newPatchNote.Message
		return true
	}

	if comparePatchNotes() {
		patchNotes = append(patchNotes, newPatchNote)
		fmt.Printf("Patches updated to version %v\n", patchNotes[len(patchNotes)-1].PatchVersion)
		return true
	}

	fmt.Printf("%v - Current patch version: %v\n", time.Now().Format(time.RFC3339), patchNotes[len(patchNotes)-1].PatchVersion)
	return false
}

// Sorts patchNotes by ID ascendingly
func sortPatchNotes() {
	sort.Slice(patchNotes, func(i, j int) bool {
		return patchNotes[i].ID < patchNotes[j].ID
	})
}

func GetPatchNotes() string {
	return patchNotes[len(patchNotes)-1].Message
}
