var loadedWebsiteData;
var editName, editProtocol, editUrl, editMethod;
var allowCheck = true;

$(document).ready(function() {
	$('#input-add-name').keypress(function(event) {
		if (event.keyCode == 13) {
			addWebsite();
		}
	});

	$('#input-add-url').keypress(function(event) {
		if (event.keyCode == 13) {
			addWebsite();
		}
	});

	$('#input-edit-name').keypress(function(event) {
		if (event.keyCode == 13) {
			saveWebsite();
		}
	});

	$('#input-edit-url').keypress(function(event) {
		if (event.keyCode == 13) {
			saveWebsite();
		}
	});

	$('#input-new-password').keypress(function(event) {
		if (event.keyCode == 13) {
			changePassword();
		}
	});

	$('#input-new-interval').keypress(function(event) {
		if (event.keyCode == 13) {
			changeInterval();
		}
	});

	loadWebsites();

	setInterval(loadWebsites, 60 * 1000);
});

function loadWebsites() {
	$.ajax({
		url: "/api/v2/websites",
		type: "GET",
		success: function(data) {
			loadedWebsiteData = data.websites;
			var dataString = '';
			for (var i = 0; i < loadedWebsiteData.length; i++) {
				dataString += '<tr><td>' + (i + 1) + '</td><td>' + loadedWebsiteData[i].name + '</td><td>';

				if (loadedWebsiteData[i].enabled) {
					dataString += ' <span class="label label-success label-action" onclick="disableWebsite(\'' + loadedWebsiteData[i].url + '\')">Enabled</span> </td><td>';
				} else {
					dataString += ' <span class="label label-warning label-action" onclick="enableWebsite(\'' + loadedWebsiteData[i].url + '\')">Disabled</span> </td><td>';
				}

				if (loadedWebsiteData[i].visible) {
					dataString += ' <span class="label label-success label-action" onclick="invisibleWebsite(\'' + loadedWebsiteData[i].url + '\')">Visbile</span> ';
				} else {
					dataString += ' <span class="label label-warning label-action" onclick="visibleWebsite(\'' + loadedWebsiteData[i].url + '\')">Invisible</span> ';
				}

				dataString += '</td><td>' + loadedWebsiteData[i].protocol + '</td><td>' + loadedWebsiteData[i].url + '</td><td><code>' + loadedWebsiteData[i].checkMethod + '</code></td><td>';

				if (loadedWebsiteData[i].status.indexOf("2") == 0) {
					dataString += ' <span class="label label-success">' + loadedWebsiteData[i].status + '</span> ';
				} else if (loadedWebsiteData[i].status.indexOf("3") == 0) {
					dataString += ' <span class="label label-warning">' + loadedWebsiteData[i].status + '</span> ';
				} else {
					dataString += ' <span class="label label-danger">' + loadedWebsiteData[i].status + '</span> ';
				}

				if (loadedWebsiteData[i].time === '0000-00-00 00:00:00') {
					dataString += '</td><td>never</td>';
				} else {
					var date = new Date(loadedWebsiteData[i].time.replace(' ', 'T'));
					dataString += '</td><td>' + date.toLocaleDateString() + ' ' + date.toLocaleTimeString() + '</td>';
				}

				if (loadedWebsiteData[i].notifications.pushbullet) {
					dataString += '<td><span class="label label-info label-action" onclick="editNotificationPushbullet(\'' + loadedWebsiteData[i].url + '\')" title="Pushbullet"><span class="fa fa-bell"></span></span> ';
				} else {
					dataString += '<td><span class="label label-info-inactive label-action" onclick="editNotificationPushbullet(\'' + loadedWebsiteData[i].url + '\')" title="Pushbullet"><span class="fa fa-bell"></span></span> ';
				}
				if (loadedWebsiteData[i].notifications.email) {
					dataString += '<span class="label label-info label-info label-action" onclick="editNotificationEmail(\'' + loadedWebsiteData[i].url + '\')" title="Email"><span class="fa fa-envelope"></span></span> ';
				} else {
					dataString += '<span class="label label-info-inactive label-info label-action" onclick="editNotificationEmail(\'' + loadedWebsiteData[i].url + '\')" title="Email"><span class="fa fa-envelope"></span></span> ';
				}
				if (loadedWebsiteData[i].notifications.telegram) {
					dataString += '<span class="label label-info label-info label-action" onclick="editNotificationTelegram(\'' + loadedWebsiteData[i].url + '\')" title="Telegram"><span class="fa fa-paper-plane"></span></span></td>';
				} else {
					dataString += '<span class="label label-info-inactive label-info label-action" onclick="editNotificationTelegram(\'' + loadedWebsiteData[i].url + '\')" title="Telegram"><span class="fa fa-paper-plane"></span></span></td>';
				}

				dataString += '<td><span class="label label-default label-action" onclick="showWebsiteDetails(\'' + loadedWebsiteData[i].url + '\')" title="More"><span class="fa fa-info"></span></span> ' +
					'<span class="label label-default label-action" onclick="showWebsiteResponseTimes(\'' + loadedWebsiteData[i].url + '\')" title="Response Times"><span class="fa fa-line-chart"></span></span> ' +
					'<span class="label label-info label-action" onclick="checkWebsite(\'' + loadedWebsiteData[i].url + '\')" title="Check Now"><span class="fa fa-repeat"></span></span> ' +
					'<span class="label label-primary label-action" onclick="editWebsite(\'' + loadedWebsiteData[i].url + '\')" title="Edit"><span class="fa fa-pencil"></span></span> ' +
					'<span class="label label-danger label-action" onclick="deleteWebsite(\'' + loadedWebsiteData[i].url + '\')" title="Delete"><span class="fa fa-trash"></span></span></td></tr>';
			}
			if (dataString === '') {
				dataString = '<tr><td colspan="11">No Websites found.</td></tr>';
			}
			$('#table-websites').html(dataString);
		},
		error: function(error) {
			$('#table-websites').html('<tr><td colspan="11">An error occurred. Please authenticate again or add a website.</td></tr>');
		}
	});
}

