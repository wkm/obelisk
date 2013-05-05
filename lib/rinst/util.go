package rinst

// extract the schema from an instrument
func FlushSchema(i Instrument, flushSize int) []Schema {
	cb := make(SchemaBuffer, flushSize)
	go func() {
		i.Schema("", cb)
		close(cb)
	}()

	return FlushSchemaBuffer(cb, flushSize)
}

// extract the schema from a SchemaBuffer
func FlushSchemaBuffer(sb SchemaBuffer, flushSize int) []Schema {
	buff := make([]Schema, 0, flushSize)
	smallbuff := make([]Schema, flushSize, flushSize)
	i := 0
	for {
		s, ok := <-sb
		if !ok {
			break
		}

		smallbuff[i%flushSize] = s
		i++

		if i%flushSize == 0 {
			buff = append(buff, smallbuff...)
		}
	}

	buff = append(buff, smallbuff[:i%flushSize]...)

	return buff
}

func FlushMeasurements(i Instrument, flushSize int) []Measurement {
	cb := make(MeasurementBuffer, flushSize)
	go func() {
		i.Measure("", cb)
		close(cb)
	}()

	return FlushMeasurementsBuffer(cb, flushSize)
}

func FlushMeasurementsBuffer(mb MeasurementBuffer, flushSize int) []Measurement {
	buff := make([]Measurement, 0, flushSize)
	smallbuff := make([]Measurement, flushSize, flushSize)
	i := 0
	for {
		m, ok := <-mb
		if !ok {
			break
		}

		smallbuff[i%flushSize] = m
		i++

		if i%flushSize == 0 {
			buff = append(buff, smallbuff...)
		}
	}

	buff = append(buff, smallbuff[:i%flushSize]...)

	return buff
}
