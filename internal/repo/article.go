package repo

import (
	"sort"
	"strconv"
	"time"

	"rwa/internal/dto"
)

func newProfileResponse(rec *tUserRec) *dto.Profile {
	if rec == nil {
		return nil
	}
	return &dto.Profile{
		BIO:       rec.bio,
		Following: rec.following,
		Image:     rec.image,
		Username:  rec.username,
	}
}

func newArticleResponse(rec *tArticleRec) *dto.Article {
	if rec == nil {
		return nil
	}
	return &dto.Article{
		CreatedAt:   rec.createdAt,
		UpdatedAt:   rec.updatedAt,
		Slug:        rec.slug,
		Title:       rec.title,
		TagList:     rec.tagList,
		Description: rec.description,
		Body:        rec.body,
		Author:      newProfileResponse(rec.author),
		// Favorited:      rec.favorited,
		// FavoritesCount: rec.favoritesCount,
	}
}

func newArticleListResponse(list []*tArticleRec) []dto.Article {
	res := make([]dto.Article, len(list))
	for i, rec := range list {
		res[i] = *newArticleResponse(rec)
	}
	return res
}

type ArticleRepo struct {
	Repo
}

func (r *ArticleRepo) genSlug(rec *tArticleRec) string {
	// TODO:
	slug := strconv.FormatUint(rec.id, 10)
	return slug
}

func (r *ArticleRepo) Create(req dto.CreateArticle) (*dto.Article, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	session, err := r.session().getByToken(req.Token)
	if err != nil {
		return nil, dto.ErrUnauthorized
	}

	r.articleID++
	now := time.Now()

	rec := &tArticleRec{
		id:          r.articleID,
		title:       req.Title,
		body:        req.Body,
		description: req.Description,
		tagList:     req.TagList,
		author:      session.user,
		createdAt:   now,
		updatedAt:   now,
	}

	rec.slug = r.genSlug(rec)
	r.articles[rec.id] = rec

	return newArticleResponse(rec), nil
}

func (r *ArticleRepo) List(req dto.ListArticle) ([]dto.Article, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var list []*tArticleRec
	for _, rec := range r.articles {
		if req.Author != "" && rec.author.username != req.Author {
			continue
		}
		if req.Tag != "" && !rec.tagList.includes(req.Tag) {
			continue
		}
		// TODO: check Favorited?
		list = append(list, rec)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].updatedAt.After(list[j].updatedAt)
	})

	min := func(a, b uint64) uint64 {
		if a > b {
			return b
		}
		return a
	}

	offset := min(req.Offset, uint64(len(list)))
	list = list[offset:]

	limit := min(req.Limit, uint64(len(list)))
	list = list[:limit]

	return newArticleListResponse(list), nil
}
