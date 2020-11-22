package commands

import (
	"github.com/jfrog/jfrog-cli-core/plugins/components"
	"github.com/magiconair/properties/assert"
	"testing"
	"time"
)

func TestSummaryConfigDefaultValues(t *testing.T) {
	ctx := &components.Context{}
	conf, _ := prepareSummaryConf(ctx)
	assert.Equal(t, conf.refreshRate, 2)
	assert.Equal(t, conf.recalculateRate, 2)
}

func TestShouldUpdateView(t *testing.T) {
	conf := summaryConfiguration{refreshRate: 2}
	lastUpdated := time.Unix(0, 0)
	assert.Equal(t, shouldUpdateView(&conf, lastUpdated), true)

	lastUpdated = time.Now()
	assert.Equal(t, shouldUpdateView(&conf, lastUpdated), false)
}

func TestShouldRecalculate(t *testing.T) {
	conf := summaryConfiguration{recalculateRate: 0}
	lastUpdated := time.Unix(0, 0)
	assert.Equal(t, shouldRecalculate(&conf, lastUpdated), false)

	conf = summaryConfiguration{recalculateRate: 2}
	assert.Equal(t, shouldRecalculate(&conf, lastUpdated), true)

	lastUpdated = time.Now()
	assert.Equal(t, shouldRecalculate(&conf, lastUpdated), false)
}
