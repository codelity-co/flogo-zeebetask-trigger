package zeebetask

import (
	"context"
	"errors"
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

//Factory of Trigger
type Factory struct{}

// New method of Trigger Factory
func (*Factory) New(config *trigger.Config) (trigger.Trigger, error) {
	return &Trigger{triggerConfig: config}, nil
}

// Metadata method of Trigger Factory
func (f *Factory) Metadata() *trigger.Metadata {
	return triggerMd
}

// Trigger struct
type Trigger struct {
	triggerConfig      *trigger.Config
	triggerInitContext trigger.InitContext
	triggerSettings    *Settings
	zeebeHandlers      []*Handler
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

	t.triggerSettings = s

	if !s.Enabled {
		return errors.New("Zeebe trigger is disabled")
	}

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

		zeebeClientConfig := &zbc.ClientConfig{
			GatewayAddress:         fmt.Sprintf("%v:%v", s.ZeebeBrokerHost, s.ZeebeBrokerPort),
			UsePlaintextConnection: s.UsePlainTextConnection,
		}

		//TODO: add credential provider

		// Connect to Zeebe broker
		zeebeClient, err = zbc.NewClient(zeebeClientConfig)
		if err != nil {
			logger.Errorf("Zeebe broker connection error: %v", err)
			return err
		}

		// Create Stop Channel
		logger.Debugf("Registering trigger handler...")

		// Create Trigger Handler
		zeebeHandler := &Handler{
			triggerInitContext:     ctx,
			zeebeClient:            zeebeClient,
			triggerHandlerSettings: handlerSettings,
			triggerHandler:         handler,
		}

		// Append handler
		t.zeebeHandlers = append(t.zeebeHandlers, zeebeHandler)
		logger.Debugf("Registered trigger handler successfully")
	}

	return nil
}

// Start implements util.Managed.Start
func (t *Trigger) Start() error {

	logger := t.triggerInitContext.Logger()

	logger.Debugf("t.triggerSettings.Enabled: %v", t.triggerSettings.Enabled)
	if !t.triggerSettings.Enabled {
		err := errors.New("Trigger is disabled")
		logger.Info(err)
		return err
	}

	for _, handler := range t.zeebeHandlers {
		err := handler.Start()
		if err != nil {
			logger.Error(err)
			return err
		}
	}

	return nil
}

// Stop implements util.Managed.Stop
func (t *Trigger) Stop() error {
	var err error

	for _, handler := range t.zeebeHandlers {
		handler.jobWorker.Close()
		handler.jobWorker.AwaitClose()
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
	triggerInitContext     trigger.InitContext
	zeebeClient            zbc.Client
	jobWorker              worker.JobWorker
	triggerHandlerSettings *HandlerSettings
	triggerHandler         trigger.Handler
}

// Start starts the handler
func (h *Handler) Start() error {

	logger := h.triggerInitContext.Logger()

	logger.Debug("Handler starting...")

	step3 := h.zeebeClient.NewJobWorker().JobType(h.triggerHandlerSettings.ServiceType).Handler(h.handleJob)
	logger.Debug("Zeebe handler has been created")

	if h.triggerHandlerSettings.JobConcurrency > 0 {
		logger.Debugf("Setting Job Concurrency: %v", h.triggerHandlerSettings.JobConcurrency)
		step3 = step3.Concurrency(h.triggerHandlerSettings.JobConcurrency)
	}

	if h.triggerHandlerSettings.MaxActiveJobs > 0 {
		logger.Debugf("Setting Max Active Jobs: %v", h.triggerHandlerSettings.MaxActiveJobs)
		step3 = step3.MaxJobsActive(h.triggerHandlerSettings.JobConcurrency)
	}

	if h.triggerHandlerSettings.PollIntervalDurationString != "" {

		logger.Debugf("Setting Poll Interval: %v", h.triggerHandlerSettings.PollIntervalDurationString)
		pollInterval, err := time.ParseDuration(h.triggerHandlerSettings.PollIntervalDurationString)
		if err != nil {
			logger.Error(err)
			return err
		}

		step3 = step3.PollInterval(pollInterval)

	}

	if h.triggerHandlerSettings.PollThreshold > 0.0 {
		logger.Debugf("Setting Poll Threshold: %v", h.triggerHandlerSettings.PollThreshold)
		step3 = step3.PollThreshold(h.triggerHandlerSettings.PollThreshold)
	}

	if h.triggerHandlerSettings.RequestTimeoutDurationString != "" {
		logger.Debugf("Setting Request Timeout: %v", h.triggerHandlerSettings.RequestTimeoutDurationString)
		requestTimeout, err := time.ParseDuration(h.triggerHandlerSettings.RequestTimeoutDurationString)
		if err != nil {
			logger.Error(err)
			return err
		}

		step3 = step3.RequestTimeout(requestTimeout)
	}

	if h.triggerHandlerSettings.TimeoutDurationString != "" {
		logger.Debugf("Setting Timeout: %v", h.triggerHandlerSettings.TimeoutDurationString)
		timeout, err := time.ParseDuration(h.triggerHandlerSettings.TimeoutDurationString)
		if err != nil {
			logger.Error(err)
			return err
		}

		step3 = step3.Timeout(timeout)
	}

	logger.Debug("Opening jobWorker")
	h.jobWorker = step3.Open()
	logger.Debug("Opened jobWorker")

	return nil
}

// Stop implements util.Managed.Stop
func (h *Handler) Stop() error {
	logger := h.triggerInitContext.Logger()

	logger.Debug("Stopping hanlder...")

	//stop servers/services if necessary
	h.jobWorker.Close()
	logger.Debug("Handler has been stopped")

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
