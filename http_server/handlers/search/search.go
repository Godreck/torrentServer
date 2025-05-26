package search

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	cache "torrentServer/cache"
	getTorrents "torrentServer/internal/handlers/request"
)

var redisCache = cache.NewRedisCache("localhost:6379", 24*time.Hour)

type PaginatedResponse struct {
	Data       []map[string]interface{} `json:"data"` // Изменено на массив объектов
	Page       int                      `json:"page"`
	PerPage    int                      `json:"per_page"`
	TotalItems int                      `json:"total_items"`
	TotalPages int                      `json:"total_pages"`
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Парсинг параметров
	query, categories, err := parseRequestParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Получаем данные (из кэша или Jackett)
	results, err := getOrFetchResults(ctx, query, categories)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Применяем пагинацию
	page, perPage := parsePaginationParams(r)
	paginatedData, totalPages := applyPagination(results, page, perPage)

	// Формируем ответ
	response := PaginatedResponse{
		Data:       paginatedData,
		Page:       page,
		PerPage:    perPage,
		TotalItems: len(results),
		TotalPages: totalPages,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Вспомогательные функции

func parseRequestParams(r *http.Request) (string, []uint, error) {
	query := r.URL.Query().Get("query")
	categoriesStr := r.URL.Query().Get("categories")

	if query == "" || categoriesStr == "" {
		return "", nil, fmt.Errorf("query and categories parameters are required")
	}

	categories := make([]uint, 0)
	for _, c := range strings.Split(categoriesStr, ",") {
		category, err := strconv.ParseUint(c, 10, 32)
		if err != nil {
			return "", nil, fmt.Errorf("invalid category format")
		}
		categories = append(categories, uint(category))
	}

	return query, categories, nil
}

func getOrFetchResults(ctx context.Context, query string, categories []uint) ([]map[string]interface{}, error) {
	cacheKey := generateCacheKey(query, categories)

	// Пытаемся получить из кэша
	var results []map[string]interface{}
	if redisCache.Get(ctx, cacheKey, &results) {
		return results, nil
	}

	// Запрос к Jackett
	jsonStr, err := getTorrents.RequestSimple(query, categories)
	if err != nil {
		return nil, err
	}

	// Парсим JSON
	if err := json.Unmarshal([]byte(jsonStr), &results); err != nil {
		return nil, fmt.Errorf("failed to parse Jackett response")
	}

	// Сохраняем в кэш
	if err := redisCache.Set(ctx, cacheKey, results); err != nil {
		log.Printf("Cache set error: %v", err)
	}

	return results, nil
}

func applyPagination(data []map[string]interface{}, page, perPage int) ([]map[string]interface{}, int) {
	totalItems := len(data)
	if totalItems == 0 {
		return []map[string]interface{}{}, 0
	}

	// Рассчитываем страницы
	totalPages := (totalItems + perPage - 1) / perPage
	if page > totalPages {
		page = totalPages
	}

	start := (page - 1) * perPage
	if start >= totalItems {
		start = totalItems - perPage
		if start < 0 {
			start = 0
		}
	}

	end := start + perPage
	if end > totalItems {
		end = totalItems
	}

	return data[start:end], totalPages
}

// Остальные функции остаются без изменений
func generateCacheKey(query string, categories []uint) string {
	sort.Slice(categories, func(i, j int) bool { return categories[i] < categories[j] })
	return fmt.Sprintf("query=%s&categories=%v", query, categories)
}

func parsePaginationParams(r *http.Request) (int, int) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	return page, perPage
}
