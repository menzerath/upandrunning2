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

	$('#input-new-redirects').keypress(function(event) {
		if (event.keyCode == 13) {
			changeRedirects();
		}
	});

	loadWebsites();

	setInterval(loadWebsites, 60 * 1000);
});

function loadWebsites() {
	$.ajax({
		url: "/api/v1/websites",
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

				dataString += '<td><span class="label label-info label-action" onclick="editNotificationPushbullet(\'' + loadedWebsiteData[i].url + '\')" title="Pushbullet"><span class="fa fa-bell"></span></span> ' +
					'<span class="label label-info label-info label-action" onclick="editNotificationEmail(\'' + loadedWebsiteData[i].url + '\')" title="Email"><span class="fa fa-envelope"></span></span></td>';

				dataString += '<td><span class="label label-default label-action" onclick="showWebsiteDetails(\'' + loadedWebsiteData[i].url + '\')" title="More"><span class="fa fa-info"></span></span> ' +
					'<span class="label label-primary label-action" onclick="editWebsite(\'' + loadedWebsiteData[i].url + '\')" title="Edit"><span class="fa fa-pencil"></span></span> ' +
					'<span class="label label-danger label-action" onclick="deleteWebsite(\'' + loadedWebsiteData[i].url + '\')" title="Delete"><span class="fa fa-trash"></span></span></td></tr>';
			}
			$('#table-websites').html(dataString);
		},
		error: function(error) {
			$('#table-websites').html('<tr><td colspan="11">An error occurred. Please authenticate again or add a website.</td></tr>');
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

function showWebsiteDetails(website) {
	if (website == "") {
		return;
	}

	$.ajax({
		url: "/api/v1/websites/" + website + "/status",
		type: "GET",
		success: function(data) {
			delete data['requestSuccess'];
			delete data['websiteData'];
			swal({
				title: website,
				text: '<pre>' + JSON.stringify(data, null, '\t') + '</pre>',
				html: true,
				confirmButtonText: "Close"
			});
		},
		error: function(error) {
			$('.bottom-right').notify({
				type: 'danger',
				message: {text: "Sorry, but I was unable to process your Request. Error: " + JSON.parse(error.responseText).message},
				fadeOut: {enabled: true, delay: 3000}
			}).show();
		}
	});
}

function addWebsite() {
	var name = $('#input-add-name').val();
	var protocol = $('#input-add-protocol').val();
	var url = $('#input-add-url').val();
	var method = $('#input-add-method').val();

	if (name.trim() && protocol.trim() && url.trim() && method.trim()) {
		$.ajax({
			url: "/api/v1/websites/" + url,
			type: "POST",
			data: {name: name, protocol: protocol, url: url, checkMethod: method},
			success: function() {
				$('#input-add-name').val('');
				$('#input-add-protocol').val('https');
				$('#input-add-method').val('HEAD');
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
		url: "/api/v1/websites/" + url + "/enabled",
		type: "PUT",
		data: {enabled: true},
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
		url: "/api/v1/websites/" + url + "/enabled",
		type: "PUT",
		data: {enabled: false},
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
		url: "/api/v1/websites/" + url + "/visibility",
		type: "PUT",
		data: {visible: true},
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
		url: "/api/v1/websites/" + url + "/visibility",
		type: "PUT",
		data: {visible: false},
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

function editNotificationPushbullet(url) {
	if (!url.trim()) return;

	$.ajax({
		url: "/api/v1/websites/" + url + "/notifications",
		type: "GET",
		success: function(data) {
			swal({
				title: "Pushbullet",
				text: "Please enter a valid <b>Pushbullet-API Key</b> in order to recieve push-messages.<br />Leave this field blank if you do not want this kind of notification.",
				html: true,
				type: "input",
				inputPlaceholder: "API key",
				inputValue: data.notifications.pushbulletKey,
				showCancelButton: true,
				confirmButtonText: "Save",
				closeOnConfirm: false
			}, function(inputValue) {
				if (inputValue === false) return;

				$.ajax({
					url: "/api/v1/websites/" + url + "/notifications",
					type: "PUT",
					data: {pushbulletKey: inputValue.trim(), email: data.notifications.email},
					success: function() {
						swal({
							title: "Done!",
							text: "Your settings have been saved.",
							timer: 2000,
							type: "success"
						});
					},
					error: function(error) {
						swal({
							title: "Oops!",
							text: JSON.parse(error.responseText).message,
							timer: 2000,
							type: "error"
						});
					}
				});
			});
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

function editNotificationEmail(url) {
	if (!url.trim()) return;

	$.ajax({
		url: "/api/v1/websites/" + url + "/notifications",
		type: "GET",
		success: function(data) {
			swal({
				title: "Email",
				text: "Please enter a valid <b>email address</b> in order to recieve email-notifications.<br />Leave this field blank if you do not want this kind of notification.",
				html: true,
				type: "input",
				inputPlaceholder: "email address",
				inputValue: data.notifications.email,
				showCancelButton: true,
				confirmButtonText: "Save",
				closeOnConfirm: false
			}, function(inputValue) {
				if (inputValue === false) return;

				$.ajax({
					url: "/api/v1/websites/" + url + "/notifications",
					type: "PUT",
					data: {pushbulletKey: data.notifications.pushbulletKey, email: inputValue.trim()},
					success: function() {
						swal({
							title: "Done!",
							text: "Your settings have been saved.",
							timer: 2000,
							type: "success"
						});
					},
					error: function(error) {
						swal({
							title: "Oops!",
							text: JSON.parse(error.responseText).message,
							timer: 2000,
							type: "error"
						});
					}
				});
			});
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
			url: "/api/v1/websites/" + editUrl,
			type: "PUT",
			data: {name: name, protocol: protocol, url: url, checkMethod: method},
			success: function() {
				cancelSaveWebsite();
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
		confirmButtonText: "Yes",
		closeOnConfirm: false
	}, function() {
		$.ajax({
			url: "/api/v1/websites/" + url,
			type: "DELETE",
			success: function() {
				loadWebsites();
				swal({
					title: "Deleted!",
					text: "This website has been deleted.",
					timer: 2000,
					type: "success"
				});
			},
			error: function(error) {
				swal("Oops!", JSON.parse(error.responseText).message, "error");
			}
		});
	});
}

function changeTitle() {
	var newTitle = $('#input-new-title').val();

	if (newTitle.trim()) {
		$.ajax({
			url: "/api/v1/settings/title",
			type: "PUT",
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
			url: "/api/v1/settings/password",
			type: "PUT",
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
			url: "/api/v1/settings/interval",
			type: "PUT",
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

function changeRedirects() {
	var newRedirects = $('#input-new-redirects').val();

	if (newRedirects.trim() && !(isNaN(newRedirects) || newRedirects < 0 || newRedirects > 10)) {
		$.ajax({
			url: "/api/v1/settings/redirects",
			type: "PUT",
			data: {redirects: newRedirects},
			success: function() {
				$('.bottom-right').notify({
					type: 'success',
					message: {text: "Amount of redirects successfully changed."},
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
			message: {text: "Please enter a valid amount of redirects (number between 0 and 10) to change it."},
			fadeOut: {enabled: true, delay: 3000}
		}).show();
	}
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
		url: "/api/v1/action/check",
		type: "GET",
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
		url: "/api/v1/auth/logout",
		type: "GET",
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