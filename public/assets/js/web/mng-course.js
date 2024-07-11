$(document).ready(function () {

    // $('#btn-submit').click(function () {
    //     var image = $("#ip-image").val();
    //     var title = $("#ip-title").val();
    //     var description = $("#ip-description").val();
    //     var price = $("#ip-price").val();
    //
    //
    //     var jsonData = {
    //         title: title,
    //         description: description,
    //         image: image,
    //         price: price,
    //     };
    //
    //     $.ajax({
    //         url: '/managements/mng-courses',
    //         method: 'POST',
    //         contentType: "application/json",
    //         data: JSON.stringify(jsonData),
    //
    //         success: function (response) {
    //             if (response === "Success") {
    //                 swal({
    //                     title: 'Successfully !',
    //                     text: 'Create account successfully',
    //                     icon: 'success',
    //                     button: {
    //                         text: "Close",
    //                         value: true,
    //                         visible: true,
    //                         className: "btn btn-primary"
    //                     }
    //                 })
    //                 $('.ip-u').val("");
    //             } else {
    //                 swal({
    //                     title: 'Error !',
    //                     text: response,
    //                     icon: 'warning',
    //                     button: {
    //                         text: "Close",
    //                         value: true,
    //                         visible: true,
    //                         className: "btn btn-danger"
    //                     }
    //                 })
    //             }
    //
    //         },
    //         error: function () {
    //             console.log('Error occurred while retrieving options');
    //         }
    //     });
    // });

    $("#btn-submit").on("click", function () {
        var formData = new FormData();
        var fileInput = $("#ip-image")[0].files[0];

        // if (!fileInput) {
        //     $("#message").text("Please select a file.");
        //     return;
        // }

        formData.append("title", $("#ip-title").val());
        formData.append("description", $("#ip-description").val());
        formData.append("price", $("#ip-price").val());
        formData.append("image", fileInput);

        $.ajax({
            type: "POST",
            url: "/managements/mng-courses",
            data: formData,
            contentType: false,
            processData: false,
            success: function (response) {
                if (response === "Success") {
                    swal({
                        title: 'Successfully !',
                        text: 'Create account successfully',
                        icon: 'success',
                        button: {
                            text: "Close",
                            value: true,
                            visible: true,
                            className: "btn btn-primary"
                        }
                    })
                    $('.ip-u').val("");
                } else {
                    swal({
                        title: 'Error !',
                        text: response,
                        icon: 'warning',
                        button: {
                            text: "Close",
                            value: true,
                            visible: true,
                            className: "btn btn-danger"
                        }
                    })
                }
            },
            error: function (error) {
                // $("#message").text("Error: " + error.statusText);
            }
        });
    });
})
