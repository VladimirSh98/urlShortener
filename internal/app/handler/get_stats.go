package handler

import (
	"encoding/json"
	"github.com/VladimirSh98/urlShortener/internal/app/middleware"
	myProto "github.com/VladimirSh98/urlShortener/proto"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"net/http"
	"strings"
)

// GetStats get stats for urls count and users count
func (h *Handler) GetStats(res http.ResponseWriter, req *http.Request) {
	sugar := zap.S()
	result, err := h.service.GetAllRecords()
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	response := getStatsResponseAPI{
		URLS:  len(result),
		Users: int(middleware.UserCount),
	}
	accept := req.Header.Get("Accept")

	switch {
	case strings.Contains(accept, "application/grpc"):
		grpcResponse := &myProto.StatsResponse{
			Urls:  int32(response.URLS),
			Users: int32(response.Users),
		}
		var grpcData []byte
		grpcData, err = proto.Marshal(grpcResponse)
		if err != nil {
			sugar.Errorln("Failed to marshal gRPC response:", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/grpc+proto")
		res.WriteHeader(http.StatusOK)
		if _, err = res.Write(grpcData); err != nil {
			sugar.Errorln("Failed to write gRPC response:", err)
		}

	default:
		var jsonResponse []byte
		jsonResponse, err = json.Marshal(response)
		if err != nil {
			sugar.Errorln("JSON marshal error:", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		if _, err = res.Write(jsonResponse); err != nil {
			sugar.Errorln("Failed to write JSON response:", err)
		}
	}
}