function reloadWebsites() {
	loadWebsites();
}

function showWebsiteDetails(website) {
	if (website == "") {
		return;
	}

	$.ajax({
		url: "/api/v2/websites/" + website + "/status",
		type: "GET",
		success: function(data) {
			delete data['requestSuccess'];
			delete data['websiteData'];
			swal({
				title: website,
				html: '<pre>' + JSON.stringify(data, null, '\t') + '</pre>',
				confirmButtonText: "Close"
			});
		},
		error: handleAjaxErrorAlert
	});
}

function showWebsiteResponseTimes(website) {
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

			swal({
				title: website,
				html: '<canvas id="graph-responsetime" height="200" width="800"></canvas>',
				width: 900,
				confirmButtonText: "Close"
			});

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
		},
		error: handleAjaxErrorAlert
	});
}

function addWebsite() {
	var name = $('#input-add-name').val();
	var protocol = $('#input-add-protocol').val();
	var url = $('#input-add-url').val();
	var method = $('#input-add-method').val();

	if (name.trim() && protocol.trim() && url.trim() && method.trim()) {
		$.ajax({
			url: "/api/v2/websites/" + url,
			type: "POST",
			data: {name: name, protocol: protocol, url: url, checkMethod: method},
			success: function() {
				$('#input-add-name').val('');
				$('#input-add-protocol').val('https');
				$('#input-add-url').val('');
				$('#input-add-method').val('HEAD');
				loadWebsites();

				showSuccessAlert("Website successfully added.");
			},
			error: handleAjaxErrorAlert
		});
	} else {
		showErrorAlert("Please fill in all fields to continue.");
	}
}

function enableWebsite(url) {
	$.ajax({
		url: "/api/v2/websites/" + url + "/enabled",
		type: "PUT",
		data: {enabled: true},
		success: function() {
			loadWebsites();
		},
		error: handleAjaxErrorAlert
	});
}

