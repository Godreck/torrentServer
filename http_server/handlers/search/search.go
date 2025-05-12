package search

import (
	"net/http"
	"strconv"
	"strings"

	getTorrents "torrentServer/internal/handlers/request"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	categoriesStr := r.URL.Query().Get("categories") // Пример категории для фильмов

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

	jsonStr, err := getTorrents.RequestSimple(query, categories)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonStr))
}
