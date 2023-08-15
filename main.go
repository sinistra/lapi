package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	argsWithoutProg := os.Args[1:]

	formattedAddress := strings.Join(argsWithoutProg, " ")
	fmt.Println(formattedAddress)
	nbnLocations, err := GetNBNSuggestions(formattedAddress)
	if err != nil {
		log.Println(err)
	}
	// spew.Dump(nbnLocations)

	prettyPrint(nbnLocations)
}

func prettyPrint(nbnLocations []NbnLapi) {
	b, err := json.MarshalIndent(nbnLocations, "", "  ")
	if err != nil {
		log.Printf("Error: %s\n", err)
		return
	}
	fmt.Println(string(b))
}

type NbnPlaces struct {
	Timestamp   int64     `json:"timestamp"`
	Source      string    `json:"source"`
	Suggestions []NbnLapi `json:"suggestions"`
}

type NbnLapi struct {
	LocID            string  `json:"id"`
	FormattedAddress string  `json:"formattedAddress" `
	Latitude         float64 `json:"latitude" `
	Longitude        float64 `json:"longitude" `
}

func GetNBNSuggestions(address string) ([]NbnLapi, error) {
	encodedAddress := url.QueryEscape(address)
	NBNUrl := "https://places.nbnco.net.au/places/v1/autocomplete?query="
	thisUrl := fmt.Sprintf("%s%s", NBNUrl, encodedAddress)
	client := &http.Client{}
	log.Println(thisUrl)
	req, err := http.NewRequest(http.MethodGet, thisUrl, nil)
	if err != nil {
		fmt.Println(err)
		return []NbnLapi{}, err
	}
	req.Header.Add("Referer", "https://www.nbnco.com.au/when-do-i-get-it/rollout-map.html")
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)

	jsonData, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}

	var lapiResponse NbnPlaces
	err = json.Unmarshal(jsonData, &lapiResponse)
	if err != nil {
		log.Println(err)
		return []NbnLapi{}, err
	}
	// spew.Dump(lapiResponse)
	return lapiResponse.Suggestions, nil
}
