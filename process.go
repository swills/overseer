package overseer

import (
	"fmt"
	"time"

	"github.com/ShinyTrinkets/go-cmd"
)

const (
	idLength          uint = 16
	defaultDelayStart uint = 25
	defaultRetryTimes uint = 3
)

// ChildProcess structure
type ChildProcess struct {
	cmd.Cmd
	id         string // The id is private
	DelayStart uint   // Nr of milli-seconds to delay the start
	RetryTimes uint   // Nr of times to restart on failure
}

// JSONProcess structure
type JSONProcess struct {
	Cmd        string    `json:"cmd"`
	PID        int       `json:"PID"`
	Complete   bool      `json:"complete"` // false if stopped or signaled
	ExitCode   int       `json:"exitCode"` // exit code of process
	Error      error     `json:"error"`    // Go error
	RunTime    float64   `json:"runTime"`  // seconds, zero if Cmd not started
	StartTime  time.Time `json:"startTime"`
	Env        []string  `json:"end"`
	Dir        string    `json:"dir"`
	DelayStart uint      `json:"delayStart"`
	RetryTimes uint      `json:"retryTimes"`
}

// NewChild returns a new child process for the given command name and arguments.
func NewChild(name string, args ...string) *ChildProcess {
	c := cmd.NewCmdOptions(
		cmd.Options{Buffered: false, Streaming: true},
		name,
		args...,
	)
	randID := randomKey(idLength)
	return &ChildProcess{*c, randID, defaultDelayStart, defaultRetryTimes}
}

// CloneChild clones a child process. All the configs are transferred,
// and the state of the original object is lost.
func (c *ChildProcess) CloneChild() *ChildProcess {
	clone := cmd.NewCmdOptions(
		cmd.Options{Buffered: false, Streaming: true},
		c.Name,
		c.Args...,
	)
	randID := randomKey(idLength)
	p := &ChildProcess{*clone, randID, defaultDelayStart, defaultRetryTimes}
	// transfer the config
	p.SetDir(c.Dir)
	p.SetEnv(c.Env)
	p.SetDelayStart(c.DelayStart)
	p.SetRetryTimes(c.RetryTimes)
	return p
}

// ToJSON returns more detailed info about a child process.
func (c *ChildProcess) ToJSON() JSONProcess {
	s := c.Status()
	cmd := fmt.Sprint(c.Name, " ", c.Args)
	startTime := time.Unix(0, s.StartTs)
	return JSONProcess{
		cmd,
		s.PID,
		s.Complete,
		s.Exit,
		s.Error,
		s.Runtime,
		startTime,
		c.Env,
		c.Dir,
		c.DelayStart,
		c.RetryTimes,
	}
}

// SetDir sets the environment variables for the launched process.
// This can only have effect before starting the process.
func (c *ChildProcess) SetDir(dir string) {
	c.Lock()
	defer c.Unlock()
	c.Dir = dir
}

// SetEnv sets the working directory of the command.
// This can only have effect before starting the process.
func (c *ChildProcess) SetEnv(env []string) {
	c.Lock()
	defer c.Unlock()
	c.Env = env
}

// SetDelayStart sets the delay start in milli-seconds.
func (c *ChildProcess) SetDelayStart(delayStart uint) {
	c.Lock()
	defer c.Unlock()
	c.DelayStart = delayStart
}

// SetRetryTimes sets the times of restart in case of failure.
func (c *ChildProcess) SetRetryTimes(retryTimes uint) {
	c.Lock()
	defer c.Unlock()
	c.RetryTimes = retryTimes
}
