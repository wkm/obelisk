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

var charts = []

function chart(options) {
	var oneDay = 60*60*24*1000

	if (options['start'] == undefined)
		options['start'] = new Date() - (oneDay/6)
	if (options['stop'] == undefined)
		options['stop'] = new Date() - 0
	if (options['resolution'] == undefined)
		options['resolution'] = 75
	if (options['rate'] == undefined)
		options['rate'] = true

	$.ajax({
		dataType: 'json',
		url: '/api/time',
		data: options,
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
			if (max == 0)
				max = 1

			chart[options['query']] = new Dygraph(
				document.getElementById('plot-'+options['query']),
				processed,
				{
					colors: [shadeColor('#5B4F9C', 20)],
					valueRange: [0, 1.2*max],
					showRoller: false,
					strokeWidth: 2,
					pointSize: 2,
					drawPointCallback: Dygraph.Circles.SQUARE,
					drawPoints: true,
					// stepPlot: true, 
					axisLineColor: shadeColor('#f2efee', -20),
					gridLineColor: shadeColor('#f2efee', -10),
					includeZero: true,
					dateWindow: [options.start, options.stop],
					yAxisLabelWidth: 30,
					errorBars: true,
					labelsKMB: true,
					sigma: 0.5
				}
			)
		}
	})
}