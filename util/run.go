package util

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	loggerUtil "github.com/xvbnm48/go-api-catatan/util/logger"
)

// RegisterHTTPHandler register endpoint to http server
type RegisterHTTPHandler func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error

// HTTPMiddleware is middleware for http handler
type HTTPMiddleware func(handler http.Handler) http.Handler

// HTTPOption are settings for http server
type HTTPOption struct {
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func HttpErrorHandler(ctx context.Context, mux *runtime.ServeMux, m runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Set("Content-Type", "application/json")

	s, ok := status.FromError(err)
	if !ok {
		s = status.New(codes.Unknown, err.Error())
	}

	errMsg := strings.Split(s.Message(), ":")
	statusCode := strconv.Itoa(int(s.Code()))
	message := s.Message()
	if len(errMsg) == 2 {
		statusCode = errMsg[0]
		message = errMsg[1]
	}
	errorMsg := InfoMessage{
		Code:    statusCode,
		Message: message,
	}
	errMsgs := []InfoMessage{}
	errMsgs = append(errMsgs, errorMsg)

	resp := Response{
		Status:   "error",
		Messages: &errMsgs,
	}
	bs, _ := json.Marshal(&resp)
	w.Write(bs)
}

func ServeGRPCAndHTTP(address, port string, grpcServer *grpc.Server,
	register RegisterHTTPHandler, creds credentials.TransportCredentials,
	cert tls.Certificate, logger log.Logger,
	option HTTPOption, handlers ...HTTPMiddleware) {

	//rgbb = getBBPattern()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	var tlsConfig *tls.Config
	if creds != nil {
		opts = []grpc.DialOption{grpc.WithTransportCredentials(creds)}
		tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
			NextProtos:   []string{"h2"},
		}
	}

	ctx := context.Background()
	mux := runtime.NewServeMux(
		runtime.WithForwardResponseOption(HttpSuccessHandler),
		runtime.WithMarshalerOption("*", &EmptyMarshaler{}),
		runtime.WithProtoErrorHandler(HttpErrorHandler),
	)

	err := register(ctx, mux, address, opts)
	if err != nil {
		logger.Log(loggerUtil.LogError, err.Error())
		return
	}

	var handler http.Handler
	handler = mux
	for _, hm := range handlers {
		handler = hm(handler)
	}

	var httpServer *http.Server

	if creds != nil {
		httpServer = &http.Server{
			Addr:      address,
			Handler:   grpcHandlerFunc(grpcServer, handler),
			TLSConfig: tlsConfig,
		}
	} else {
		httpServer = &http.Server{
			Addr:    address,
			Handler: handler,
		}
	}

	if option.ReadTimeout > 0 {
		httpServer.ReadTimeout = option.ReadTimeout
	}

	if option.WriteTimeout > 0 {
		httpServer.WriteTimeout = option.WriteTimeout
	}

	conn, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		logger.Log(loggerUtil.LogError, err.Error())
		return
	}

	if creds != nil {
		err = httpServer.Serve(tls.NewListener(conn, httpServer.TLSConfig))
		if err != nil {
			logger.Log(loggerUtil.LogError, err.Error())
			return
		}
	} else {
		mux := cmux.New(conn)
		grpcL := mux.Match(cmux.HTTP2(), cmux.HTTP2HeaderField("content-type", "application/grpc"))
		httpL := mux.Match(cmux.HTTP1Fast())

		// Use the muxed listeners for your servers.
		go grpcServer.Serve(grpcL)
		go httpServer.Serve(httpL)
		// Start serving!
		mux.Serve()
	}

}

func HttpSuccessHandler(ctx context.Context, w http.ResponseWriter, p proto.Message) error {
	rsp := &Response{
		Status: "success",
		Data:   p,
	}
	rsp.Data = p
	buf, _ := json.Marshal(rsp)
	w.Write(buf)
	return nil
}

func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	})
}
