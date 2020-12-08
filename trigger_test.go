package zeebetask

import (
	"testing"

	"github.com/project-flogo/core/support"
	"github.com/project-flogo/core/trigger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	// "github.com/zeebe-io/zeebe/clients/go/pkg/zbc"
)

type ZeebetaskTriggerTestSuite struct {
	suite.Suite
}

func TestZeebetaskTriggerTestSuite(t *testing.T) {
	suite.Run(t, new(ZeebetaskTriggerTestSuite))
}

func (suite *ZeebetaskTriggerTestSuite) TestZeebetaskTrigger_Register() {

	ref := support.GetRef(&Trigger{})
	f := trigger.GetFactory(ref)
	assert.NotNil(suite.T(), f)
}
