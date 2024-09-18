package repo

import (
	"sync"
	"time"
	"unsafe"
)

type tSession struct {
	id   uint64
	user *tUserRec
}

type tUserRec struct {
	id        uint64
	bio       string
	email     string
	image     string
	username  string
	password  string
	createdAt time.Time
	updatedAt time.Time
	following bool
}

type tTagList []string

func (l tTagList) includes(tag string) bool {
	for _, v := range l {
		if v == tag {
			return true
		}
	}
	return false
}

type tArticleRec struct {
	id          uint64
	createdAt   time.Time
	updatedAt   time.Time
	author      *tUserRec
	slug        string
	title       string
	tagList     tTagList
	description string
	body        string
	// favorited      bool
	// favoritesCount uint
}

type Repo struct {
	mu        sync.RWMutex
	users     map[uint64]*tUserRec
	userID    uint64
	sessions  map[uint64]*tSession
	sessionID uint64
	articles  map[uint64]*tArticleRec
	articleID uint64
}

func New() *Repo {
	return &Repo{
		users:    map[uint64]*tUserRec{},
		articles: map[uint64]*tArticleRec{},
		sessions: map[uint64]*tSession{},
	}
}

func (r *Repo) UserRepo() *UserRepo {
	return (*UserRepo)(unsafe.Pointer(r))
}

func (r *Repo) ArticleRepo() *ArticleRepo {
	return (*ArticleRepo)(unsafe.Pointer(r))
}

func (r *Repo) session() *Session {
	return (*Session)(unsafe.Pointer(r))
}
