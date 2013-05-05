package rinst

const FlushBufferSize = 100

func FlushSchema(sb SchemaBuffer) []Schema {

	...
	
	buff := make([]Schema, 0, FlushBufferSize)
	smallbuff := make([]Schema, 0, FlushBufferSize)
	i := 0

	for {
		s, ok := <-sb
		if !ok {
			break
		}
	}

	buff = append(buff, smallbuff...)
	return buff
}

... FlushMeasurements()