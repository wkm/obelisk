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
	overlay: 1,         // overlay factor
	duration: 4*oneHour,  // duration of the plot
	resolution: 75,     // number of data points per plot
	layoutFactor: 1     // increase the resolution based on layout
}
var charts = []

var baseOptions = {
	colors: [shadeColor('#5B4F9C', 20), shadeColor('#f2efee', -20)],
	showRoller: false,
	strokeWidth: 2,
	pointSize: 2,
	// stepPlot: true, 
	axisLineColor: shadeColor('#f2efee', -20),
	gridLineColor: shadeColor('#f2efee', -10),
	includeZero: true,
	yAxisLabelWidth: 30,
	errorBars: true,
	labelsKMB: true,
	sigma: 0.5,
	logscale: false,
	labels: ['date', 'a', 'b'],
	'a': {
		drawPoints: true
	},
	'b': {
		drawPoints: false,
		strokePattern: Dygraph.DOTTED_LINE,
		strokeWidth: 1
	}
}

function chart(options) {
	overlay = settings.overlay
	start = new Date() - settings.duration * overlay
	stop = new Date() - 0
	resolution = settings.resolution * settings.layoutFactor * overlay

	if (options['rate'] == undefined)
		options['rate'] = true

	$.ajax({
		dataType: 'json',
		url: '/api/time',
		data: $.extend({}, options, {
			start: start,
			stop: stop,
			resolution: resolution,
			overlay: overlay // ignored for now (but logged)
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

			// build overlay factor
			var overlayed = []
			for (var i = 0; i < resolution/overlay; i++) {
				// get the original date shown
				overlayed[i] = [processed[overlay * i][0]]

				// fill out data points for this date across all overlays, but the first
				// entry should be the most recent
				for (var j = 0; j < overlay; j++) {
					overlayed[i][overlay - j] = processed[j*resolution/overlay + i][1]
				}
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
			}

			if (charts[options['query']].graph == null) {
				charts[options['query']].div = document.getElementById('plot-'+options['query'])
				charts[options['query']].graph = new Dygraph(
					charts[options['query']].div,
					overlayed,
					$.extend({}, baseOptions, {
						valueRange: [0, 1.2*max],
						dateWindow: [start, stop]
					})
				)
			} else {
				charts[options['query']].graph.updateOptions({
					file: overlayed,
					valueRange: [0, 1.2*max],
					dateWindow: [start, stop]
				})
			}

			$(charts[options['query']].div).removeClass('drawing')
		}
	})
}

function setOverlay(overlay) {
	settings.overlay = overlay
	redraw(true)
}

function setDuration(hours) {
	settings.duration = hours * oneHour
	redraw(false)
}

function setResolution(points) {
	settings.resolution = points
	redraw(false)
}

function redraw(destroy) {
	// freeze all charts
	for (var name in charts) {
		$(charts[name].div).addClass('drawing')
		if (destroy) {
			charts[name].graph.destroy()
			charts[name].graph = null
		}
	}
	for (var name in charts) {
		chart(charts[name].options)
	}
}