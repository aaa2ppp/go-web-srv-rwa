package api

import (
	"encoding/json"
	"net/http"

	"rwa/internal/dto"
)

type iUserRepo interface {
	Create(dto.CreateUser) (*dto.User, error)
	Get(token string) (*dto.User, error)
	Update(dto.UpdateUser) (*dto.User, error)
	Login(dto.LoginUser) (*dto.User, error)
	Logout(token string) error
}

type tUserResponse struct {
	User *dto.User `json:"user,omitempty"`
}

type User struct {
	repo iUserRepo
}

func NewUser(repo iUserRepo) *User {
	return &User{
		repo: repo,
	}
}

func (h *User) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		User dto.CreateUser `json:"user"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, httpError{http.StatusBadRequest, ""})
		return
	}

	user, err := h.repo.Create(req.User)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSONResponse(w, http.StatusCreated, tUserResponse{user})
}

func (h *User) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		User dto.LoginUser `json:"user"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, httpError{http.StatusBadRequest, ""})
		return
	}

	user, err := h.repo.Login(req.User)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSONResponse(w, http.StatusOK, tUserResponse{user})
}

func (h *User) Logout(w http.ResponseWriter, r *http.Request) {
	token, err := getAuthorizationToken(r)
	if err != nil {
		writeError(w, err)
		return
	}

	if err := h.repo.Logout(token); err != nil {
		writeError(w, err)
		return
	}
}

func (h *User) Get(w http.ResponseWriter, r *http.Request) {
	token, err := getAuthorizationToken(r)
	if err != nil {
		writeError(w, err)
		return
	}

	user, err := h.repo.Get(token)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSONResponse(w, http.StatusOK, tUserResponse{user})
}

func (h *User) Update(w http.ResponseWriter, r *http.Request) {
	token, err := getAuthorizationToken(r)
	if err != nil {
		writeError(w, err)
		return
	}

	var req struct {
		User dto.UpdateUser `json:"user"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, httpError{http.StatusBadRequest, ""})
		return
	}
	req.User.Token = &token

	user, err := h.repo.Update(req.User)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSONResponse(w, http.StatusOK, tUserResponse{user})
}
