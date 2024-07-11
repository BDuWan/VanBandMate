$(document).ready(function () {
    $('#sb-student').click(function () {

        var jsonData = {
            first_name: $("#first-name-stu").val(),
            last_name: $("#last-name-stu").val(),
            email: $("#email-stu").val(),
            phone_number: $("#phone-number-stu").val(),
            address: $("#address-stu").val(),
            password: $("#password-stu").val(),
            referral_code: $("#referral-code-stu").val(),

        };

        $.ajax({
            url: '/signup/stu',
            method: 'POST',
            contentType: "application/json",
            data: JSON.stringify(jsonData),

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
                    $('.input-stu').val("");
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
            error: function () {
                console.log('Error occurred while retrieving options');
            }
        });
    });

    $('#sb-instructor').click(function () {

        var jsonData = {
            first_name: $("#first-name-instr").val(),
            last_name: $("#last-name-instr").val(),
            email: $("#email-instr").val(),
            phone_number: $("#phone-number-instr").val(),
            address: $("#address-instr").val(),
            password: $("#password-instr").val(),
            referral_code: $("#referral-code-instr").val(),

        };

        $.ajax({
            url: '/signup/instr',
            method: 'POST',
            contentType: "application/json",
            data: JSON.stringify(jsonData),

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
                    $('.input-instr').val("");
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
            error: function () {
                console.log('Error occurred while retrieving options');
            }
        });
    });

    $('#sb-sale').click(function () {

        var jsonData = {
            first_name: $("#first-name-sale").val(),
            last_name: $("#last-name-sale").val(),
            email: $("#email-sale").val(),
            phone_number: $("#phone-number-sale").val(),
            address: $("#address-sale").val(),
            password: $("#password-sale").val(),

        };

        $.ajax({
            url: '/signup/sale',
            method: 'POST',
            contentType: "application/json",
            data: JSON.stringify(jsonData),

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
                    $('.input-sale').val("");
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
            error: function () {
                console.log('Error occurred while retrieving options');
            }
        });
    });

    $('#sb-business').click(function () {

        var jsonData = {
            name_business: $("#name-business").val(),
            full_name_representative: $("#full-name-representative").val(),
            email: $("#email-business").val(),
            phone_number: $("#phone-number-business").val(),
            address: $("#address-business").val(),
            password: $("#password-business").val(),

        };

        $.ajax({
            url: '/signup/business',
            method: 'POST',
            contentType: "application/json",
            data: JSON.stringify(jsonData),

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
                    $('.input-business').val("");
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
            error: function () {
                console.log('Error occurred while retrieving options');
            }
        });
    });
});