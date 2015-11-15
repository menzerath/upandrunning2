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

	$('.bottom-right').notify({
		type: 'warning',
		message: {text: "Processing..."},
		fadeOut: {enabled: true, delay: 3000}
	}).show();

	if (password.trim()) {
		$.ajax({
			url: "/api/admin/login",
			type: "POST",
			data: {"password": password},
			success: function() {
				window.location.replace("/admin");
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