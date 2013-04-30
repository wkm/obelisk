function chart(metric) {
	var start = new Date() - (60*60*24*1*1000) // a day ago
	var stop = new Date() - 0

	$.ajax({
		dataType: 'json',
		url: '/api/time',
		data: {
			query: metric,
			start: start,
			stop:  stop,
			resolution: 50
		},
		success: function (data, status, xhr ) {
			var pts = data['points']
			for (var i = pts.length - 1; i >= 0; i--) {
				pts[i][0] = new Date(pts[i][0])
				pts[i][1] = [pts[i][1], pts[i][2]]
			}

			var g = new Dygraph(
				document.getElementById('plot-'+metric),
				pts,
				{
					colors: ['#204a87'],
					errorBars: false,
					showRoller: false,
					strokeWidth: 1.5,
					pointSize: 2,
					drawPoints: true,
					// stepPlot: true,
					axisLineColor: '#d3d7cf',
					gridLineColor: 'rgb(236,236,234)',
					includeZero: true,
					dateWindow: [start,stop],
					yAxisLabelWidth: 30,
					errorBars: true,
					sigma: 0.5
				}
			)
		}
	})
}