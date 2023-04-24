package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type TagWithCount struct {
	Tag   string `json:"tag"`
	Count int    `json:"count"`
}

type InputParams struct {
	Lat    float64 `json:"lat"`
	Lng    float64 `json:"lng"`
	Radius int     `json:"radius"`
}

type NearbyCountResult struct {
	InputParams InputParams    `json:"input_params"`
	Results     []TagWithCount `json:"results"`
}

func validate_coordinates(lat, lng float64) bool {
	return lat >= -90 && lat <= 90 && lng >= -180 && lng <= 180
}

func read_tags_from_file(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var types []string

	for scanner.Scan() {
		types = append(types, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return types
}

func get_nearby_places(lat, lng float64, radius int, tag, key string, next_page_token *string) (map[string]interface{}, error) {
	url := "https://maps.googleapis.com/maps/api/place/nearbysearch/json"

	var params map[string]string
	if next_page_token == nil {
		params = map[string]string{
			"location": fmt.Sprintf("%f,%f", lat, lng),
			"radius":   fmt.Sprintf("%d", radius),
			"type":     tag,
			"key":      key,
		}
	} else {
		params = map[string]string{
			"key":       key,
			"pagetoken": *next_page_token,
		}
	}

	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}

	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	err = json.Unmarshal(bodyBytes, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func get_nearyby_tag_count(tag string, lat, lng float64, radius int, key string) (int, error) {
	var count int = 0

	results, err := get_nearby_places(lat, lng, radius, tag, key, nil)
	if err != nil {
		return -1, err
	}

	count += len(results["results"].([]interface{}))

	for {
		if next_page_token, ok := results["next_page_token"].(string); ok {
			// API is rate limited. Sleep for 2 seconds.
			time.Sleep(2000 * time.Millisecond)
			results, err = get_nearby_places(lat, lng, radius, tag, key, &next_page_token)
			if err != nil {
				return -1, err
			}
			count += len(results["results"].([]interface{}))
		} else {
			break
		}
	}

	return count, nil
}

func get_nearby_tags_count(tags []string, lat, lng float64, radius int, key string) (NearbyCountResult, error) {
	var wg sync.WaitGroup
	results := []TagWithCount{}

	for _, tag := range tags {
		wg.Add(1)

		go func(tag string) {
			defer wg.Done()
			data, _ := get_nearyby_tag_count(tag, lat, lng, int(radius), key)
			results = append(results, TagWithCount{tag, data})
		}(tag)
	}

	wg.Wait()

	response := NearbyCountResult{
		InputParams: InputParams{
			Lat:    lat,
			Lng:    lng,
			Radius: int(radius),
		},
		Results: results,
	}

	return response, nil
}
