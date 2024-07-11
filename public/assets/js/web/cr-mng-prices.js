$(document).ready(function () {
    $('#create').click(function() {
        var price= parseInt($("#price").val(), 10);
        var commission= parseInt($("#commission").val(), 10);
        var description =  $("#description").val();
        var start_time= $("#start-time").val();
        var end_time= $("#end-time").val();

        var jsonData = {
            price: price,
            commission: commission,
            description: description,
            start_time: start_time,
            end_time: end_time,
        };

        $.ajax({
            url: '/managements/mng-payments/prices',
            method: 'POST',
            contentType: "application/json",
            data: JSON.stringify(jsonData),

            success: function (response) {
                if(response === "Success"){
                    swal({
                        title: 'Successfully !',
                        text: 'Create price program successfully',
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
                        title: "Error !",
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