function disableWebsite(url) {
	$.ajax({
		url: "/api/v2/websites/" + url + "/enabled",
		type: "PUT",
		data: {enabled: false},
		success: function() {
			loadWebsites();
		},
		error: handleAjaxErrorAlert
	});
}

function visibleWebsite(url) {
	$.ajax({
		url: "/api/v2/websites/" + url + "/visibility",
		type: "PUT",
		data: {visible: true},
		success: function() {
			loadWebsites();
		},
		error: handleAjaxErrorAlert
	});
}

function invisibleWebsite(url) {
	$.ajax({
		url: "/api/v2/websites/" + url + "/visibility",
		type: "PUT",
		data: {visible: false},
		success: function() {
			loadWebsites();
		},
		error: handleAjaxErrorAlert
	});
}

function editNotificationPushbullet(url) {
	if (!url.trim()) return;

	$.ajax({
		url: "/api/v2/websites/" + url + "/notifications",
		type: "GET",
		success: function(data) {
			swal({
				title: "Pushbullet",
				html: "Please enter a valid <b>Pushbullet-API Key</b> in order to receive push-messages.<br />Leave this field blank if you do not want this kind of notification.<br /><br /><input class='form-control' type='text' id='input-pushbullet' placeholder='API key' value=" + data.notifications.pushbulletKey + ">",
				showCancelButton: true,
				confirmButtonText: "Save"
			}).then(
				function() {
					var inputValue = $('#input-pushbullet').val();
					if (inputValue === false) return;

					$.ajax({
						url: "/api/v2/websites/" + url + "/notifications",
						type: "PUT",
						data: {
							pushbulletKey: inputValue.trim(),
							email: data.notifications.email,
							telegramId: data.notifications.telegramId
						},
						success: function() {
							loadWebsites();
							showSuccessAlert("Settings have been updated.");
						},
						error: handleAjaxErrorAlert
					});
				},
				function(dismiss) {
				}
			);
		},
		error: handleAjaxErrorAlert
	});
}

function editNotificationEmail(url) {
	if (!url.trim()) return;

	$.ajax({
		url: "/api/v2/websites/" + url + "/notifications",
		type: "GET",
		success: function(data) {
			swal({
				title: "Email",
				html: "Please enter a valid <b>email address</b> in order to receive email-notifications.<br />Leave this field blank if you do not want this kind of notification.<br /><br /><input class='form-control' type='text' id='input-email' placeholder='email address' value=" + data.notifications.email + ">",
				showCancelButton: true,
				confirmButtonText: "Save"
			}).then(
				function() {
					var inputValue = $('#input-email').val();
					if (inputValue === false) return;

					$.ajax({
						url: "/api/v2/websites/" + url + "/notifications",
						type: "PUT",
						data: {
							pushbulletKey: data.notifications.pushbulletKey,
							email: inputValue.trim(),
							telegramId: data.notifications.telegramId
						},
						success: function() {
							loadWebsites();
							showSuccessAlert("Settings have been updated.");
						},
						error: handleAjaxErrorAlert
					});
				},
				function(dismiss) {
				}
			);
		},
		error: handleAjaxErrorAlert
	});
}

function editNotificationTelegram(url) {
	if (!url.trim()) return;

	$.ajax({
		url: "/api/v2/websites/" + url + "/notifications",
		type: "GET",
		success: function(data) {
			swal({
				title: "Telegram",
				html: "Please enter a valid <b>Telegram user-id</b> in order to receive Telegram-messages.<br />Leave this field blank if you do not want this kind of notification.<br /><br /><input class='form-control' type='text' id='input-telegram' placeholder='Telegram user id' value=" + data.notifications.telegramId + ">",
				showCancelButton: true,
				confirmButtonText: "Save"
			}).then(
				function() {
					var inputValue = $('#input-telegram').val();
					if (inputValue === false) return;

					$.ajax({
						url: "/api/v2/websites/" + url + "/notifications",
						type: "PUT",
						data: {
							pushbulletKey: data.notifications.pushbulletKey,
							email: data.notifications.email,
							telegramId: inputValue.trim()
						},
						success: function() {
							loadWebsites();
							showSuccessAlert("Settings have been updated.");
						},
						error: handleAjaxErrorAlert
					});
				},
				function(dismiss) {
				}
			);
		},
		error: handleAjaxErrorAlert
	});
}

