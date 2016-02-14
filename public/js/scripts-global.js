function showSuccessAlert(text) {
	swal("Success!", text, "success");
}

function showErrorAlert(text) {
	swal("Oops!", text, "error");
}

function handleAjaxErrorAlert(error) {
	if (error.status === 0) {
		swal("Oops!", "Could not connect to API.", "error");
	} else {
		swal("Oops!", JSON.parse(error.responseText).message, "error");
	}
}