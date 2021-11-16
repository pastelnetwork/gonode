package messaging

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/gogo/protobuf/proto"
	"github.com/pastelnetwork/gonode/common/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Remoter struct {
	remoter *remote.Remote
	context *actor.RootContext
	mapPID  map[string]*actor.PID
}

func (r *Remoter) Send(ctx context.Context, pid *actor.PID, message proto.Message) error {
	if actorContext := ctx.GetActorContext(); actorContext != nil {
		actorContext.Send(pid, message)
	} else {
		r.context.Send(pid, message)
	}
	return nil
}

func (r *Remoter) RegisterActor(a actor.Actor, name string) (*actor.PID, error) {
	pid, err := r.context.SpawnNamed(actor.PropsFromProducer(func() actor.Actor { return a }), name)
	if err != nil {
		return nil, err
	}
	r.mapPID[name] = pid
	return pid, nil
}

func (r *Remoter) DeregisterActor(name string) {
	if r.mapPID[name] != nil {
		r.context.Stop(r.mapPID[name])
	}
	delete(r.mapPID, name)
}

func (r *Remoter) Start() {
	r.remoter.Start()
}

func (r *Remoter) GracefulStop() {
	for name, pid := range r.mapPID {
		r.context.StopFuture(pid).Wait()
		delete(r.mapPID, name)
	}
}

type Config struct {
	Host              string                           `mapstructure:"host"`
	Port              int                              `mapstructure:"port"`
	clientSecureCreds credentials.TransportCredentials `mapstructure:"-"`
	serverSecureCreds credentials.TransportCredentials `mapstructure:"-"`
}

func (c *Config) WithClientSecureCreds(s credentials.TransportCredentials) *Config {
	c.clientSecureCreds = s
	return c
}

func (c *Config) WithServerSecureCreds(s credentials.TransportCredentials) *Config {
	c.serverSecureCreds = s
	return c
}

func NewRemoter(system *actor.ActorSystem, cfg Config) *Remoter {
	if cfg.Host == "" {
		cfg.Host = "0.0.0.0"
	}
	if cfg.Port == 0 {
		cfg.Port = 9000
	}
	remoterConfig := remote.Configure(cfg.Host, cfg.Port)
	clientCreds := []grpc.DialOption{grpc.WithBlock()}
	serverCreds := []grpc.ServerOption{}
	if cfg.clientSecureCreds != nil {
		clientCreds = append(clientCreds, grpc.WithTransportCredentials(cfg.clientSecureCreds))
	}
	if cfg.serverSecureCreds != nil {
		serverCreds = append(serverCreds, grpc.Creds(cfg.serverSecureCreds))
	}
	remoterConfig = remoterConfig.WithServerOptions(serverCreds...).WithDialOptions(clientCreds...)
	return &Remoter{
		remoter: remote.NewRemote(system, remoterConfig),
		context: system.Root,
		mapPID:  make(map[string]*actor.PID),
	}
}
