package grpc

// import (
//     "context"
//     "testing"
//     "github.com/Renal37/musthave_shortener_tpl.git/internal/services/mocks" 
//     pb "github.com/Renal37/musthave_shortener_tpl.git/proto"
//     "github.com/stretchr/testify/assert"
//     "github.com/stretchr/testify/mock"
//     "google.golang.org/grpc"
//     "net"
// )

// func TestGRPCServer_ShortenURL(t *testing.T) {
//     // Создаем мок для ShortenerService
//     mockService := new(mocks.NewShortenerService)

//     // Подготовим тестовые данные и ожидаемый результат
//     req := &pb.ShortenURLRequest{Url: "http://example.com"}
//     expectedShortURL := "short.ly/abc123"

//     // Настроим мок для метода GetExistURL
//     mockService.On("GetExistURL", req.Url, mock.Anything).Return(expectedShortURL, nil)

//     // Создаем gRPC сервер с мок-сервисом
//     server := NewGRPCServer(mockService)

//     // Запускаем сервер в горутине
//     go func() {
//         grpcServer := grpc.NewServer()
//         pb.RegisterURLShortenerServer(grpcServer, server)
//         lis, err := net.Listen("tcp", ":0") // Слушаем на произвольном доступном порту
//         if err != nil {
//             t.Fatalf("Failed to listen: %v", err)
//         }
//         grpcServer.Serve(lis)
//     }()

//     // Создаем gRPC клиент
//     conn, err := grpc.Dial(":0", grpc.WithInsecure())
//     if err != nil {
//         t.Fatalf("Did not connect: %v", err)
//     }
//     defer conn.Close()
//     client := pb.NewURLShortenerClient(conn)

//     // Делаем запрос к серверу
//     resp, err := client.ShortenURL(context.Background(), req)
//     assert.NoError(t, err)
//     assert.Equal(t, expectedShortURL, resp.ShortUrl)

//     // Проверяем, что метод GetExistURL был вызван
//     mockService.AssertExpectations(t)
// }
