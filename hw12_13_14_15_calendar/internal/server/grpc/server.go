package internalgrpc

import (
	"context"
	"net"

	"github.com/Andrewmakmaer/goOtusHM/hw12_13_14_15_calendar/internal/server/grpc/pb"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/kit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type Application interface {
	CreateEvent(context.Context, string, string, string, string, string, string, string) error
	UpdateEvent(context.Context, string, string, string, string, string, string, string) error
	DeleteEvent(string, string) error
	ListEventDay(string) (string, error)
	ListEventWeek(string) (string, error)
	ListEventsMonth(string) (string, error)
}

type Logger interface {
	Debug(...interface{})
	Info(...interface{})
	Warn(...interface{})
	Error(...interface{})
	Log(...interface{}) error
}

type Server struct {
	pb.UnimplementedCalendarServer
	app      Application
	logg     Logger
	listener net.Listener
	serv     *grpc.Server
}

func NewServer(logger Logger, app Application, port string) Server {
	lsn, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logger.Error("message", err.Error())
	}

	serv := grpc.NewServer(grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(kit.UnaryServerInterceptor(logger))))
	newServer := &Server{app: app, logg: logger, listener: lsn, serv: serv}
	return *newServer
}

func (s *Server) Start(ctx context.Context) error {
	pb.RegisterCalendarServer(s.serv, s)
	reflection.Register(s.serv)

	s.logg.Info("message", "starting grpc server")
	if err := s.serv.Serve(s.listener); err != nil {
		s.logg.Error("message", err.Error())
		return err
	}

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.serv.GracefulStop()
	s.listener.Close()
	return nil
}

func (s *Server) AddEvent(ctx context.Context, req *pb.AddEventRequest) (*pb.StatusMessageResponce, error) {
	eventFields := req.Change
	if eventFields == nil {
		return nil, status.Error(codes.InvalidArgument, "error during add event")
	}

	err := s.app.CreateEvent(ctx, eventFields.Id, eventFields.UserId, eventFields.Title,
		eventFields.Description, eventFields.StartTime, eventFields.EndTime, eventFields.CallDuration)
	if err != nil {
		s.logg.Error("message", err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &pb.StatusMessageResponce{Message: "event created"}, nil
}

func (s *Server) UpdateEvent(ctx context.Context, req *pb.UpdateEventRequest) (*pb.StatusMessageResponce, error) {
	eventFields := req.Change
	err := s.app.UpdateEvent(ctx, eventFields.Id, eventFields.UserId, eventFields.Title,
		eventFields.Description, eventFields.StartTime, eventFields.EndTime, eventFields.CallDuration)
	if err != nil {
		s.logg.Error("message", err.Error())
		return nil, err
	}
	return &pb.StatusMessageResponce{Message: "Success updated event"}, nil
}

func (s *Server) DeleteEvent(ctx context.Context, req *pb.DeleteEventRequest) (*pb.StatusMessageResponce, error) {
	userID := req.UserId
	id := req.Id

	err := s.app.DeleteEvent(userID, id)
	if err != nil {
		s.logg.Error("message", err.Error(), "id", id, "userID", userID)
		return nil, err
	}

	s.logg.Info("message", "Event successfully deleted")
	return &pb.StatusMessageResponce{Message: "Success delete event"}, nil
}

func (s *Server) ListEvent(ctx context.Context, req *pb.ListEventRequest) (*pb.EventsListResponce, error) {
	userID := req.UserId
	timeBy := req.TimeBy
	if userID == "" || timeBy == "" {
		return nil, status.Error(codes.InvalidArgument, "error during list event, undefined UserID or TimeBy")
	}

	var err error
	var result string
	switch timeBy {
	case "day":
		result, err = s.app.ListEventDay(userID)
	case "week":
		result, err = s.app.ListEventWeek(userID)
	case "month":
		result, err = s.app.ListEventsMonth(userID)
	}

	if err != nil {
		s.logg.Error("message", err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &pb.EventsListResponce{Message: result}, nil
}
