package tft

import (
	"discordbot/logs"
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

func parsePatchNotes(response *http.Response) (PatchNote, error) {
	var newPatchNote PatchNote

	var (
		titleContentMap = make(map[string][]string, 0)
		headings        []string
		setVersion      string
	)
	count := -1

	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return PatchNote{}, err
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
		return PatchNote{},
			fmt.Errorf("could not find any matches for all patch notes from lolchess.gg, re-running in 30 minutes")
	}
	extractedJSON := matches[1]

	// Parse the extracted JSON data
	err = json.Unmarshal([]byte(extractedJSON), &AllPatchInfo)
	if err != nil {
		return PatchNote{}, err
	}
	AllPatchInfo = sortPatchNotes(AllPatchInfo)
	newPatchNote = AllPatchInfo[len(AllPatchInfo)-1]

	//// Write formatted patch notes to file
	// Write set version
	parts := regexp.MustCompile(`([a-zA-Z]+)((\d+\.\d+)|(\d+))`).FindStringSubmatch(newPatchNote.Season)
	setVersion = strings.ToUpper(strings.Join(parts[1:3], " "))

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

	return newPatchNote, nil
}

// Returns true if latest patchNotes patch version differs from newly fetched patch version
func comparePatchNotes(newPatchNote PatchNote) bool {
	return patchNotes[len(patchNotes)-1].PatchVersion != newPatchNote.PatchVersion
}

// Compares previously saved patch notes from newly fetched ones, returns true if change was found
func UpdatePatches() bool {
	response, err := fetchPatchNotes(0)
	logs.Check(err)

	newPatchNote, err := parsePatchNotes(response)
	if !logs.Check(err) || (newPatchNote == PatchNote{}) {
		return false
	}

	if len(patchNotes) == 0 {
		logs.WriteLogFile("Initializing: Setting up patch notes list")
		patchNotes = AllPatchInfo[len(AllPatchInfo)-5:] // Fetch 5 latest patch notes
		logs.WriteLogFile(fmt.Sprintf("Initializing: Added patch notes from version %v to %v",
			patchNotes[0].PatchVersion, patchNotes[len(patchNotes)-1].PatchVersion))
		patchNotes[len(patchNotes)-1].Message = newPatchNote.Message
		return true
	}

	if comparePatchNotes(newPatchNote) {
		patchNotes = append(patchNotes, newPatchNote)
		logs.WriteLogFile(fmt.Sprintf("Patches updated to version %v",
			patchNotes[len(patchNotes)-1].PatchVersion))
		return true
	}

	logs.WriteLogFile(fmt.Sprintf("%v - Current patch version: %v",
		time.Now().Format(time.RFC3339), patchNotes[len(patchNotes)-1].PatchVersion))
	return false
}

// Sorts patchNotes by ID ascendingly
func sortPatchNotes(notes []PatchNote) []PatchNote {
	sort.Slice(notes, func(i, j int) bool {
		return notes[i].ID < notes[j].ID
	})
	return notes
}

func GetPatchNotes() string {
	return patchNotes[len(patchNotes)-1].Message
}
