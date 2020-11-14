package zeebetask

import (
	"github.com/project-flogo/core/data/coerce"
)

// Settings struct
type Settings struct {
	ZeebeBrokerHost string `md:"zeebeBrokerHost,required"`
	ZeebeBrokerPort int    `md:"zeebeBrokerPort,required"`
	BpmnProcessID   string `md:"bpmnProcessID,required"`
	ServiceType     string `md:"serviceType,required"`
	UsePlainTextConnection bool `md:"usePlainTextConnection"`
}

// FromMap method of Settings
func (s *Settings) FromMap(values map[string]interface{}) error {
	var (
		err             error
		zeebeBrokerHost string
		zeebeBrokerPort int
		bpmnProcessID   string
		serviceType     string
		usePlainTextConnection bool
	)

	zeebeBrokerHost, err = coerce.ToString(values["zeebeBrokerHost"])
	if err != nil {
		return err
	}
	s.ZeebeBrokerHost = zeebeBrokerHost

	zeebeBrokerPort, err = coerce.ToInt(values["zeebeBrokerPort"])
	if err != nil {
		return err
	}
	s.ZeebeBrokerPort = zeebeBrokerPort

	bpmnProcessID, err = coerce.ToString(values["bpmnProcessID"])
	if err != nil {
		return err
	}
	s.BpmnProcessID = bpmnProcessID

	serviceType, err = coerce.ToString(values["serviceType"])
	if err != nil {
		return err
	}
	s.ServiceType = serviceType

	usePlainTextConnection, err = coerce.ToBool(values["usePlainTextConnection"])
	if err != nil {
		return err
	}
	s.UsePlainTextConnection = usePlainTextConnection

	return nil
}

// ToMap method of Settings
func (s *Settings) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"zeebeBrokerHost": s.ZeebeBrokerHost,
		"zeebeBrokerPort": s.ZeebeBrokerPort,
		"bpmnProcessID":   s.BpmnProcessID,
		"serviceType":   	 s.ServiceType,
		"usePlainTextConnection": s.UsePlainTextConnection,
	}
}

// HandlerSettings struct
type HandlerSettings struct{
	JobConcurrency int `md:"jobConcurrency"`
	MaxActiveJobs int `md:"maxActiveJobs"`
	PollInterval string `md:"pollInterval"` 
	PollThreshold float64 `md:"pollThreshold"`
	RequestTimeout string `md:"requestTimeout"`
	Timeout string `md:"timeout"`
}

// FromMap method of HandlerSettings
func (hs *HandlerSettings) FromMap(values map[string]interface{}) error {
	var (
		err  error
		jobConcurrency int
		maxActiveJobs int
		pollInterval string
		pollThreshold float64
		requestTimeout string
		timeout string
	)

	jobConcurrency, err = coerce.ToInt(values["jobConcurrency"])
	if err != nil {
		return err
	}

	maxActiveJobs, err = coerce.ToInt(values["maxActiveJobs"])
	if err != nil {
		return err
	}

	pollInterval, err = coerce.ToString(values["pollInterval"])
	if err != nil {
		return err
	}

	pollThreshold, err = coerce.ToFloat64(values["pollThreshold"])
	if err != nil {
		return err
	}

	requestTimeout, err = coerce.ToString(values["requestTimeout"])
	if err != nil {
		return err
	}

	timeout, err = coerce.ToString(values["timeout"])
	if err != nil {
		return err
	}

	hs.JobConcurrency = jobConcurrency
	hs.MaxActiveJobs = maxActiveJobs
	hs.PollInterval = pollInterval
	hs.PollThreshold = pollThreshold
	hs.RequestTimeout = requestTimeout
	hs.Timeout = timeout
	return nil
}

// ToMap method of HandlerSettings
func (hs *HandlerSettings) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"jobConcurrency": hs.JobConcurrency,
		"maxActiveJobs": hs.MaxActiveJobs,
		"pollInterval": hs.PollInterval,
		"pollThreshold": hs.PollThreshold,
		"requestTimeout": hs.RequestTimeout,
		"timeout": hs.Timeout,
	}
}

// Output struct
type Output struct {
	Status string `md:"status,required"`
	Result interface{} `md:"result,required"`
}

// FromMap method of Output
func (o *Output) FromMap(values map[string]interface{}) error {
	var (
		err  error
		status string
		result map[string]interface{}
	)

	status, err = coerce.ToString(values["status"])
	if err != nil {
		return err
	}

	result, err = coerce.ToObject(values["result"])
	if err != nil {
		return err
	}

	o.Status = status
	o.Result = result
	return nil
}

// ToMap method of Output
func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"status": o.Status,
		"result": o.Result,
	}
}

// Reply struct
type Reply struct {
	Status string      `md:"status,required"`
	Result interface{} `md:"result,required"`
}

// FromMap method of Reply
func (r *Reply) FromMap(values map[string]interface{}) error {
	var (
		err    error
		status string
		result interface{}
	)

	status, err = coerce.ToString(values["status"])
	if err != nil {
		return err
	}
	r.Status = status

	result, err = coerce.ToAny(values["result"])
	if err != nil {
		return err
	}
	r.Result = result

	return nil
}

// ToMap method of Reply
func (r *Reply) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"status": r.Status,
		"result": r.Result,
	}
}
