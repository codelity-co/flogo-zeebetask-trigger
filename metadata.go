package zeebetask

import (
	"github.com/project-flogo/core/data/coerce"
)

// Settings struct
type Settings struct {
	Enabled								 bool   `md:"enabbled,required"`
	ZeebeBrokerHost        string `md:"zeebeBrokerHost,required"`
	ZeebeBrokerPort        int    `md:"zeebeBrokerPort,required"`
	UsePlainTextConnection bool   `md:"usePlainTextConnection"`
	CaCertificatePath      string `md:"caCertificatePath"`
	Token                  string `md:"token"`
	ClientID               string `md:"clientID"`
	ClientSecret           string `md:"clientSecret"`
	AudienceEndpoint       string `md:"audienceEndpoint"`
	AuthorizationServerUrl string `md:"authorizationServerUrl"`
	TimeoutDurationString  string `md:"timeoutDurationString"`
}

// FromMap method of Settings
func (s *Settings) FromMap(values map[string]interface{}) error {
	var (
		err                    error
		enabled                bool
		zeebeBrokerHost        string
		zeebeBrokerPort        int
		usePlainTextConnection bool
		caCertificatePath      string
		token                  string
		clientID               string
		clientSecret           string
		audienceEndpoint       string
		authorizationServerUrl string
		timeoutDurationString  string
	)

	enabled, err = coerce.ToBool(values["enabled"])
	if err != nil {
		return err
	}
	s.Enabled = enabled

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

	usePlainTextConnection, err = coerce.ToBool(values["usePlainTextConnection"])
	if err != nil {
		return err
	}
	s.UsePlainTextConnection = usePlainTextConnection

	caCertificatePath, err = coerce.ToString(values["caCertificatePath"])
	if err != nil {
		return err
	}
	s.CaCertificatePath = caCertificatePath

	token, err = coerce.ToString(values["token"])
	if err != nil {
		return err
	}
	s.Token = token

	clientID, err = coerce.ToString(values["clientID"])
	if err != nil {
		return err
	}
	s.ClientID = clientID

	clientSecret, err = coerce.ToString(values["clientSecret"])
	if err != nil {
		return err
	}
	s.ClientSecret = clientSecret

	audienceEndpoint, err = coerce.ToString(values["audienceEndpoint"])
	if err != nil {
		return err
	}
	s.AudienceEndpoint = audienceEndpoint

	authorizationServerUrl, err = coerce.ToString(values["authorizationServerUrl"])
	if err != nil {
		return err
	}
	s.AuthorizationServerUrl = authorizationServerUrl

	timeoutDurationString, err = coerce.ToString(values["timeoutDurationString"])
	if err != nil {
		return err
	}
	s.TimeoutDurationString = timeoutDurationString

	return nil
}

// ToMap method of Settings
func (s *Settings) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"enabled":                s.Enabled,
		"zeebeBrokerHost":        s.ZeebeBrokerHost,
		"zeebeBrokerPort":        s.ZeebeBrokerPort,
		"usePlainTextConnection": s.UsePlainTextConnection,
		"caCertificatePath":      s.CaCertificatePath,
		"token":                  s.Token,
		"clientID":               s.ClientID,
		"clientSecret":           s.ClientSecret,
		"audienceEndpoint":       s.AudienceEndpoint,
		"authorizationServerUrl": s.AuthorizationServerUrl,
		"timeoutDurationString":  s.TimeoutDurationString,
	}
}

// HandlerSettings struct
type HandlerSettings struct {
	ServiceType                  string  `md:"serviceType,required"`
	JobConcurrency               int     `md:"jobConcurrency"`
	MaxActiveJobs                int     `md:"maxActiveJobs"`
	PollIntervalDurationString   string  `md:"pollIntervalDurationString"`
	PollThreshold                float64 `md:"pollThreshold"`
	RequestTimeoutDurationString string  `md:"requestTimeoutDurationString"`
	TimeoutDurationString        string  `md:"timeoutDurationString"`
}

