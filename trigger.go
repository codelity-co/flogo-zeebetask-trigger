package zeebetask

import (
	"fmt"

	"github.com/project-flogo/core/trigger"
	"github.com/zeebe-io/zeebe/clients/go/pkg/zbc"
)

var triggerMd = trigger.NewMetadata(&Settings{}, &HandlerSettings{}, &Output{}, &Reply{})

func init() {
	_ = trigger.Register(&Trigger{}, &Factory{})
}

type (

	// Factory of Trigger
	Factory struct{}

	// Trigger struct
	Trigger struct {
		triggerConfig *trigger.Config
		natsHandlers  []*Handler
	}
)

// New method of Trigger Factory
func (*Factory) New(config *trigger.Config) (trigger.Trigger, error) {
	return &Trigger{triggerConfig: config}, nil
}

// Metadata method of Trigger Factory
func (f *Factory) Metadata() *trigger.Metadata {
	return triggerMd
}

// Metadata implements trigger.Trigger.Metadata
func (t *Trigger) Metadata() *trigger.Metadata {
	return triggerMd
}

// Initialize method of Trigger
func (t *Trigger) Initialize(ctx trigger.InitContext) error {

	var (
		err         error
		zeebeClient zbc.Client
	)

	logger := ctx.Logger()

	s := &Settings{}
	err = s.FromMap(t.triggerConfig.Settings)
	if err != nil {
		return err
	}
	logger.Debugf("Settings: %v", s)

	// Init handlers
	for _, handler := range ctx.GetHandlers() {

		// Create handler settings
		logger.Infof("Mapping handler settings...")
		handlerSettings := &HandlerSettings{}
		if err := handlerSettings.FromMap(handler.Settings()); err != nil {
			return err
		}
		logger.Debugf("handlerSettings: %v", handlerSettings)
		logger.Infof("Mapped handler settings successfully")

		// Connect to Zeebe broker
		zeebeClient, err = zbc.NewClient(&zbc.ClientConfig{
			GatewayAddress:         fmt.Sprintf("%v:%v", s.ZeebeBrokerHost, s.ZeebeBrokerPort),
			UsePlaintextConnection: true,
		})
		if err != nil {
			logger.Errorf("Zeebe broker connection error: %v", err)
			return err
		}

		// Create Stop Channel
		logger.Debugf("Registering trigger handler...")
		stopChannel := make(chan bool)

		// Create Trigger Handler
		natsHandler := &Handler{
			zeebeClient:   zeebeClient,
			bpmnProcessID: s.BpmnProcessID,
			serviceType:   s.ServiceType,
			command:       s.Command,
			stopChannel:   stopChannel,
		}

		// Append handler
		t.natsHandlers = append(t.natsHandlers, natsHandler)
		logger.Debugf("Registered trigger handler successfully")
	}

	return nil
}

// Start implements util.Managed.Start
func (t *Trigger) Start() error {
	for _, handler := range t.natsHandlers {
		_ = handler.Start()
	}
	return nil
}

// Stop implements util.Managed.Stop
func (t *Trigger) Stop() error {
	for _, handler := range t.natsHandlers {
		_ = handler.Stop()
	}
	return nil
}

// Handler is Zeebe task handler
type Handler struct {
	zeebeClient   zbc.Client
	bpmnProcessID string
	serviceType   string
	command       string
	stopChannel   chan bool
}

// Start starts the handler
func (h *Handler) Start() error {
	return nil
}

// Stop implements util.Managed.Stop
func (h *Handler) Stop() error {
	//stop servers/services if necessary
	return nil
}
