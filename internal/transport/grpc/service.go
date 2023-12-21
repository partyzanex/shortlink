package grpc

import (
	"context"
	"net/url"
	"time"

	"github.com/partyzanex/shortlink/internal/link"
	pb "github.com/partyzanex/shortlink/pkg/proto/go/shorts"
	"github.com/partyzanex/shortlink/pkg/ptr"
	"github.com/pkg/errors"
)

type LinkService interface {
	Create(ctx context.Context, uri *url.URL, expiredAt *time.Time) (*link.ID, error)
	Get(ctx context.Context, id *link.ID) (*link.Link, error)
}

type Service struct {
	pb.UnimplementedShortsServer

	linkService LinkService
}

func NewService(linkService LinkService) *Service {
	return &Service{
		UnimplementedShortsServer: pb.UnimplementedShortsServer{},
		linkService:               linkService,
	}
}

func (s *Service) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	uri, err := url.Parse(req.GetTargetUrl())
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse target URL")
	}

	if uri.Scheme == "" {
		uri.Scheme = "http"
	}

	var expiredAt *time.Time

	if req.ExpiredAt != nil {
		expiredAt = ptr.Ptr(time.Unix(req.GetExpiredAt(), 0))
	}

	shortID, err := s.linkService.Create(ctx, uri, expiredAt)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create short link")
	}

	return &pb.CreateResponse{
		ShortLink: *shortID,
	}, nil
}

func (s *Service) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	result, err := s.linkService.Get(ctx, ptr.Ptr(req.ShortLink))
	if err != nil {
		return nil, errors.Wrap(err, "cannot get link")
	}

	uri, err := url.Parse(result.RawURL())
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse link result")
	}

	return &pb.GetResponse{
		TargetUrl: uri.String(),
	}, nil
}
