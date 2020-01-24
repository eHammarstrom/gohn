package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type jHNItem struct {
	Author  string `json:"by"`
	ID      uint   `json:"id"`
	Score   uint   `json:"score"`
	JobText string `json:"text"`
	Time    uint   `json:"time"`
	Title   string `json:"title"`
	TypeStr string `json:"type"`
	URL     string `json:"url"`
}

type jHNItemIDs []uint

const baseURL string = "https://hacker-news.firebaseio.com/v0"
const topStories string = baseURL + "/topstories.json"
const numStories uint = 10

func eprint(fmtStr string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, fmtStr+"\n", a...)
}

func itemIDToURL(id uint) string {
	return baseURL + "/item/" + strconv.FormatUint(uint64(id), 10) + ".json"
}

func getItemByItemID(id uint) (*jHNItem, error) {
	resp, err := http.Get(itemIDToURL(id))
	if err != nil {
		return nil, err
	}

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var item *jHNItem = new(jHNItem)
	if err := json.Unmarshal(bs, item); err != nil {
		return nil, err
	}

	return item, nil
}

func main() {
	// get all top items (stories/jobs)
	resp, err := http.Get(topStories)
	if err != nil || resp.StatusCode != 200 {
		eprint("Failed to retrieve top stories from %s", topStories)
		return
	}

	// read http.Body into []byte
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		eprint("Failed to parse HTTP body from %s", topStories)
		return
	}

	// parse the items ([]byte) into json (jHNItems)
	var itemIDs jHNItemIDs
	if err = json.Unmarshal(bs, &itemIDs); err != nil {
		eprint("Got invalid json from %s", topStories)
		eprint("Data: %s", bs)
		return
	}

	// retrieve the story/job
	var items = make([]string, 0, numStories)
	for idx, itemID := range itemIDs[:numStories] {
		itemStr := strconv.Itoa(idx+1) + ": "

		item, err := getItemByItemID(itemID)
		if err != nil {
			itemStr += "Error retrieving item."
			itemStr += "\n\tReason: " + err.Error()
		} else {
			itemStr += item.Title + "\n\t (" + item.URL + ")"
		}

		if idx < int(numStories)-1 {
			itemStr += "\n"
		}

		items = append(items, itemStr)
	}

	for _, item := range items {
		fmt.Println(item)
	}
}
