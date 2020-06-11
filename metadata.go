package zeebetask

import (
	"github.com/project-flogo/core/data/coerce"
)

type (

	// Settings struct
	Settings struct {
		ZeebeBrokerHost string `md:"zeebeBrokerHost,required"`
		ZeebeBrokerPort int    `md:"zeebeBrokerPort,required"`
		BpmnProcessID   string `md:"bpmnProcessID,required"`
		ServiceType     string `md:"serviceType,required"`
		Command         string `md:"command,required"`
		UsePlainTextConnection bool `md:"usePlainTextConnection"`
	}

	// HandlerSettings struct
	HandlerSettings struct{}

	// Output struct
	Output struct {
		Data map[string]interface{} `md:"data,required"`
	}

	// Reply struct
	Reply struct {
		Status string      `md:"status,required"`
		Result interface{} `md:"result,required"`
	}
)

// FromMap method of HandlerSettings
func (hs *HandlerSettings) FromMap(values map[string]interface{}) error {
	return nil
}

// ToMap method of HandlerSettings
func (hs *HandlerSettings) ToMap() map[string]interface{} {
	return make(map[string]interface{})
}

// FromMap method of Settings
func (s *Settings) FromMap(values map[string]interface{}) error {
	var (
		err             error
		zeebeBrokerHost string
		zeebeBrokerPort int
		bpmnProcessID   string
		serviceType     string
		command         string
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

	command, err = coerce.ToString(values["command"])
	if err != nil {
		return err
	}
	s.Command = command

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
		"command":         s.Command,
		"usePlainTextConnection": s.UsePlainTextConnection,
	}
}

// FromMap method of Output
func (o *Output) FromMap(values map[string]interface{}) error {
	var (
		err  error
		data map[string]interface{}
	)

	data, err = coerce.ToObject(values["data"])
	if err != nil {
		return err
	}

	o.Data = data
	return nil
}

// ToMap method of Output
func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"data": o.Data,
	}
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
