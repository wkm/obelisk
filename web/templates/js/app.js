function chart(metric) {
	var start = new Date() - (60*60*24*1*1000) // a day ago
	var stop = new Date() - 0

	$.ajax({
		dataType: 'json',
		url: '/api/time',
		data: {
			query: metric,
			start: start,
			stop:  stop
		},
		success: function (data, status, xhr ) {
			console.log('plotting ...')
			var pts = data['points']
			for (var i = pts.length - 1; i >= 0; i--) {
				pts[i][0] = new Date(pts[i][0])
			}

			var g = new Dygraph(
				document.getElementById('plot-'+metric),
				pts,
				{
					colors: ['#204a87'],
					errorBars: false,
					showRoller: false,
					strokeWidth: 2.5,
					stepPlot: true,
					axisLineColor: '#babdb6',
					gridLineColor: '#d3d7cf',
					includeZero: true,
					dateWindow: [start,stop],
					yAxisLabelWidth: 30
				}
			)
		}
	})
}