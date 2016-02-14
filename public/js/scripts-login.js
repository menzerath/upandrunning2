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
				swal("Success!", "Logging you in...", "success");
				setTimeout(function() {
					window.location.replace("/admin");
				}, 1000);
			},
			error: function(error) {
				swal("Oops!", JSON.parse(error.responseText).message, "error");
			}
		});
	}
}