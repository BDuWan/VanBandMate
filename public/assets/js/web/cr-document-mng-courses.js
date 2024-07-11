$(document).ready(function () {
    $("#ip-name").change(function () {
        // Get the selected file
        var file = this.files[0];

        if (file) {
            // Create a FileReader to read the file
            var reader = new FileReader();

            // Set up the FileReader onload event
            reader.onload = function (e) {
                // Update the src attribute of the img element
                $("#img-file").attr("src", e.target.result);
            };

            // Read the file as a data URL (base64 encoded)
            reader.readAsDataURL(file);
        }
    });
});