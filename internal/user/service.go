package user

import (
	"context"
	"errors"
	"esi/internal/pkg/conf"
	"esi/internal/pkg/token"
	"esi/internal/pkg/util"
	"time"

	"github.com/google/uuid"
)

type Repoer interface {
	create(ctx context.Context, u *User) error
	fetch(ctx context.Context) ([]*User, error)
	getByID(ctx context.Context, id UUID) (*User, error)
	getByEmail(ctx context.Context, email string) (*User, error)
	update(ctx context.Context, u *User) error
}

// Interface check
var _ Repoer = (*Repo)(nil)

/* ----- */ /* ----- */ /* ----- */

type TokenUtil interface {
	GenTokenPair(u *User) (string, string, error)
	RefreshToken(ctx context.Context, token_str string) (string, error)
	BlockToken(ctx context.Context, token_str string) error
}

// Interface check
var _ TokenUtil = (*token.Token)(nil)

/* ===== */ /* ===== */ /* ===== */

type Service struct {
	repo    Repoer
	token   TokenUtil
	timeout int
}

func NewService(cfg *conf.Config, r Repoer, t TokenUtil) *Service {
	return &Service{
		repo:    r,
		token:   t,
		timeout: cfg.Server.CtxTimeout,
	}
}

func (s *Service) login(ctx context.Context, req *LoginReq) (*LoginResp, error) {
	user, err := s.getUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if !util.ValidatePwd(req.Pwd, user.Pwd) {
		return nil, errors.New("invalid password")
	}
	access_token, refresh_token, err := s.token.GenTokenPair(user)
	if err != nil {
		return nil, err
	}
	var resp = &LoginResp{
		AccessToken:  access_token,
		RefreshToken: refresh_token,
	}
	return resp, nil
}

func (s *Service) logout(ctx context.Context, access_id string) error {
	if err := s.token.BlockToken(ctx, access_id); err != nil {
		return err
	}
	return nil
}

func (s *Service) register(ctx context.Context, req *RegReq) (*RegResp, error) {
	if _, err := s.getUserByEmail(ctx, req.Email); err == nil {
		return nil, errors.New("user exists")
	}
	if err := util.EncryptPwd(&req.Pwd); err != nil {
		return nil, err
	}
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	var user = &User{
		ID:       id,
		Username: req.Username,
		Email:    req.Email,
		Pwd:      req.Pwd,
	}
	if err := s.repo.create(ctx, user); err != nil {
		return nil, err
	}
	access_token, refresh_token, err := s.token.GenTokenPair(user)
	if err != nil {
		return nil, err
	}
	var resp = &RegResp{
		AccessToken:  access_token,
		RefreshToken: refresh_token,
	}
	return resp, nil
}

func (s *Service) refresh(ctx context.Context, refresh_token string) (*RefreshResp, error) {
	access_token, err := s.token.RefreshToken(ctx, refresh_token)
	if err != nil {
		return nil, err
	}
	var resp = &RefreshResp{
		AccessToken: access_token,
	}
	return resp, nil
}

func (s *Service) profile(ctx context.Context) ([]*User, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(s.timeout)*time.Second)
	defer cancel()
	return s.repo.fetch(ctx)
}

func (s *Service) getUserByID(ctx context.Context, id UUID) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(s.timeout)*time.Second)
	defer cancel()
	return s.repo.getByID(ctx, id)
}

func (s *Service) getUserByEmail(ctx context.Context, email string) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(s.timeout)*time.Second)
	defer cancel()
	return s.repo.getByEmail(ctx, email)
}
