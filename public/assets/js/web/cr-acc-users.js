$(document).ready(function () {
    $("#form-referral-code").hide();

    $('#ip-type-user').change(function () {
        var idTypeUser = parseInt($(this).val(), 10);
        if (idTypeUser === 3) {
            $("#lb-first-name").text("Name of Business");
            $("#lb-last-name").text("Full name of the representative");
            $("#ip-first-name").attr("placeholder", "Enter name of Business");
            $("#ip-last-name").attr("placeholder", "Enter full name of the representative");
            $("#ip-email").attr("placeholder", "Enter business email");
            $("#ip-phone-number").attr("placeholder", "Enter business phone number");
            $("#ip-address").attr("placeholder", "Enter business address");
            $("#ip-password").attr("placeholder", "Enter business password");
        } else {
            $("#lb-first-name").text("First Name");
            $("#lb-last-name").text("Last Name");
            $("#ip-first-name").attr("placeholder", "Enter your first name");
            $("#ip-last-name").attr("placeholder", "Enter your last name");
            $("#ip-email").attr("placeholder", "Enter your email");
            $("#ip-phone-number").attr("placeholder", "Enter your phone number");
            $("#ip-address").attr("placeholder", "Enter your address");
            $("#ip-password").attr("placeholder", "Enter your password");
        }

        if (idTypeUser !== 5){
            $("#form-referral-code").hide();
        }else{
            $("#form-referral-code").show();
        }
    });

    $('#btn-submit').click(function() {
        var type_user_id= parseInt($("#ip-type-user").val(), 10);
        var role_id= parseInt($("#ip-role").val(), 10);
        var first_name =  $("#ip-first-name").val();
        var last_name= $("#ip-last-name").val();
        var name_business =  $("#ip-first-name").val();
        var full_name_representative= $("#ip-last-name").val();
        var email= $("#ip-email").val();
        var phone_number= $("#ip-phone-number").val();
        var address= $("#ip-address").val();
        var password= $("#ip-password").val();
        var referral_code= $("#ip-referral-code").val();

        if (type_user_id !== 5){
            referral_code = "";
        }
        if(type_user_id === 3){
            first_name = "";
            last_name = "";
        }else{
            full_name_representative = "";
            name_business = "";
        }

        var jsonData = {
            type_user_id: type_user_id,
            role_id: role_id,
            first_name: first_name,
            last_name: last_name,
            full_name_representative: full_name_representative,
            name_business:name_business,
            email: email,
            phone_number: phone_number,
            address: address,
            password: password,
            referral_code: referral_code,
        };

        $.ajax({
            url: '/accounts/users',
            method: 'POST',
            contentType: "application/json",
            data: JSON.stringify(jsonData),

            success: function (response) {
                if(response === "Success"){
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