package pipeline

import (
	"errors"
	"strings"
)

var (
	_ Status = (*statusImpl)(nil)
)

// Default implementation of Status
type statusImpl struct {
	code         StatusCode
	err          error
	reasons      []string
	failedPlugin Plugin
	failedStage  string
}

func NewStatus(code StatusCode, reasons ...string) Status {
	status := statusImpl{
		code:    code,
		reasons: reasons,
	}
	if code == InternalError {
		status.err = errors.New(status.Message())
	}
	return &status
}

func NewSuccessStatus() Status {
	return NewStatus(Success)
}

func NewInternalErrorStatus(err error) Status {
	status := statusImpl{
		code:    InternalError,
		err:     err,
		reasons: []string{err.Error()},
	}
	return &status
}

func (s *statusImpl) Code() StatusCode {
	return s.code
}

func (s *statusImpl) CodeAsString() string {
	return statusCodeStrings[s.code]
}

func (s *statusImpl) Error() error {
	return s.err
}

func (s *statusImpl) FailedPlugin() Plugin {
	return s.failedPlugin
}

func (s *statusImpl) FailedStage() string {
	return s.failedStage
}

func (s *statusImpl) SetFailedPlugin(plugin Plugin, stage string) {
	s.failedPlugin = plugin
	s.failedStage = stage
}

func (s *statusImpl) Message() string {
	var msg string
	if s.failedPlugin != nil {
		msg = s.failedPlugin.Name()
	}
	if len(s.failedStage) > 0 {
		msg = msg + " @ " + s.failedStage + ": "
	}
	if s.reasons != nil {
		msg = msg + strings.Join(s.reasons, ", ")
	}
	return msg
}

func (s *statusImpl) Reasons() []string {
	return s.reasons
}
