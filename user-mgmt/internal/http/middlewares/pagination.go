package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const (
	PageKey   string = "page"
	LimitKey  string = "limit"
	SortKey   string = "sort"
	FilterKey string = "filter"
	SearchKey string = "search"
)

const (
	DefaultLimit int    = 10
	DefaultPage  int    = 1
	DefaultSort  string = "created_at asc"
)

func Paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// defaults
		limit := DefaultLimit
		page := DefaultPage
		sort := DefaultSort
		filter := map[string]string{}
		search := map[string]string{}

		query := r.URL.Query()
		for key, value := range query {
			queryValue := value[len(value)-1]
			switch key {
			case "limit":
				limit, _ = strconv.Atoi(queryValue)
				break
			case "page":
				page, _ = strconv.Atoi(queryValue)
				break
			case "sort":
				formattedSort, err := validateSort(queryValue)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				sort = formattedSort
				break
			case "filter":
				formattedFilter, err := validateFilter(queryValue)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				filter = formattedFilter
				break
			case "search":
				formatedSearch, err := validateSearch(queryValue)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				search = formatedSearch
				break
			}
		}

		// set final pagination context values
		ctx := context.WithValue(r.Context(), LimitKey, limit)
		ctx = context.WithValue(ctx, PageKey, page)
		ctx = context.WithValue(ctx, SortKey, sort)
		ctx = context.WithValue(ctx, FilterKey, filter)
		ctx = context.WithValue(ctx, SearchKey, search)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func validateFilter(filter string) (map[string]string, error) {
	splits := strings.Split(filter, ".")
	if len(splits) != 2 {
		return nil, fmt.Errorf("malformed filter query parameter, should be field.filter")
	}

	field, value := splits[0], splits[1]
	return map[string]string{field: value}, nil
}

func validateSearch(filter string) (map[string]string, error) {
	splits := strings.Split(filter, ".")
	if len(splits) != 2 {
		return nil, fmt.Errorf("malformed search query parameter, should be field.value")
	}

	field, value := splits[0], splits[1]
	return map[string]string{field: value}, nil
}

func validateSort(sort string) (string, error) {
	splits := strings.Split(sort, ".")
	if len(splits) != 2 {
		return "", fmt.Errorf("malformed sort query, should be field.orderdirection")
	}

	field, order := splits[0], splits[1]
	if order != "desc" && order != "asc" {
		return "", fmt.Errorf("malformed order in sort query, should be asc or desc")
	}

	return fmt.Sprintf("%s %s", field, strings.ToUpper(order)), nil
}
