package grpc

import (
	"context"
	"log"
	"net"

	"github.com/Renal37/musthave_shortener_tpl.git/internal/services"
	pb "github.com/Renal37/musthave_shortener_tpl.git/proto/shortener"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	pb.UnimplementedURLShortenerServer
	service *services.ShortenerService
}

func NewGRPCServer(service *services.ShortenerService) *GRPCServer {
	return &GRPCServer{service: service}
}

func (s *GRPCServer) ShortenURL(ctx context.Context, req *pb.ShortenURLRequest) (*pb.ShortenURLResponse, error) {
	shortURL, err := s.service.GetExistURL(req.Url, nil)
	if err != nil {
		log.Printf("Error shortening URL: %v", err)
		return nil, err
	}
	return &pb.ShortenURLResponse{ShortUrl: shortURL}, nil
}

func (s *GRPCServer) GetOriginalURL(ctx context.Context, req *pb.GetOriginalURLRequest) (*pb.GetOriginalURLResponse, error) {
	originalURL, err := s.service.GetRep(req.ShortId, "")
	if err != nil {
		log.Printf("Error getting original URL: %v", err)
		return nil, err
	}
	return &pb.GetOriginalURLResponse{OriginalUrl: originalURL}, nil
}

func (s *GRPCServer) DeleteUserURLs(ctx context.Context, req *pb.DeleteUserURLsRequest) (*pb.DeleteUserURLsResponse, error) {
	err := s.service.DeleteURLsRep(req.UserId, req.ShortUrls)
	if err != nil {
		log.Printf("Error deleting user URLs: %v", err)
		return nil, err
	}
	return &pb.DeleteUserURLsResponse{}, nil
}

func (s *GRPCServer) GetUserURLs(ctx context.Context, req *pb.GetUserURLsRequest) (*pb.GetUserURLsResponse, error) {
	urls, err := s.service.GetFullRep(req.UserId)
	if err != nil {
		log.Printf("Error getting user URLs: %v", err)
		return nil, err
	}
	var urlResponses []*pb.URLResponse
	for _, url := range urls {
		urlResponses = append(urlResponses, &pb.URLResponse{
			CorrelationId: url["uuid"],
			ShortUrl:      url["short_url"],
		})
	}
	return &pb.GetUserURLsResponse{Urls: urlResponses}, nil
}

func (s *GRPCServer) GetStats(ctx context.Context, req *pb.GetStatsRequest) (*pb.GetStatsResponse, error) {
	urlCount, err := s.service.GetURLCount()
	if err != nil {
		log.Printf("Error getting URL count: %v", err)
		return nil, err
	}
	userCount, err := s.service.GetUserCount()
	if err != nil {
		log.Printf("Error getting user count: %v", err)
		return nil, err
	}
	return &pb.GetStatsResponse{
		UrlCount:  int32(urlCount),
		UserCount: int32(userCount),
	}, nil
}

func StartGRPCServer(ctx context.Context, service *services.ShortenerService, address string) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
		return err
	}

	grpcServer := grpc.NewServer()
	pb.RegisterURLShortenerServer(grpcServer, NewGRPCServer(service))

	log.Printf("gRPC сервер слушает на %v", lis.Addr())
	go grpcServer.Serve(lis)

	<-ctx.Done()
	// Остановка gRPC сервера
	grpcServer.Stop()
	return nil
}
