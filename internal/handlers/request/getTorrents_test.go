package getTorrents_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"torrentServer/internal/services/jackett"
)

// В пакете jackett (или в getTorrents, если не хотите менять jackett)
type JackettClient interface {
	Fetch(ctx context.Context, req *jackett.FetchRequest) (*jackett.FetchResponse, error)
	FilterResults(results []jackett.Result) ([]byte, error)
}

func GetJackettInstance() *jackett.Jackett {
	jackettInstance := jackett.NewJackett(&jackett.Settings{
		ApiURL: "",
		ApiKey: "",
	})
	return jackettInstance
}

// Мок JackettClient
type mockJackettClient struct {
	fetchFunc         func(ctx context.Context, req *jackett.FetchRequest) (*jackett.FetchResponse, error)
	filterResultsFunc func(results []jackett.Result) ([]byte, error)
}

func (m *mockJackettClient) Fetch(ctx context.Context, req *jackett.FetchRequest) (*jackett.FetchResponse, error) {
	return m.fetchFunc(ctx, req)
}

func (m *mockJackettClient) FilterResults(results []jackett.Result) ([]byte, error) {
	return m.filterResultsFunc(results)
}

// Переменная для подмены клиента в тестах
var testJackettClient JackettClient

// Обновим функции для использования testJackettClient, если он задан
func getClient() JackettClient {
	if testJackettClient != nil {
		return testJackettClient
	}
	return GetJackettInstance()
}

// Обновлённые функции для теста (пример для RequestSimple)
func RequestSimple(query string, categories []uint) (string, error) {
	ctx := context.Background()
	j := getClient()
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

func RequestAll(query string, categories []uint) (string, error) {
	ctx := context.Background()
	j := getClient()
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

// Тесты

func TestRequestSimple_Success(t *testing.T) {
	testJackettClient = &mockJackettClient{
		fetchFunc: func(ctx context.Context, req *jackett.FetchRequest) (*jackett.FetchResponse, error) {
			return &jackett.FetchResponse{
				Results: []jackett.Result{
					{Title: "Test1"},
					{Title: "Test2"},
				},
			}, nil
		},
		filterResultsFunc: func(results []jackett.Result) ([]byte, error) {
			return []byte(`[{"title":"Test1"},{"title":"Test2"}]`), nil
		},
	}
	defer func() { testJackettClient = nil }()

	res, err := RequestSimple("query", []uint{1, 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := `[{"title":"Test1"},{"title":"Test2"}]`
	if res != expected {
		t.Errorf("expected %s, got %s", expected, res)
	}
}

func TestRequestSimple_FetchError(t *testing.T) {
	testJackettClient = &mockJackettClient{
		fetchFunc: func(ctx context.Context, req *jackett.FetchRequest) (*jackett.FetchResponse, error) {
			return nil, errors.New("fetch error")
		},
		filterResultsFunc: nil,
	}
	defer func() { testJackettClient = nil }()

	_, err := RequestSimple("query", []uint{1})
	if err == nil || err.Error() != "fetch error" {
		t.Errorf("expected fetch error, got %v", err)
	}
}

func TestRequestSimple_FilterResultsError(t *testing.T) {
	testJackettClient = &mockJackettClient{
		fetchFunc: func(ctx context.Context, req *jackett.FetchRequest) (*jackett.FetchResponse, error) {
			return &jackett.FetchResponse{
				Results: []jackett.Result{{Title: "Test"}},
			}, nil
		},
		filterResultsFunc: func(results []jackett.Result) ([]byte, error) {
			return nil, errors.New("filter error")
		},
	}
	defer func() { testJackettClient = nil }()

	_, err := RequestSimple("query", []uint{1})
	if err == nil || err.Error() != "filter error" {
		t.Errorf("expected filter error, got %v", err)
	}
}

func TestRequestAll_Success(t *testing.T) {
	testJackettClient = &mockJackettClient{
		fetchFunc: func(ctx context.Context, req *jackett.FetchRequest) (*jackett.FetchResponse, error) {
			return &jackett.FetchResponse{
				Results: []jackett.Result{
					{Title: "Test1"},
					{Title: "Test2"},
				},
			}, nil
		},
		filterResultsFunc: nil,
	}
	defer func() { testJackettClient = nil }()

	res, err := RequestAll("query", []uint{1, 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Проверяем, что результат - валидный JSON с ожидаемыми данными
	var results []jackett.Result
	if err := json.Unmarshal([]byte(res), &results); err != nil {
		t.Fatalf("failed to unmarshal result: %v", err)
	}
	if len(results) != 2 || results[0].Title != "Test1" {
		t.Errorf("unexpected results: %+v", results)
	}
}

func TestRequestAll_FetchError(t *testing.T) {
	testJackettClient = &mockJackettClient{
		fetchFunc: func(ctx context.Context, req *jackett.FetchRequest) (*jackett.FetchResponse, error) {
			return nil, errors.New("fetch error")
		},
	}
	defer func() { testJackettClient = nil }()

	_, err := RequestAll("query", []uint{1})
	if err == nil || err.Error() != "fetch error" {
		t.Errorf("expected fetch error, got %v", err)
	}
}
