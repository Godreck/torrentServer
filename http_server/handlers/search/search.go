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

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// парсим параметры
	query := r.URL.Query().Get("query")
	categoriesStr := r.URL.Query().Get("categories")

	if query == "" || categoriesStr == "" {
		http.Error(w, "Query and categories parameters are required", http.StatusBadRequest)
		return
	}

	categories := make([]uint, 0)
	for _, c := range strings.Split(categoriesStr, ",") {
		category, err := strconv.ParseUint(c, 10, 32)
		if err != nil {
			http.Error(w, "Invalid category format", http.StatusBadRequest)
			return
		}
		categories = append(categories, uint(category))
	}

	// Генерируем ключ кэша
	cacheKey := generateCacheKey(query, categories)

	// Проверяем кэш
	var jsonStr string
	if ok := redisCache.Get(ctx, cacheKey, &jsonStr); ok {
		log.Printf("Cache hit for key: %s", cacheKey)
		json.NewEncoder(w).Encode(jsonStr)
		return
	}

	log.Printf("Cache miss for key: %s", cacheKey)
	jsonStr, err := getTorrents.RequestSimple(query, categories)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Сохраняем в кэш
	if err := redisCache.Set(ctx, cacheKey, jsonStr); err != nil {
		// Логируем ошибку, но не прерываем выполнение
		log.Printf("Cache set error: %v", err)
	}

	// w.Header().Set("Content-Type", "application/json")
	// w.Write([]byte(jsonStr))
	json.NewEncoder(w).Encode(jsonStr) //попробовать этот вариант
}

func generateCacheKey(query string, categories []uint) string {
	sort.Slice(categories, func(i, j int) bool { return categories[i] < categories[j] })
	return fmt.Sprintf("query=%s&categories=%v", query, categories)
}
