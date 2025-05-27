package getTorrents

import (
	"context"
	"encoding/json"
	"sync"
	"torrentServer/internal/services/jackett"
)

var (
	jackettInstance *jackett.Jackett
	once            sync.Once
)

func GetJackettInstance() *jackett.Jackett {
	once.Do(func() {
		jackettInstance = jackett.NewJackett(&jackett.Settings{
			ApiURL: "http://jackett:9117",              // http://localhost:9117
			ApiKey: "4k1y8yde4djn0kzzlcjlx47syqrlzzst", //v7mi9deijytd3qzz359phbs1i2s31xwo
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
