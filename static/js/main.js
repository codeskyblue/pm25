$(function() {
	var LineDraw = function(labels, flow) {
		var width = $("#statusDiv").width();
		var Minimum = function(a, b) {
			return a < b ? a: b;
		}
		width = Minimum(width, 800);

		var data = [{
			name: 'PV',
			value: flow,
			color: '#0d8ecf',
			line_width: 2
		}];
		var line = new iChart.LineBasic2D({
			render: 'statusDiv',
			data: data,
			align: 'center',
			title: 'AQI Status',
			width: width,
			height: 300,
			padding: 30,
			sub_option: {
				smooth: true,
				point_size: 10
			},
			tip: {
				enable: true,
				shadow: true
			},
			legend: {
				enable: false
			},
			crosshair: {
				enable: true,
				line_color: '#62bce9'
			},
			coordinate: {
				width: 600,
				valid_width: 500,
				height: 260,
				axis: {
					color: '#9f9f9f',
					width: [0, 0, 2, 2]
				},
				grids: {
					vertical: {
						way: 'share_alike',
						value: 12
					}
				},
				scale: [{
					position: 'left',
					start_scale: 0,
					end_scale: 500,
					scale_space: 100,
					scale_size: 2,
					scale_color: '#9f9f9f'
				},
				{
					position: 'bottom',
					labels: labels
				}]
			}
		});
		line.draw(); 
	}

	var apiUrl = "api/v2/pm25" // /city_list | history
	$.get(apiUrl + "/history", {
		"loc": "beijing"
	},
	function(ret, textStatus) {
		if (ret.error !== null) {
			alert("load history fail");
			return
		}
		var pm25 = [];
		var aqi = [];
		var labels = [];
		for (var i = 0; i < ret.data.length; i++) {
			d = ret.data[i];
			pm25.push(d.Pm25);
			aqi.push(d.Aqi);
			var timePoint = i;//d.timePoint.substr(1, 2);
			labels.push(timePoint);
		}
		var chartData = {
			labels: labels,
			datasets: [{
				fillColor: "rgba(151,187,205,0.5)",
				strokeColor: "rgba(151,187,205,1)",
				data: pm25,
			},
			],
		}

		var $chart = $("#myChart");
		$chart.width(700);
		var ctx = $chart.get(0).getContext("2d");
		new Chart(ctx).Line(chartData);

		LineDraw(labels, aqi);
	},
	"json");

}); // end of jQuery

