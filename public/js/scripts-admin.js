var loadedWebsiteData;
var editUrl;
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

	$('#input-new-title').keypress(function(event) {
		if (event.keyCode == 13) {
			changeTitle();
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
});

function loadWebsites() {
	$.ajax({
		url: "/api/admin/websites",
		type: "GET",
		success: function(data) {
			loadedWebsiteData = data.websites;
			var dataString = '';
			for (var i = 0; i < loadedWebsiteData.length; i++) {
				dataString += '<tr><td>' + (i + 1) + '</td><td>' + loadedWebsiteData[i].name + '</td><td>';

				if (loadedWebsiteData[i].enabled) {
					dataString += ' <span class="label label-success" id="label-action" onclick="disableWebsite(\'' + loadedWebsiteData[i].url + '\')">Enabled</span> </td><td>';
				} else {
					dataString += ' <span class="label label-warning" id="label-action" onclick="enableWebsite(\'' + loadedWebsiteData[i].url + '\')">Disabled</span> </td><td>';
				}

				if (loadedWebsiteData[i].visible) {
					dataString += ' <span class="label label-success" id="label-action" onclick="invisibleWebsite(\'' + loadedWebsiteData[i].url + '\')">Visbile</span> ';
				} else {
					dataString += ' <span class="label label-warning" id="label-action" onclick="visibleWebsite(\'' + loadedWebsiteData[i].url + '\')">Invisible</span> ';
				}

				dataString += '</td><td>' + loadedWebsiteData[i].protocol + '</td><td>' + loadedWebsiteData[i].url + '</td><td>';

				if (loadedWebsiteData[i].status.indexOf("200") > -1) {
					dataString += ' <span class="label label-success">' + loadedWebsiteData[i].status + '</span> ';
				} else if (loadedWebsiteData[i].status.indexOf("301") > -1 || loadedWebsiteData[i].status.indexOf("302") > -1) {
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

				dataString += '<td>' + loadedWebsiteData[i].avgAvail + '</td>';

				dataString += '<td><span class="label label-primary" id="label-action" onclick="editWebsite(\'' + loadedWebsiteData[i].url + '\')">Edit</span> <span class="label label-danger" id="label-action" onclick="deleteWebsite(\'' + loadedWebsiteData[i].url + '\')">Delete</span></td></tr>';
			}
			$('#table-websites').html(dataString);
		},
		error: function(error) {
			$('#table-websites').html('<tr><td colspan="10">An error occured. Please authenticate again or add a website.</td></tr>');
		}
	});
}

function reloadWebsites() {
	$('.bottom-right').notify({
		type: 'success',
		message: {text: "Reloading websites..."},
		fadeOut: {enabled: true, delay: 3000}
	}).show();
	loadWebsites();
}

function addWebsite() {
	var name = $('#input-add-name').val();
	var protocol = $('#input-add-protocol').val();
	var url = $('#input-add-url').val();

	if (name.trim() && protocol.trim() && url.trim()) {
		$.ajax({
			url: "/api/admin/websites/add",
			type: "POST",
			data: {name: name, protocol: protocol, url: url},
			success: function() {
				$('#input-add-name').val('');
				$('#input-add-protocol').val('https');
				$('#input-add-url').val('');
				loadWebsites();

				$('.bottom-right').notify({
					type: 'success',
					message: {text: "Website successfully added."},
					fadeOut: {enabled: true, delay: 3000}
				}).show();
			},
			error: function(error) {
				$('.bottom-right').notify({
					type: 'danger',
					message: {text: JSON.parse(error.responseText).message},
					fadeOut: {enabled: true, delay: 3000}
				}).show();
			}
		});
	} else {
		$('.bottom-right').notify({
			type: 'danger',
			message: {text: "Please fill in all fields to add a new website."},
			fadeOut: {enabled: true, delay: 3000}
		}).show();
	}
}

function enableWebsite(url) {
	$.ajax({
		url: "/api/admin/websites/enable",
		type: "POST",
		data: {url: url},
		success: function() {
			loadWebsites();
		},
		error: function(error) {
			$('.bottom-right').notify({
				type: 'danger',
				message: {text: JSON.parse(error.responseText).message},
				fadeOut: {enabled: true, delay: 3000}
			}).show();
		}
	});
}

function disableWebsite(url) {
	$.ajax({
		url: "/api/admin/websites/disable",
		type: "POST",
		data: {url: url},
		success: function() {
			loadWebsites();
		},
		error: function(error) {
			$('.bottom-right').notify({
				type: 'danger',
				message: {text: JSON.parse(error.responseText).message},
				fadeOut: {enabled: true, delay: 3000}
			}).show();
		}
	});
}

function visibleWebsite(url) {
	$.ajax({
		url: "/api/admin/websites/visible",
		type: "POST",
		data: {url: url},
		success: function() {
			loadWebsites();
		},
		error: function(error) {
			$('.bottom-right').notify({
				type: 'danger',
				message: {text: JSON.parse(error.responseText).message},
				fadeOut: {enabled: true, delay: 3000}
			}).show();
		}
	});
}

function invisibleWebsite(url) {
	$.ajax({
		url: "/api/admin/websites/invisible",
		type: "POST",
		data: {url: url},
		success: function() {
			loadWebsites();
		},
		error: function(error) {
			$('.bottom-right').notify({
				type: 'danger',
				message: {text: JSON.parse(error.responseText).message},
				fadeOut: {enabled: true, delay: 3000}
			}).show();
		}
	});
}

function editWebsite(url) {
	editUrl = url;
	$('#form-edit-website').fadeIn(200);

	for (var i = 0; i < loadedWebsiteData.length; i++) {
		if (url === loadedWebsiteData[i].url) {
			$('#input-edit-name').val(loadedWebsiteData[i].name);
			$('#input-edit-protocol').val(loadedWebsiteData[i].protocol);
			$('#input-edit-url').val(loadedWebsiteData[i].url);
		}
	}
}

function saveWebsite() {
	var name = $('#input-edit-name').val();
	var protocol = $('#input-edit-protocol').val();
	var url = $('#input-edit-url').val();

	if (name.trim() && protocol.trim() && url.trim()) {
		$.ajax({
			url: "/api/admin/websites/edit",
			type: "POST",
			data: {oldUrl: editUrl, name: name, protocol: protocol, url: url},
			success: function() {
				cancleSaveWebsite();
				loadWebsites();

				$('.bottom-right').notify({
					type: 'success',
					message: {text: "Website successfully edited."},
					fadeOut: {enabled: true, delay: 3000}
				}).show();
			},
			error: function(error) {
				$('.bottom-right').notify({
					type: 'danger',
					message: {text: JSON.parse(error.responseText).message},
					fadeOut: {enabled: true, delay: 3000}
				}).show();
			}
		});
	} else {
		$('.bottom-right').notify({
			type: 'danger',
			message: {text: "Please fill in all fields to save this edited website."},
			fadeOut: {enabled: true, delay: 3000}
		}).show();
	}
}

function cancleSaveWebsite() {
	$('#form-edit-website').fadeOut(200);
}

function deleteWebsite(url) {
	if (window.confirm("Are you sure?")) {
		$.ajax({
			url: "/api/admin/websites/delete",
			type: "POST",
			data: {url: url},
			success: function() {
				loadWebsites();

				$('.bottom-right').notify({
					type: 'success',
					message: {text: "Website successfully deleted."},
					fadeOut: {enabled: true, delay: 3000}
				}).show();
			},
			error: function(error) {
				$('.bottom-right').notify({
					type: 'danger',
					message: {text: JSON.parse(error.responseText).message},
					fadeOut: {enabled: true, delay: 3000}
				}).show();
			}
		});
	}
}

function changeTitle() {
	var newTitle = $('#input-new-title').val();

	if (newTitle.trim()) {
		$.ajax({
			url: "/api/admin/settings/title",
			type: "POST",
			data: {title: newTitle},
			success: function() {
				$(document).attr("title", "Administration | " + newTitle);
				$('#navbar-title').text(newTitle);

				$('.bottom-right').notify({
					type: 'success',
					message: {text: "Title successfully changed."},
					fadeOut: {enabled: true, delay: 3000}
				}).show();
			},
			error: function(error) {
				$('.bottom-right').notify({
					type: 'danger',
					message: {text: JSON.parse(error.responseText).message},
					fadeOut: {enabled: true, delay: 3000}
				}).show();
			}
		});
	} else {
		$('.bottom-right').notify({
			type: 'danger',
			message: {text: "Please enter a valid title to change it."},
			fadeOut: {enabled: true, delay: 3000}
		}).show();
	}
}

function changePassword() {
	var newPassword = $('#input-new-password').val();

	if (newPassword.trim()) {
		$.ajax({
			url: "/api/admin/settings/password",
			type: "POST",
			data: {password: newPassword},
			success: function() {
				$('#input-new-password').val('');

				$('.bottom-right').notify({
					type: 'success',
					message: {text: "Password successfully changed."},
					fadeOut: {enabled: true, delay: 3000}
				}).show();
			},
			error: function(error) {
				$('.bottom-right').notify({
					type: 'danger',
					message: {text: JSON.parse(error.responseText).message},
					fadeOut: {enabled: true, delay: 3000}
				}).show();
			}
		});
	} else {
		$('.bottom-right').notify({
			type: 'danger',
			message: {text: "Please enter a valid password to change it."},
			fadeOut: {enabled: true, delay: 3000}
		}).show();
	}
}

function changeInterval() {
	var newInterval = $('#input-new-interval').val();

	if (newInterval.trim() && !(isNaN(newInterval) || newInterval < 1 || newInterval > 600)) {
		$.ajax({
			url: "/api/admin/settings/interval",
			type: "POST",
			data: {interval: newInterval},
			success: function() {
				$('.bottom-right').notify({
					type: 'success',
					message: {text: "Interval successfully changed."},
					fadeOut: {enabled: true, delay: 3000}
				}).show();
			},
			error: function(error) {
				$('.bottom-right').notify({
					type: 'danger',
					message: {text: JSON.parse(error.responseText).message},
					fadeOut: {enabled: true, delay: 3000}
				}).show();
			}
		});
	} else {
		$('.bottom-right').notify({
			type: 'danger',
			message: {text: "Please enter a valid interval (numbers between 1 and 600) to change it."},
			fadeOut: {enabled: true, delay: 3000}
		}).show();
	}
}

function changePBKey() {
	var newKey = $('#input-new-pb_key').val();

	$.ajax({
		url: "/api/admin/settings/pbkey",
		type: "POST",
		data: {key: newKey},
		success: function() {
			$('.bottom-right').notify({
				type: 'success',
				message: {text: "API-Key successfully changed."},
				fadeOut: {enabled: true, delay: 3000}
			}).show();
		},
		error: function(error) {
			$('.bottom-right').notify({
				type: 'danger',
				message: {text: JSON.parse(error.responseText).message},
				fadeOut: {enabled: true, delay: 3000}
			}).show();
		}
	});
}

function checkNow() {
	if (!allowCheck) {
		$('.bottom-right').notify({
			type: 'danger',
			message: {text: "Please wait a few seconds before trying this operation again."},
			fadeOut: {enabled: true, delay: 3000}
		}).show();
		return;
	}

	allowCheck = false;
	$.ajax({
		url: "/api/admin/check",
		type: "POST",
		success: function() {
			$('.bottom-right').notify({
				type: 'success',
				message: {text: "Check triggered. Reload in three seconds."},
				fadeOut: {enabled: true, delay: 3000}
			}).show();
			setTimeout(function() {
				loadWebsites();
			}, 3000);
			setTimeout(function() {
				allowCheck = true;
			}, 10000);
		},
		error: function(error) {
			$('.bottom-right').notify({
				type: 'danger',
				message: {text: JSON.parse(error.responseText).message},
				fadeOut: {enabled: true, delay: 3000}
			}).show();
			allowCheck = true;
		}
	});
}

function logout() {
	$.ajax({
		url: "/api/admin/logout",
		type: "POST",
		success: function() {
			window.location.replace("/");
		},
		error: function(error) {
			$('.bottom-right').notify({
				type: 'danger',
				message: {text: JSON.parse(error.responseText).message},
				fadeOut: {enabled: true, delay: 3000}
			}).show();
		}
	});
}