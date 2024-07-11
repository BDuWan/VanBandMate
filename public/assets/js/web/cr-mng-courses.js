$(document).ready(function () {
    $('#btn-submit').click(function() {
        var name =  $("#ip-name").val();
        var study_program_id= parseInt($("#ip-study_program").val(), 10);
        var user_id= parseInt($("#ip-user").val(), 10);
        var description = $("#ip-description").val();

        var jsonData = {
            name: name,
            study_program_id: study_program_id,
            user_id: user_id,
            description: description,
        };

        $.ajax({
            url: '/managements/mng-courses',
            method: 'POST',
            contentType: "application/json",
            data: JSON.stringify(jsonData),

            success: function (response) {
                if(response === "Success"){
                    swal({
                        title: 'Successfully !',
                        text: 'Create courses successfully',
                        icon: 'success',
                        button: {
                            text: "Close",
                            value: true,
                            visible: true,
                            className: "btn btn-primary"
                        }
                    })
                    $('.ip-u').val("");
                }else{
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
            error: function () {
                console.log('Error occurred while retrieving options');
            }
        });
    });
});