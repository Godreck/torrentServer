package getTorrents

import (
	"context"
	"encoding/json"
	"os"
	"sync"
	jackett "torrentServer/internal/services/jackett"
)

var (
	jackettInstance *jackett.Jackett
	once            sync.Once
)

var (
	apiURL = os.Getenv("JACKETT_API_URL")
	apiKey = os.Getenv("JACKETT_API_KEY")
)

// func init() {
// 	if v, ok := os.LookupEnv("JACKETT_API_URL"); ok {
// 		apiURL = v
// 		log.Print("jackett api_url found status:", ok)
// 	} else {
// 		log.Print("jackett api_url NOT FOUND")
// 	}
// 	if v, ok := os.LookupEnv("JACKETT_API_KEY"); ok {
// 		apiKey = v
// 		log.Print("jackett api_key found status:", ok)
// 	} else {
// 		log.Print("jackett api_key NOT FOUND")

// 	}
// }

func GetJackettInstance() *jackett.Jackett {
	once.Do(func() {
		jackettInstance = jackett.NewJackett(&jackett.Settings{
			ApiURL: apiURL,
			ApiKey: apiKey,
		})
	})
	return jackettInstance
}

func RequestAll(query string, categories []uint) (string, error) {
	ctx := context.Background()
	j := GetJackettInstance()
	resp, err := j.Fetch(ctx, &jackett.FetchRequest{
		Categories: categories,
		Query:      query,
	})
	if err != nil {
		return "", err
	}

	jsonBytes, err := json.MarshalIndent(resp.Results, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

func RequestSimple(query string, categories []uint) (string, error) {
	ctx := context.Background()
	j := GetJackettInstance()
	resp, err := j.Fetch(ctx, &jackett.FetchRequest{
		Categories: categories,
		Query:      query,
	})
	if err != nil {
		return "", err
	}

	simpleRes, err := j.FilterResults(resp.Results)
	if err != nil {
		return "", err
	}

	return string(simpleRes), nil
}
