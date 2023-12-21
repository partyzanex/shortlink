package link

import (
	"context"
	"database/sql"
	"net/url"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/partyzanex/shortlink/internal/boiler"
	"github.com/partyzanex/shortlink/pkg/enc"
	"github.com/partyzanex/shortlink/pkg/ptr"
)

var (
	ErrLinkNotFound = errors.New("link not found")
)

type Service struct {
	db *sql.DB

	readTimeout  time.Duration
	writeTimeout time.Duration
	hashLength   int
}

func NewService(db *sql.DB, readTimeout, writeTimeout time.Duration, hashLength int) *Service {
	return &Service{
		db:           db,
		readTimeout:  readTimeout,
		writeTimeout: writeTimeout,
		hashLength:   hashLength,
	}
}

func (s *Service) Create(ctx context.Context, uri *url.URL, expiredAt *time.Time) (*ID, error) {
	n := s.hashLength

CreateID:
	id, err := enc.Hash(n)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create hash")
	}

	model := &boiler.Link{
		ID:        id,
		Schema:    uri.Scheme,
		Domain:    uri.Host,
		URI:       uri.RequestURI(),
		CreatedAt: time.Time{},
		ExpiredAt: null.Time{},
	}

	if expiredAt != nil {
		model.ExpiredAt.SetValid(*expiredAt)
	}

	err = s.create(ctx, model)
	if err != nil {
		if e, ok := errors.Cause(err).(*pq.Error); ok && e.Code == "23505" {
			n++
			goto CreateID
		}

		return nil, err
	}

	return &id, err
}

func (s *Service) create(ctx context.Context, model *boiler.Link) error {
	ctx, cancel := context.WithTimeout(ctx, s.writeTimeout)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return errors.Wrap(err, "cannot begin transaction")
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else if txErr := tx.Commit(); txErr != nil {
			err = multierror.Append(err, txErr)
		}
	}()

	err = model.Insert(ctx, tx, boil.Infer())
	if err != nil {
		return errors.Wrap(err, "cannot insert link")
	}

	return err
}

func (s *Service) Get(ctx context.Context, id *ID) (*Link, error) {
	ctx, cancel := context.WithTimeout(ctx, s.readTimeout)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadUncommitted,
		ReadOnly:  true,
	})
	if err != nil {
		return nil, errors.Wrap(err, "cannot begin transaction")
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else if txErr := tx.Commit(); txErr != nil {
			err = multierror.Append(err, txErr)
		}
	}()

	model, err := boiler.FindLink(ctx, tx, *id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrLinkNotFound
		}

		return nil, errors.Wrap(err, "cannot find link")
	}

	var expiredAt *time.Time

	if model.ExpiredAt.Valid {
		expiredAt = ptr.Ptr(model.ExpiredAt.Time)
	}

	return &Link{
		ID:        model.ID,
		Schema:    Schema(model.Schema),
		Domain:    model.Domain,
		URI:       model.URI,
		CreatedAt: model.CreatedAt,
		ExpiredAt: expiredAt,
	}, err
}
