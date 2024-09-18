package repo

import (
	"errors"
	"time"

	"rwa/internal/dto"
)

func newUserResponse(rec *tUserRec, token string) *dto.User {
	if rec == nil {
		return nil
	}
	return &dto.User{
		BIO:       rec.bio,
		Email:     rec.email,
		Image:     rec.image,
		Token:     token,
		Username:  rec.username,
		CreatedAt: rec.createdAt,
		UpdatedAt: rec.updatedAt,
	}
}

type UserRepo struct {
	Repo
}

func (r *UserRepo) Create(req dto.CreateUser) (*dto.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.getByEmail(req.Email) != nil {
		return nil, errors.New("user already exists")
	}

	r.userID++
	now := time.Now()
	rec := &tUserRec{
		id:        r.userID,
		email:     req.Email,
		password:  req.Password,
		username:  req.Username,
		createdAt: now,
		updatedAt: now,
	}

	token, err := r.session().create(rec)
	if err != nil {
		return nil, err
	}

	r.users[rec.id] = rec
	return newUserResponse(rec, token), nil
}

func (r *UserRepo) Get(token string) (*dto.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	session, err := r.session().getByToken(token)
	if err != nil {
		return nil, dto.ErrUnauthorized
	}

	return newUserResponse(session.user, token), nil
}

func (r *UserRepo) Update(req dto.UpdateUser) (*dto.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if req.Token == nil {
		return nil, errors.New("token requered")
	}

	session, err := r.session().getByToken(*req.Token)
	if err != nil {
		return nil, dto.ErrUnauthorized
	}

	rec := session.user

	if req.BIO != nil {
		rec.bio = *req.BIO
	}
	if req.Email != nil {
		rec.email = *req.Email
	}
	if req.Image != nil {
		rec.image = *req.Image
	}
	if req.Username != nil {
		rec.username = *req.Username
	}
	rec.updatedAt = time.Now()

	return newUserResponse(rec, *req.Token), nil
}

func (r *UserRepo) getByEmail(email string) *tUserRec {
	for _, rec := range r.users {
		if rec.email == email {
			return rec
		}
	}
	return nil
}

func (r *UserRepo) Login(req dto.LoginUser) (*dto.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	rec := r.getByEmail(req.Email)
	if rec == nil || rec.password != req.Password {
		return nil, errors.New("bad login or password")
	}

	token, err := r.session().create(rec)
	if err != nil {
		return nil, err
	}

	return newUserResponse(rec, token), nil
}

func (r *UserRepo) Logout(token string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.session().delete(token)
}
