package zeebetask

import (
	"context"
	"fmt"
	"time"

	"github.com/project-flogo/core/trigger"
	"github.com/zeebe-io/zeebe/clients/go/pkg/entities"
	"github.com/zeebe-io/zeebe/clients/go/pkg/worker"
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
		triggerInitContext trigger.InitContext
		zeebeHandlers  []*Handler
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

	t.triggerInitContext = ctx
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
			UsePlaintextConnection: s.UsePlainTextConnection,
		})
		if err != nil {
			logger.Errorf("Zeebe broker connection error: %v", err)
			return err
		}

		// Create Stop Channel
		logger.Debugf("Registering trigger handler...")
		stopChannel := make(chan bool)

		// Create Trigger Handler
		zeebeHandler := &Handler{
			triggerInitContext: ctx,
			zeebeClient: zeebeClient,
			bpmnProcessID: s.BpmnProcessID,
			serviceType:   s.ServiceType,
			stopChannel:   stopChannel,
			triggerHandlerSettings: handlerSettings,
			triggerHandler: handler,
		}

		// Append handler
		t.zeebeHandlers = append(t.zeebeHandlers, zeebeHandler)
		logger.Debugf("Registered trigger handler successfully")
	}

	return nil
}

// Start implements util.Managed.Start
func (t *Trigger) Start() error {
	for _, handler := range t.zeebeHandlers {
		_ = handler.Start()
	}
	return nil
}

// Stop implements util.Managed.Stop
func (t *Trigger) Stop() error {
	var err error
	for _, handler := range t.zeebeHandlers {
		err = handler.Stop()
		if err != nil {
			t.triggerInitContext.Logger().Errorf("Trigger stop error: %v", err)
			return err
		}
		handler.zeebeClient.Close()
	}
	return nil
}

// Handler is Zeebe task handler
type Handler struct {
	triggerInitContext trigger.InitContext
	zeebeClient   zbc.Client
	bpmnProcessID string
	serviceType   string
	stopChannel   chan bool
	jobWorker 		worker.JobWorker
	triggerHandlerSettings *HandlerSettings
	triggerHandler trigger.Handler
}

// Start starts the handler
func (h *Handler) Start() error {

	step3 := h.zeebeClient.NewJobWorker().JobType(h.serviceType).Handler(h.handleJob)

	if h.triggerHandlerSettings.JobConcurrency > 0 {
		step3 = step3.Concurrency(h.triggerHandlerSettings.JobConcurrency)
	}

	if h.triggerHandlerSettings.MaxActiveJobs > 0 {
		step3 = step3.MaxJobsActive(h.triggerHandlerSettings.JobConcurrency)
	}

	if h.triggerHandlerSettings.PollInterval != "" {

		pollInterval, err := time.ParseDuration(h.triggerHandlerSettings.PollInterval)
		if err != nil {
			return err
		}

		step3 = step3.PollInterval(pollInterval)

	}

	if h.triggerHandlerSettings.PollThreshold > 0.0 {
		step3 = step3.PollThreshold(h.triggerHandlerSettings.PollThreshold)
	}

	if h.triggerHandlerSettings.RequestTimeout != ""  {
		requestTimeout, err := time.ParseDuration(h.triggerHandlerSettings.RequestTimeout)
		if err != nil {
			return err
		}

		step3 = step3.RequestTimeout(requestTimeout)
	}

	if h.triggerHandlerSettings.Timeout != ""  {
		timeout, err := time.ParseDuration(h.triggerHandlerSettings.Timeout)
		if err != nil {
			return err
		}

		step3 = step3.Timeout(timeout)
	}

	h.jobWorker = step3.Open()

	return nil
}

// Stop implements util.Managed.Stop
func (h *Handler) Stop() error {
	h.stopChannel <- true
	//stop servers/services if necessary
	h.jobWorker.Close()
	return nil
}

func (h *Handler) handleJob(client worker.JobClient, job entities.Job) {
	logger := h.triggerInitContext.Logger()

	jobKey := job.GetKey()
	logger.Debug("JobKey: ", jobKey)

	headers, err := job.GetCustomHeadersAsMap()
	if err != nil {
		logger.Errorf("job.GetCustomerHeadersAsMap error: %v", err)
		failJob(client, job, err)
		return
	}
	logger.Debug("Headers: ", headers)

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		logger.Errorf("job.GetVariablesAsMap error: %v", err)
		failJob(client, job, err)
		return
	}
	logger.Debug("Variables: ", variables)

	output := &Output{}
	result, err := h.triggerHandler.Handle(context.Background(), output.ToMap())
	if err != nil {
		logger.Errorf("triggerHandler.Handle error: %v", err)
		failJob(client, job, err)
		return
	}

	reply := &Reply{}
	err = reply.FromMap(result)
	if err != nil {
		logger.Errorf("Parsing reply error: %v", err)
		failJob(client, job, err)
		return
	}
	
	request, err := client.NewCompleteJobCommand().JobKey(jobKey).VariablesFromMap(result)
	if err != nil {
		// failed to set the updated variables
		failJob(client, job, err)
		return
	}

	_, err = request.Send(context.Background())
	if err != nil {
		logger.Errorf("Complete job request send error: %v", err)
		failJob(client, job, err)
		return
	}

	logger.Info("Complete job", jobKey, "of type", job.Type)
}

func failJob(client worker.JobClient, job entities.Job, err error) {
	request := client.NewFailJobCommand().JobKey(job.GetKey()).Retries(job.Retries - 1).ErrorMessage(err.Error())
	_, _ = request.Send(context.Background())
}