package micro

import (
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/config/source"
)

// Service represents a microservice definition.
type Service struct {
	micro.Service
	Name string
}

// NewService initializes a new microservice instance with the given identifier
// (e.g. com.example.svc.greeter), and given a configuration struct containing
// default values, loads configuration values from the given sources into the
// struct.
func NewService(name string, conf interface{}, sources ...source.Source) (*Service, error) {
	service := micro.NewService(
		micro.Name(name),
		micro.WrapClient(
			requestIDClientWrapper,
			logClientWrapper,
		),
		micro.WrapHandler(
			requestIDHandlerWrapper,
			logHandlerWrapper,
		),
	)
	service.Init()

	//log.DefaultLogger = zerolog.NewLogger(
	//	zerolog.UseAsDefault(),
	//	log.WithFields(map[string]interface{}{
	//		"node":    service.Server().Options().Id,
	//		"service": name,
	//	}),
	//	zerolog.WithTimeFormat(time.RFC3339Nano),
	//)

	if conf != nil {
		if err := config.Load(sources...); err != nil {
			return nil, err
		}

		if err := config.Scan(conf); err != nil {
			return nil, err
		}
	}

	return &Service{service, name}, nil
}
