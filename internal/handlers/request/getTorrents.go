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
			ApiURL: "http://localhost:9117",
			ApiKey: "4k1y8yde4djn0kzzlcjlx47syqrlzzst",
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
