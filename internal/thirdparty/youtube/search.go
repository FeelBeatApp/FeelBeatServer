package youtube

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/buger/jsonparser"
)

type SearchResult struct {
	VideoId  string
	Title    string
	Duration time.Duration
}

const ytBaseUrl = "https://www.youtube.com/results?search_query="
const loopSearchTreshold = 5

func fetchRawPage(query string) ([]byte, error) {
	res, err := http.DefaultClient.Get(ytBaseUrl + url.QueryEscape(query))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	rawHtml, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return rawHtml, nil
}

func extractJsonContent(rawHtml []byte) string {
	firstSplits := strings.Split(string(rawHtml), `var ytInitialData = `)
	return strings.Split(firstSplits[1], `;</script>`)[0]

}

func searchVideosInJson(jsonString string) ([]byte, error) {
	for index := 0; index < loopSearchTreshold; index++ {
		vidoes, _, _, err := jsonparser.Get([]byte(jsonString), fmt.Sprintf("[%d]", index), "itemSectionRenderer", "contents")
		if err == nil {
			return vidoes, nil
		}

		index++
	}

	return nil, errors.New("Couldn't find results in youtube response")
}

func parseJson(jsonString string) ([]SearchResult, error) {
	videosSection, _, _, err := jsonparser.Get([]byte(jsonString), "contents", "twoColumnSearchResultsRenderer", "primaryContents", "sectionListRenderer", "contents")
	if err != nil {
		return nil, err
	}
	videosListJson, err := searchVideosInJson(string(videosSection))
	if err != nil {
		return nil, err
	}

	result := make([]SearchResult, 0)

	_, err = jsonparser.ArrayEach(videosListJson, func(value []byte, _ jsonparser.ValueType, _ int, _ error) {
		video, _, _, err := jsonparser.Get(value, "videoRenderer")
		if err != nil {
			return
		}

		id, err := jsonparser.GetString(video, "videoId")
		if err != nil {
			return
		}
		title, err := jsonparser.GetString(video, "title", "runs", "[0]", "text")
		if err != nil {
			return
		}
		durationString, err := jsonparser.GetString(video, "lengthText", "simpleText")
		if err != nil {
			return
		}
		durations := strings.Split(durationString, ":")
		if len(durations) > 2 {
			return
		}

		minutes, err := strconv.Atoi(durations[0])
		if err != nil {
			return
		}

		seconds, err := strconv.Atoi(durations[1])
		if err != nil {
			return
		}

		result = append(result, SearchResult{
			VideoId:  id,
			Title:    title,
			Duration: time.Duration(seconds)*time.Second + time.Duration(minutes)*time.Minute,
		})
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func Search(query string) ([]SearchResult, error) {
	rawHtml, err := fetchRawPage(query)
	if err != nil {
		return nil, fmt.Errorf("Youtube search failed: %w", err)
	}

	jsonString := extractJsonContent(rawHtml)

	results, err := parseJson(jsonString)
	if err != nil {
		return nil, fmt.Errorf("Youtbe search parsing failed: %w", err)
	}

	return results, nil
}
