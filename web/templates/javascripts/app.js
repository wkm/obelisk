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
					colors: [shadeColor('#4F909C', 30)],
					valueRange: [0, 1.2*max],
					showRoller: false,
					strokeWidth: 2.5,
					pointSize: 2,
					drawPoints: true,
					// stepPlot: true,
					axisLineColor: shadeColor('#A3A0A1', -20),
					gridLineColor: shadeColor('#A3A0A1', -40),
					includeZero: true,
					dateWindow: [start,stop],
					yAxisLabelWidth: 30,
					errorBars: true,
					labelsKMB: true,
					sigma: 0.5
				}
			)
		}
	})
}