// FromMap method of HandlerSettings
func (hs *HandlerSettings) FromMap(values map[string]interface{}) error {
	var (
		err                          error
		serviceType                  string
		jobConcurrency               int
		maxActiveJobs                int
		pollIntervalDurationString   string
		pollThreshold                float64
		requestTimeoutDurationString string
		timeoutDurationString        string
	)

	serviceType, err = coerce.ToString(values["serviceType"])
	if err != nil {
		return err
	}

	jobConcurrency, err = coerce.ToInt(values["jobConcurrency"])
	if err != nil {
		return err
	}

	maxActiveJobs, err = coerce.ToInt(values["maxActiveJobs"])
	if err != nil {
		return err
	}

	pollIntervalDurationString, err = coerce.ToString(values["pollIntervalDurationString"])
	if err != nil {
		return err
	}

	pollThreshold, err = coerce.ToFloat64(values["pollThreshold"])
	if err != nil {
		return err
	}

	requestTimeoutDurationString, err = coerce.ToString(values["requestTimeoutDurationString"])
	if err != nil {
		return err
	}

	timeoutDurationString, err = coerce.ToString(values["timeoutDurationString"])
	if err != nil {
		return err
	}

	hs.ServiceType = serviceType
	hs.JobConcurrency = jobConcurrency
	hs.MaxActiveJobs = maxActiveJobs
	hs.PollIntervalDurationString = pollIntervalDurationString
	hs.PollThreshold = pollThreshold
	hs.RequestTimeoutDurationString = requestTimeoutDurationString
	hs.TimeoutDurationString = timeoutDurationString
	return nil
}

// ToMap method of HandlerSettings
func (hs *HandlerSettings) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"serviceType":                  hs.ServiceType,
		"jobConcurrency":               hs.JobConcurrency,
		"maxActiveJobs":                hs.MaxActiveJobs,
		"pollIntervalDurationString":   hs.PollIntervalDurationString,
		"pollThreshold":                hs.PollThreshold,
		"requestTimeoutDurationString": hs.RequestTimeoutDurationString,
		"timeoutDurationString":        hs.TimeoutDurationString,
	}
}

// Output struct
type Output struct {
	JobKey         int64                  `md:"jobKey"`
	Headers        map[string]interface{} `md:"headers"`
	InputVariables map[string]interface{} `md:"inputVarialbes"`
}

// FromMap method of Output
func (o *Output) FromMap(values map[string]interface{}) error {
	var (
		err            error
		jobKey         int64
		headers        map[string]interface{}
		inputVariables map[string]interface{}
	)

	jobKey, err = coerce.ToInt64(values["jobKey"])
	if err != nil {
		return err
	}

	headers, err = coerce.ToObject(values["headers"])
	if err != nil {
		return err
	}

	inputVariables, err = coerce.ToObject(values["inputVariables"])
	if err != nil {
		return err
	}

	o.JobKey = jobKey
	o.Headers = headers
	o.InputVariables = inputVariables
	return nil
}

// ToMap method of Output
func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"jobKey":         o.JobKey,
		"headers":        o.Headers,
		"inputVariables": o.InputVariables,
	}
}

// Reply struct
type Reply struct {
	ApplicationMessageType string                 `md:"applicationMessageType"`
	ApplicationMessageCode string                 `md:"applicationMessageCode"`
	ApplicationMessageText string                 `md:"applicationMessageText"`
	OutputVariables        map[string]interface{} `md:"outputVariables"`
}

// FromMap method of Reply
func (r *Reply) FromMap(values map[string]interface{}) error {
	var (
		err                    error
		applicationMessageType string
		applicationMessageCode string
		applicationMessageText string
		outputVariables        map[string]interface{}
	)

	applicationMessageType, err = coerce.ToString(values["applicationMessageType"])
	if err != nil {
		return err
	}
	r.ApplicationMessageType = applicationMessageType

	applicationMessageCode, err = coerce.ToString(values["applicationMessageCode"])
	if err != nil {
		return err
	}
	r.ApplicationMessageCode = applicationMessageCode

	applicationMessageText, err = coerce.ToString(values["applicationMessageText"])
	if err != nil {
		return err
	}
	r.ApplicationMessageText = applicationMessageText

	outputVariables, err = coerce.ToObject(values["outputVariables"])
	if err != nil {
		return err
	}
	r.OutputVariables = outputVariables

	return nil
}

// ToMap method of Reply
func (r *Reply) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"applicationMessageType": r.ApplicationMessageType,
		"applicationMessageCode": r.ApplicationMessageCode,
		"applicationMessageText": r.ApplicationMessageText,
		"outputVariables":        r.OutputVariables,
	}
}
