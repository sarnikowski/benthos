package pipeline

import (
	"fmt"
	"strconv"

	iprocessor "github.com/Jeffail/benthos/v3/internal/component/processor"
	"github.com/Jeffail/benthos/v3/internal/interop"
	"github.com/Jeffail/benthos/v3/internal/old/processor"
)

// Config is a configuration struct for creating parallel processing pipelines.
// The number of resuling parallel processing pipelines will match the number of
// threads specified. Processors are executed on each message in the order that
// they are written.
//
// In order to fully utilise each processing thread you must either have a
// number of parallel inputs that matches or surpasses the number of pipeline
// threads, or use a memory buffer.
type Config struct {
	Threads    int                `json:"threads" yaml:"threads"`
	Processors []processor.Config `json:"processors" yaml:"processors"`
}

// NewConfig returns a configuration struct fully populated with default values.
func NewConfig() Config {
	return Config{
		Threads:    -1,
		Processors: []processor.Config{},
	}
}

//------------------------------------------------------------------------------

// New creates an input type based on an input configuration.
func New(conf Config, mgr interop.Manager) (Type, error) {
	processors := make([]iprocessor.V1, len(conf.Processors))
	for j, procConf := range conf.Processors {
		var err error
		pMgr := mgr.IntoPath("processors", strconv.Itoa(j))
		processors[j], err = processor.New(procConf, pMgr, pMgr.Logger(), pMgr.Metrics())
		if err != nil {
			return nil, fmt.Errorf("failed to create processor '%v': %v", procConf.Type, err)
		}
	}
	if conf.Threads == 1 {
		return NewProcessor(processors...), nil
	}
	return newPoolV2(conf.Threads, mgr.Logger(), processors...)
}