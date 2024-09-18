package main

import (
	"net/http"

	"rwa/internal/api"
	"rwa/internal/http_debug"
	"rwa/internal/repo"
)

// сюда писать код

type tMethods map[string]http.Handler

func (m tMethods) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for k, h := range m {
		if r.Method == k {
			h.ServeHTTP(w, r)
			return
		}
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func GetApp() http.Handler {
	mux := http.NewServeMux()
	repo := repo.New()

	{
		h := api.NewUser(repo.UserRepo())
		mux.Handle("/api/user", tMethods{
			"GET": http.HandlerFunc(h.Get),
			"PUT": http.HandlerFunc(h.Update),
		})
		mux.Handle("/api/users", tMethods{
			"POST": http.HandlerFunc(h.Register),
		})
		mux.Handle("/api/users/login", tMethods{
			"POST": http.HandlerFunc(h.Login),
		})
		mux.Handle("/api/user/logout", tMethods{
			"POST": http.HandlerFunc(h.Logout),
		})
	}

	{
		h := api.NewArticle(repo.ArticleRepo())
		mux.Handle("/api/articles", tMethods{
			"GET":  http.HandlerFunc(h.List),
			"POST": http.HandlerFunc(h.Create),
		})
		// mux.HandleFunc("/articles/feed", h.ListFollow)
		// mux.Handle("/api/articles/", http.StripPrefix("/api/articles", tMethodMux{
		// 	"GET":    http.HandlerFunc(h.Get),
		// 	"PUT":    http.HandlerFunc(h.Update),
		// 	"DELETE": http.HandlerFunc(h.Delete),
		// }))
	}

	return http_debug.Handle(mux)
}
