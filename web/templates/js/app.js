function chart(metric) {
	var oneDay = 60*60*24*1000
	var start = new Date() - (oneDay/2) 
	var stop = new Date() - 0

	$.ajax({
		dataType: 'json',
		url: '/api/time',
		data: {
			query: metric,
			start: start,
			stop:  stop,
			resolution: 75
		},
		success: function (data, status, xhr ) {
			var pts = data['points']
			var max = 0
			for (var i = pts.length - 1; i >= 0; i--) {
				pts[i][0] = new Date(pts[i][0])
				pts[i][1] = [pts[i][1], pts[i][2]]

				if (pts[i][1][0] + pts[i][1][1] > max) {
					max = pts[i][1][0] + pts[i][1][1]
				}
			}

			var g = new Dygraph(
				document.getElementById('plot-'+metric),
				pts,
				{
					colors: ['#3465a4'],
					valueRange: [0, 1.2*max],
					showRoller: false,
					strokeWidth: 1.5,
					pointSize: 1.5,
					drawPoints: true,
					// stepPlot: true,
					axisLineColor: '#d3d7cf',
					gridLineColor: 'rgb(236,236,234)',
					includeZero: true,
					dateWindow: [start,stop],
					yAxisLabelWidth: 30,
					errorBars: true,
					labelsKMB: true,
					sigma: 1.0
				}
			)
		}
	})
}