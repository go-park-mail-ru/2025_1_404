package logger

type StubLogger struct{}

func NewStub() Logger {
	return &StubLogger{}
}

func (s *StubLogger) WithFields(_ LoggerFields) Logger { return s }
func (s *StubLogger) Info(_ string)                    {}
func (s *StubLogger) Error(_ string)                   {}
func (s *StubLogger) Warn(_ string)                    {}
func (s *StubLogger) Debug(_ string)                   {}
