package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"rwa/internal/dto"
)

type tArticleResponse struct {
	Article *dto.Article `json:"article,omitempty"`
}

type tArticleListResponse struct {
	Articles      []dto.Article `json:"articles,omitempty"`
	ArticlesCount int           `json:"articlesCount,omitempty"`
}

type iArticleRepo interface {
	Create(dto.CreateArticle) (*dto.Article, error)
	List(dto.ListArticle) ([]dto.Article, error)
	// ListFollow(tListFollowRequest) ([]model.Article, error)
	// Get(tGetRequest) (model.Article, error)
	// Update(tUpdateRequest) (bool, error)
	// Delete(tDeleteRequest) (bool, error)
}

type Article struct {
	repo iArticleRepo
}

func NewArticle(repo iArticleRepo) *Article {
	return &Article{repo: repo}
}

func (h *Article) Create(w http.ResponseWriter, r *http.Request) {
	token, err := getAuthorizationToken(r)
	if err != nil {
		writeError(w, err)
		return
	}

	var req struct {
		Article dto.CreateArticle `json:"article,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, httpError{http.StatusBadRequest, ""})
		return
	}
	req.Article.Token = token

	article, err := h.repo.Create(req.Article)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSONResponse(w, http.StatusCreated, tArticleResponse{
		Article: article,
	})
}

func (h *Article) List(w http.ResponseWriter, r *http.Request) {
	req := dto.ListArticle{
		Tag:       r.FormValue("tag"),
		Author:    r.FormValue("author"),
		Favorited: r.FormValue("favorited"),
		Limit:     defaultLimit,
	}

	if v := r.FormValue("limit"); v != "" {
		limit, err := strconv.Atoi(v)
		if err != nil || limit <= 0 {
			writeError(w, httpError{http.StatusBadRequest, "limit must be int > 0"})
			return
		}
		req.Limit = uint64(limit)
	}

	if v := r.FormValue("offset"); v != "" {
		offset, err := strconv.Atoi(v)
		if err != nil || offset < 0 {
			writeError(w, httpError{http.StatusBadRequest, "offset must be int >= 0"})
			return
		}
		req.Offset = uint64(offset)
	}

	articles, err := h.repo.List(req)
	if err != nil {
		writeError(w, httpError{http.StatusInternalServerError, "internal error"})
		return
	}

	writeJSONResponse(w, http.StatusOK, tArticleListResponse{
		Articles:      articles,
		ArticlesCount: len(articles),
	})
}

// func (h *Articles) ListFollow(w http.ResponseWriter, r *http.Request) {
// 	req := tListFollowRequest{
// 		Limit: defaultLimit,
// 	}

// 	if v := r.FormValue("limit"); v != "" {
// 		limit, err := strconv.Atoi(v)
// 		if err != nil || limit <= 0 {
// 			writeError(w, httpError{http.StatusBadRequest, "limit must be int > 0"})
// 			return
// 		}
// 		req.Limit = uint(limit)
// 	}

// 	if v := r.FormValue("offset"); v != "" {
// 		offset, err := strconv.Atoi(v)
// 		if err != nil || offset < 0 {
// 			writeError(w, httpError{http.StatusBadRequest, "offset must be int >= 0"})
// 			return
// 		}
// 		req.Offset = uint(offset)
// 	}

// 	articles, err := h.repo.ListFollow(req)
// 	if err != nil {
// 		writeError(w, httpError{http.StatusInternalServerError, "internal error"})
// 		return
// 	}

// 	writeResponse(w, tListResponse{
// 		Articles:      articles,
// 		ArticlesCount: len(articles),
// 	})
// }

// func (h *Articles) getSlug(r *http.Request) (string, error) {
// 	const op = "Articles.getSlug"
// 	var slug string
// 	if v := r.URL.Path; v == "" {
// 		log.Printf("%s: unexpected request with empty path", op)
// 		return "", httpError{http.StatusInternalServerError, "internal error"}
// 	} else {
// 		slug = v[1:]
// 	}
// 	if slug == "" {
// 		return "", httpError{http.StatusBadRequest, "slug requred"}
// 	}
// 	return slug, nil
// }

// func (h *Articles) Get(w http.ResponseWriter, r *http.Request) {
// 	slug, err := h.getSlug(r)
// 	if err != nil {
// 		writeError(w, err)
// 		return
// 	}

// 	article, err := h.repo.Get(tGetRequest{
// 		Slug: slug,
// 	})
// 	if err != nil {
// 		writeError(w, err)
// 	}
// 	writeResponse(w, tGetResponse{
// 		Article: article,
// 	})
// }

// func (h *Articles) Update(w http.ResponseWriter, r *http.Request) {
// 	slug, err := h.getSlug(r)
// 	if err != nil {
// 		writeError(w, err)
// 		return
// 	}
// 	var req tUpdateRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		writeError(w, httpError{http.StatusBadRequest, "bad request"})
// 		return
// 	}
// 	req.Slug = slug

// 	ok, err := h.repo.Update(req)
// 	if err != nil {
// 		writeError(w, err)
// 		return
// 	}
// 	if !ok {
// 		writeError(w, httpError{http.StatusNotFound, "not found"})
// 		return
// 	}
// }

// func (h *Articles) Delete(w http.ResponseWriter, r *http.Request) {
// 	slug, err := h.getSlug(r)
// 	if err != nil {
// 		writeError(w, err)
// 		return
// 	}

// 	ok, err := h.repo.Delete(tDeleteRequest{
// 		Slug: slug,
// 	})
// 	if err != nil {
// 		writeError(w, err)
// 		return
// 	}
// 	if !ok {
// 		writeError(w, httpError{http.StatusNotFound, "not found"})
// 		return
// 	}
// }
