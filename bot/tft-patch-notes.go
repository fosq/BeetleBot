package tft

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var (
	titleContentMap = make(map[string][]string, 0)
	headings        []string
)

func CreatePatchNotes() {
	count := -1
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Create and modify HTTP request before sending
	request, err := http.NewRequest("GET", "https://lolchess.gg/guide/patch-notes", nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Accept-Language", "en")
	request.Header.Set("User-Agent", "BeetleBot - TFT Patch notes scraper for a private Discord server")

	// Make request
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}

	// Find patch notes' headings and contents from <section>
	document.Find("section[class*='ep0afc42']").Each(func(index int, items *goquery.Selection) {
		items.Find("div[class='css-1qhzwuq e18ae7l60']").Each(func(index int, item *goquery.Selection) {
			item.Find("div[class='css-kvkqpd e18ae7l61']").Each(func(index int, element *goquery.Selection) {
				content := strings.TrimSpace(element.Text())
				headings = append(headings, content)
			})

			item.Find("li[class='css-ouvf7b e18ae7l65']").Each(func(index int, element *goquery.Selection) {
				if index == 0 {
					count++
				}
				titleContentMap[headings[count]] = append(titleContentMap[headings[count]], element.Text())
			})
		})
		fmt.Println()
	})

	// Create patch notes file
	f, err := os.Create("patchnotes.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := range headings {
		_, err := f.WriteString("__**" + headings[i] + "**__\n")
		checkErr(err)
		for j := range titleContentMap[headings[i]] {
			_, err := f.WriteString("- " + titleContentMap[headings[i]][j] + "\n")
			checkErr(err)
		}
		_, err = f.WriteString("\n")
		checkErr(err)
	}
}

func ComparePatchNotes() bool {
	newData, err := os.ReadFile("patchnotes.txt")
	checkErr(err)
	oldData, err := os.ReadFile("activepatchnotes.txt")
	checkErr(err)

	if bytes.Equal(newData, oldData) {
		err := os.Remove("patchnotes.txt")
		checkErr(err)
		return false
	} else {
		err := os.Rename("patchnotes.txt", "activepatchnotes.txt")
		checkErr(err)
		return true
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