function editWebsite(url) {
	editUrl = url;
	$('#row-edit-website').fadeIn(200);

	for (var i = 0; i < loadedWebsiteData.length; i++) {
		if (url === loadedWebsiteData[i].url) {
			editName = loadedWebsiteData[i].name;
			editProtocol = loadedWebsiteData[i].protocol;
			editMethod = loadedWebsiteData[i].checkMethod;

			$('#input-edit-name').val(loadedWebsiteData[i].name);
			$('#input-edit-protocol').val(loadedWebsiteData[i].protocol);
			$('#input-edit-url').val(loadedWebsiteData[i].url);
			$('#input-edit-method').val(loadedWebsiteData[i].checkMethod);
		}
	}
}

function saveWebsite() {
	var name = $('#input-edit-name').val();
	var protocol = $('#input-edit-protocol').val();
	var url = $('#input-edit-url').val();
	var method = $('#input-edit-method').val();

	if (name == editName && protocol == editProtocol && editUrl == url && editMethod == method) {
		cancelSaveWebsite();
		return;
	}

	if (name.trim() && protocol.trim() && url.trim() && method.trim()) {
		$.ajax({
			url: "/api/v2/websites/" + editUrl,
			type: "PUT",
			data: {name: name, protocol: protocol, url: url, checkMethod: method},
			success: function() {
				cancelSaveWebsite();
				loadWebsites();

				showSuccessAlert("Website successfully edited.")
			},
			error: handleAjaxErrorAlert
		});
	} else {
		showErrorAlert("Please fill in all fields to continue.");
	}
}

function cancelSaveWebsite() {
	$('#row-edit-website').fadeOut(200);
}

function deleteWebsite(url) {
	if (!url.trim()) return;
	swal({
		title: "Are you sure?",
		text: "The website's settings and check-results will be lost forever. You can not undo this operation.",
		type: "warning",
		showCancelButton: true,
		confirmButtonColor: "#DD6B55",
		confirmButtonText: "Yes"
	}).then(
		function() {
			$.ajax({
				url: "/api/v2/websites/" + url,
				type: "DELETE",
				success: function() {
					loadWebsites();
					showSuccessAlert("Website successfully deleted.");
				},
				error: handleAjaxErrorAlert
			});
		},
		function(dismiss) {
		}
	);
}

function checkWebsite(url) {
	if (!url.trim()) return;
	$.ajax({
		url: "/api/v2/websites/" + url + "/check",
		type: "GET",
		success: function() {
			loadWebsites();
		},
		error: function(error) {
			handleAjaxErrorAlert(error);
			allowCheck = true;
		}
	});
}

function changePassword() {
	var newPassword = $('#input-new-password').val();

	if (newPassword.trim()) {
		$.ajax({
			url: "/api/v2/settings/password",
			type: "PUT",
			data: {password: newPassword},
			success: function() {
				$('#input-new-password').val('');

				showSuccessAlert("Settings have been updated.");
			},
			error: handleAjaxErrorAlert
		});
	} else {
		showErrorAlert("Please enter a valid password to continue.");
	}
}

function changeInterval() {
	var newInterval = $('#input-new-interval').val();

	if (newInterval.trim() && !(isNaN(newInterval) || newInterval < 1 || newInterval > 600)) {
		$.ajax({
			url: "/api/v2/settings/interval",
			type: "PUT",
			data: {interval: newInterval},
			success: function() {
				showSuccessAlert("Settings have been updated.");
			},
			error: handleAjaxErrorAlert
		});
	} else {
		showErrorAlert("Please enter a valid interval (between 1 and 600 seconds) to continue.");
	}
}

function logout() {
	$.ajax({
		url: "/api/v2/auth/logout",
		type: "GET",
		success: function() {
			window.location.replace("/");
		},
		error: handleAjaxErrorAlert
	});
}