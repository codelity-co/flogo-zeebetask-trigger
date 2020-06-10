package zeebetask

import (
	"context"
	"fmt"
	"os/exec"
	"testing"
	"time"

	"github.com/project-flogo/core/action"
	"github.com/project-flogo/core/support"
	"github.com/project-flogo/core/support/test"
	"github.com/project-flogo/core/trigger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/zeebe-io/zeebe/clients/go/pkg/zbc"
)

type ZeebetaskTriggerTestSuite struct {
	suite.Suite
}

func TestZeebetaskTriggerTestSuite(t *testing.T) {
	suite.Run(t, new(ZeebetaskTriggerTestSuite))
}

func (suite *ZeebetaskTriggerTestSuite) SetupSuite() {
	command := exec.Command("docker", "start", "zeebe")
	err := command.Run()
	if err != nil {
		fmt.Println(err.Error())
		command := exec.Command("docker", "run", "-p", "26500-26502:26500-26502", "--name", "zeebe", "-d", "camunda/zeebe:latest")
		err := command.Run()
		if err != nil {
			fmt.Println(err.Error())
			panic(err)
		}
		time.Sleep(10 * time.Second)
	}

}

func (suite *ZeebetaskTriggerTestSuite) BeforeTest(suiteName, testName string) {

	switch testName {
	case "TestZeebetaskTrigger_CreateWorkflowInstance":

		zeebeClient, err := zbc.NewClient(&zbc.ClientConfig{
			GatewayAddress:         "localhost:26500",
			UsePlaintextConnection: true,
		})
		if err != nil {
			panic(err)
		}
		response, err := zeebeClient.NewDeployWorkflowCommand().AddResourceFile("./test/order-process.bpmn").Send(context.Background())
		if err != nil {
			panic(err)
		}
		fmt.Println(fmt.Sprintf("response text: %v", response.String()))
	}

}

func (suite *ZeebetaskTriggerTestSuite) TestZeebetaskTrigger_Register() {

	ref := support.GetRef(&Trigger{})
	f := trigger.GetFactory(ref)
	assert.NotNil(suite.T(), f)
}

func (suite *ZeebetaskTriggerTestSuite) TestZeebetaskTrigger_Initialize() {
	t := suite.T()

	factory := &Factory{}
	config := &trigger.Config{}
	actions := map[string]action.Action{"dummy": test.NewDummyAction(func() {
	})}

	trg, err := test.InitTrigger(factory, config, actions)
	assert.Nil(t, err)
	assert.NotNil(t, trg)

	err = trg.Start()
	assert.Nil(t, err)

	err = trg.Stop()
	assert.Nil(t, err)
}
