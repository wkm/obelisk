/*
	A library for computing approximate percentiles from streaming data.
	Based on the algorithms presented in A Fast Algorithm for Approximate
	Quantiles in High Speed Data Streams, Q. Zhang, et.al (2007) and
	Space-Efficient Online Computation of Quantile Summaries, M. Greenwald, et.al
	(2001).

	Briefly, this package maintains a sampled subset of a stream of
	floating point numbers, letting you derive histograms with known
	and configurable error.

	The datastructure uses roughly O(log(log(N))) memory and O(N log(N))
	time.


			// populate data
			s := NewStreamSummaryStructure(0.001)
			for i := 0; i < 5000; i++ {
				s.Update(rand.Float64())
			}

			// extract histogram
			h := s.Histogram()
			for p := 1; p <= 5000; p += 500 {
				println(p, h.Quantile(p))
			}

	The summary structure is not thread safe.
*/
package streamhist
