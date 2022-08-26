package iot

type Converter interface {
	Convert(in RawMessage) (MetricMessage, error)
}
