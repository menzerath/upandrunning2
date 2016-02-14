$(document).ready(function() {
	var pwField = $('#input-password');
	pwField.keypress(function(event) {
		if (event.keyCode == 13) {
			login();
		}
	});
	pwField.focus();
});

function login() {
	var password = $('#input-password').val();

	if (password.trim()) {
		$.ajax({
			url: "/api/v1/auth/login",
			type: "POST",
			data: {"password": password},
			success: function() {
				showSuccessAlert("Logging you in...");
				setTimeout(function() {
					window.location.replace("/admin");
				}, 1000);
			},
			error: handleAjaxErrorAlert
		});
	} else {
		showErrorAlert("Please enter a Password to continue.");
	}
}