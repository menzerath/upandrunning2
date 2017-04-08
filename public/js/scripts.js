$(document).ready(function() {
	if (location.pathname.split("/")[1] == "status") {
		if (location.pathname.split("/")[2] !== undefined && location.pathname.split("/")[2] !== "") {
			showInformation(location.pathname.split("/")[2]);
		} else {
			history.replaceState('data', '', '/');
		}
	} else if (location.pathname.split("/")[1] == "results") {
		if (location.pathname.split("/")[2] !== undefined && location.pathname.split("/")[2] !== "") {
			showResponseTimeGraph(location.pathname.split("/")[2]);
		} else {
			history.replaceState('data', '', '/');
		}
	}

	loadWebsiteData();

	setInterval(loadWebsiteData, 5 * 60 * 1000);
});

function showInformation(website) {
	if (website == "") {
		return;
	}

	$.ajax({
		url: "/api/v2/websites/" + website + "/status",
		type: "GET",
		success: function(data) {
			var dataString = '<div class="well"><legend>More Information</legend>';
			dataString += '<p>The website at <a href="' + data.websiteData.url + '">' + data.websiteData.url + '</a> is called <b>"' + data.websiteData.name + '"</b>, was checked <b>' + data.availability.total + ' times</b> and has an average availability of <b>' + data.availability.average + '</b>.</p>';

			if (data.lastCheckResult.status !== '0 - unknown') {
				var dateRecent = new Date(data.lastCheckResult.time.replace(' ', 'T'));
				dataString += '<p>The most recent check on <b>' + dateRecent.toLocaleDateString() + '</b> at <b>' + dateRecent.toLocaleTimeString() + '</b> got the following response after <b>' + data.lastCheckResult.responseTime.replace(/\B(?=(\d{3})+(?!\d))/g, ".") + '</b>: <b>' + data.lastCheckResult.status + '</b>.</p>';
			}

			if (data.lastFailedCheckResult.status !== '0 - unknown') {
				var dateFail = new Date(data.lastFailedCheckResult.time.replace(' ', 'T'));
				dataString += '<p>The last failed check on <b>' + dateFail.toLocaleDateString() + '</b> at <b>' + dateFail.toLocaleTimeString() + '</b> failed after <b>' + data.lastFailedCheckResult.responseTime.replace(/\B(?=(\d{3})+(?!\d))/g, ".") + '</b> because of this response: <b>' + data.lastFailedCheckResult.status + '</b>.</p>';
			}

			dataString += '<button class="btn btn-primary" onclick="hideInformation()">Close</button></div>';

			$('#col-form-information').html(dataString);

			// show everything to the user
			hideResponseTime();
			$('#bc-feature').css('display', 'inline-block').text('Status');
			$('#bc-site').css('display', 'inline-block').text(website).html('<a href="/status/' + website + '">' + website + '</a>');
			history.replaceState('data', '', '/status/' + website + '/');

			$('#row-information').fadeIn(200);
		},
		error: handleAjaxErrorAlert
	});
}

function hideInformation() {
	$('#row-information').hide();

	$('#bc-feature').css('display', 'none').text('');
	$('#bc-site').css('display', 'none').text('');
	history.replaceState('data', '', '/');
}

function showResponseTimeGraph(website) {
	if (website == "") {
		return;
	}

	if (typeof responseTimeGraph !== 'undefined') {
		responseTimeGraph.destroy();
	}

	$.ajax({
		url: "/api/v2/websites/" + website + "/results?limit=100",
		type: "GET",
		success: function(data) {
			var chartValuesResponseTimes = [];
			var chartValuesDatestamps = [];
			for (var i = data.results.length - 1; i >= 0; i--) {
				chartValuesResponseTimes.push(parseInt(data.results[i].responseTime.substr(0, data.results[i].responseTime.length - 3)));
				chartValuesDatestamps.push(data.results[i].time);
			}

			var chartData = {
				labels: chartValuesDatestamps,
				datasets: [
					{
						label: "Response Time (in ms)",
						fill: true,
						backgroundColor: "rgba(220,220,220,0.3)",
						borderColor: "rgba(220,220,220,1)",
						pointBorderColor: "rgba(220,220,220,1)",
						pointBackgroundColor: "#fff",
						pointBorderWidth: 1,
						pointHoverRadius: 5,
						pointHoverBackgroundColor: "rgba(220,220,220,1)",
						pointHoverBorderColor: "rgba(220,220,220,1)",
						pointHoverBorderWidth: 2,
						data: chartValuesResponseTimes
					}
				]
			};

			hideInformation();
			$('#row-responsetime').fadeIn(200);

			var ctx = document.getElementById("graph-responsetime").getContext("2d");
			window.responseTimeGraph = new Chart(ctx, {
				type: 'line',
				data: chartData,
				options: {
					legend: {
						display: false
					},
					scales: {
						xAxes: [{
							ticks: {
								display: false
							}
						}],
						yAxes: [{
							ticks: {
								beginAtZero: true,
								fontFamily: 'Roboto'
							}
						}]
					},
					responsive: true
				}
			});

			$('#bc-feature').css('display', 'inline-block').text('Response Times');
			$('#bc-site').css('display', 'inline-block').text(website).html('<a href="/results/' + website + '">' + website + '</a>');
			history.replaceState('data', '', '/results/' + website + '/');
		},
		error: handleAjaxErrorAlert
	});
}

