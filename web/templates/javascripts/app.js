function shadeColor(color, porcent) {

    var R = parseInt(color.substring(1,3),16)
    var G = parseInt(color.substring(3,5),16)
    var B = parseInt(color.substring(5,7),16);

    R = parseInt(R * (100 + porcent) / 100);
    G = parseInt(G * (100 + porcent) / 100);
    B = parseInt(B * (100 + porcent) / 100);

    R = (R<255)?R:255;  
    G = (G<255)?G:255;  
    B = (B<255)?B:255;  

    var RR = ((R.toString(16).length==1)?"0"+R.toString(16):R.toString(16));
    var GG = ((G.toString(16).length==1)?"0"+G.toString(16):G.toString(16));
    var BB = ((B.toString(16).length==1)?"0"+B.toString(16):B.toString(16));

    return "#"+RR+GG+BB;
}

var oneHour = 60*60*1000
var settings = {
	duration: 4*oneHour,  // duration of the plot
	resolution: 75,     // number of data points per plot
	layoutFactor: 1     // increase the resolution based on layout
}
var charts = []

var baseOptions = {
	colors: [shadeColor('#5B4F9C', 20)],
	showRoller: false,
	strokeWidth: 2,
	pointSize: 2,
	drawPointCallback: Dygraph.Circles.SQUARE,
	drawPoints: true,
	// stepPlot: true, 
	axisLineColor: shadeColor('#f2efee', -20),
	gridLineColor: shadeColor('#f2efee', -10),
	includeZero: true,
	yAxisLabelWidth: 30,
	errorBars: true,
	labelsKMB: true,
	sigma: 0.5
}

function chart(options) {
	start = new Date() - settings.duration
	stop = new Date() - 0
	resolution = settings.resolution * settings.layoutFactor

	if (options['rate'] == undefined)
		options['rate'] = true

	$.ajax({
		dataType: 'json',
		url: '/api/time',
		data: $.extend({}, options, {
			start: start,
			stop: stop,
			resolution: resolution
		}),
		success: function (data, status, xhr) {
			var pts = data['points']
			var processed = []
			var max = 0
			for (var i = pts.length - 1; i >= 0; i--) {
				pts[i][0] = new Date(pts[i][0])
				if (pts[i][1] == null) {
					pts[i][1] = 0
				} else {
					pts[i][1] = [pts[i][1], pts[i][2]]
					if (pts[i][1][0] + pts[i][1][1] > max)
						max = pts[i][1][0] + pts[i][1][1]
				}
				processed[processed.length] = [pts[i][0], pts[i][1]] 
			}

			// don't want a plot range of [0, 0]
			if (max == 0)
				max = 1

			if (charts[options['query']] == undefined) {
				// initialize chart object
				charts[options['query']] = {
					options: options,
					graph: null
				}

				charts[options['query']].div = document.getElementById('plot-'+options['query'])
				charts[options['query']].graph = new Dygraph(
					charts[options['query']].div,
					processed,
					$.extend({}, baseOptions, {
						valueRange: [0, 1.2*max],
						dateWindow: [start, stop]
					})
				)
			} else {
				charts[options['query']].graph.updateOptions({
					file: processed,
					valueRange: [0, 1.2*max],
					dateWindow: [start, stop]
				})
			}

			$(charts[options['query']].div).removeClass('drawing')
		}
	})
}

function setDuration(hours) {
	settings.duration = hours * oneHour
	redraw()
}

function setResolution(points) {
	settings.resolution = points
	redraw()
}

function redraw() {
	// freeze all charts
	for (var name in charts) {
		$(charts[name].div).addClass('drawing')
	}
	for (var name in charts) {
		chart(charts[name].options)
	}
}