function hideResponseTime() {
	$('#row-responsetime').hide();

	$('#bc-feature').css('display', 'none').text('');
	$('#bc-site').css('display', 'none').text('');
	history.replaceState('data', '', '/');
}

function loadWebsiteData() {
	$.ajax({
		url: "/api/v2/websites",
		type: "GET",
		success: function(data) {
			loadedWebsiteData = data.websites;
			var dataStringUp = '', dataStringDown = '', newEntry = '';
			var countUp = 1, countDown = 1;
			for (var i = 0; i < loadedWebsiteData.length; i++) {
				newEntry = '<td><a href="' + loadedWebsiteData[i].protocol + '://' + loadedWebsiteData[i].url + '" target="_blank">' + loadedWebsiteData[i].name + '</a></td><td>';

				if (loadedWebsiteData[i].status.indexOf("2") == 0) {
					newEntry += ' <span class="label label-success">' + loadedWebsiteData[i].status + '</span> ';
				} else if (loadedWebsiteData[i].status.indexOf("3") == 0) {
					newEntry += ' <span class="label label-warning">' + loadedWebsiteData[i].status + '</span> ';
				} else {
					newEntry += ' <span class="label label-danger">' + loadedWebsiteData[i].status + '</span> ';
				}

				newEntry += '</td><td>';

				var responseTime = loadedWebsiteData[i].responseTime.split(' ')[0];
				if (responseTime >= 500) {
					newEntry += ' <span class="label label-danger">' + loadedWebsiteData[i].responseTime + '</span> ';
				} else if (responseTime >= 200) {
					newEntry += ' <span class="label label-warning">' + loadedWebsiteData[i].responseTime + '</span> ';
				} else if (responseTime >= 0) {
					newEntry += ' <span class="label label-success">' + loadedWebsiteData[i].responseTime + '</span> ';
				} else {
					newEntry += ' <span class="label label-info">' + loadedWebsiteData[i].responseTime + '</span> ';
				}

				newEntry += '</td><td> <span class="label label-primary label-action" onclick="showInformation(\'' + loadedWebsiteData[i].url + '\')" title="More"><span class="fa fa-info"></span></span> ' +
					'<span class="label label-primary label-action" onclick="showResponseTimeGraph(\'' + loadedWebsiteData[i].url + '\')" title="Response Times"><span class="fa fa-line-chart"></span></span> </td>';

				if (loadedWebsiteData[i].status.indexOf("2") == 0 || loadedWebsiteData[i].status.indexOf("3") == 0) {
					dataStringUp += '<tr><td>' + countUp + '</td>' + newEntry + '</tr>';
					countUp++;
				} else {
					dataStringDown += '<tr><td>' + countDown + '</td>' + newEntry + '</tr>';
					countDown++;
				}
			}

			if (dataStringUp === '') {
				dataStringUp = '<tr><td colspan="5">No Websites found.</td></tr>';
			}
			if (dataStringDown === '') {
				dataStringDown = '<tr><td colspan="5">No Websites found.</td></tr>';
			}

			$('#table-websites-up').html(dataStringUp);
			$('#table-websites-down').html(dataStringDown);
		},
		error: function(error) {
			$('#table-websites-up').html('<tr><td colspan="5">An Error occurred: ' + JSON.parse(error.responseText).message + '</td></tr>');
			$('#table-websites-down').html('<tr><td colspan="5">An Error occurred: ' + JSON.parse(error.responseText).message + '</td></tr>');
		}
	});
}