$(document).ready(function () {
    $.fn.dataTable.ext.errMode = 'none';
    var table_student = $('#table-student').DataTable({
        ajax: {
            url: '/managements/students/api',
            type: 'POST',
            "contentType": "application/json",
            "data": function (d) {
                return JSON.stringify(d);
            }

        },
        columns: [
            {
                "data": null,
                "searchable": false,
                "orderable": false,
                "render": function (data, type, full, meta) {
                    return table_student.rows().count() > 0 ?
                        meta.row + meta.settings._iDisplayStart + 1 : 0;
                }
            },
            {data: 'first_name'},
            {data: 'last_name'},
            {data: 'email'},
            {data:'phone_number'},
            {data: 'address'},
            {
                data: 'verify',
                render: function (data, type, full, meta) {
                    if (data == 0) {
                        return 'Not verified';
                    } else {
                        return 'Verified';
                    }
                }
            },
            {
                data: 'paid',
                render: function (data, type, full, meta) {
                    if (data == 0) {
                        return `<div class="row" style="margin:0;">
                                    <div class="col-sm-12 t-center">
                                        <button class="btn btn-primary item-state" data-state="${full.user_id}">Confirm Payment</button>
                                    </div>
                                </div>`;
                    } else {
                        return 'Has Paid';
                    }
                }
            },
            {
                data: 'user_id',
                "searchable": false,
                "orderable": false,
                render: function (data, type, full, row) {
                    if (type === 'display') {

                        return `<div class="row" style="margin:0;">
                                          <div class="col-sm-12 t-center">
                                            <a href="/managements/students/${data}">
                                              <i class="fas fa-info-circle cl-ed" aria-hidden="true" title="Information"></i>
                                            </a>
                                          </div>                                       
                                        </div>`;


                    }
                    return data;
                }
            }

        ],

        "paging": true, 
        "searching": true,
        "ordering": true,
        "order": [],
        "drawCallback": function (settings) {
            var api = this.api();
            api.column(0, {
                order: 'applied'
            }).nodes().each(function (cell, i) {
                cell.innerHTML = i + 1;
            });
        }
    });

    $('#table-student tbody').on('click', '.item-delete', function() {
        var PriceID = $(this).data('delete');
        swal({
            title: 'Confirm Delete ?',
            icon: 'warning',
            buttons: {
                cancel: {
                    text: "Cancel",
                    value: null,
                    visible: true,
                    className: "btn btn-danger",
                    closeModal: true,
                },
                confirm: {
                    text: "OK",
                    value: true,
                    visible: true,
                    className: "btn btn-primary",
                    closeModal: true
                }
            }
        }).then((result=>{
            if(result){
                $.ajax({
                    url: '/managements/mng-payments/prices/' + PriceID,
                    type: 'DELETE',
                    success: function (response) {
                        if(response === "Success"){
                            swal({
                                title: 'Successfully !',
                                icon: 'success',
                                button: {
                                    text: "Close",
                                    value: true,
                                    visible: true,
                                    className: "btn btn-primary"
                                }
                            })
                            table_price.ajax.reload();
                        }else {
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
                    error: function (xhr, status, error) {
                        console.error("Error deleting row:", status, error);
                    }
                });
            }

        }))
    });

    $('#table-student tbody').on('click', '.item-state', function () {
        var id = $(this).data("state");
        
        swal({
            title: "Are you sure?",
            text: "Do you really want to confirm this payment?",
            icon: "warning",
            buttons: {
                cancel: {
                    text: "Cancel",
                    value: null,
                    visible: true,
                    className: "btn btn-secondary",
                    closeModal: true,
                },
                confirm: {
                    text: "Yes, confirm it!",
                    value: true,
                    visible: true,
                    className: "btn btn-primary",
                    closeModal: false
                }
            }
        }).then((willConfirm) => {
            if (willConfirm) {
                $.ajax({
                    url: '/managements/students/confirm-payment/' + id,
                    method: 'PUT',
                    dataType: 'json',
                    success: function (response) {
                        if (response === "Success") {
                            swal({
                                title: 'Successfully!',
                                icon: 'success',
                                button: {
                                    text: "Close",
                                    value: true,
                                    visible: true,
                                    className: "btn btn-primary"
                                }
                            }).then(() => {
                                $('#table-student').DataTable().ajax.reload();
                            });
                        } else if(response === "Success1"){
                            swal({
                                title: 'This student has not entered the referral code. Payment confirmation will not add commission to any sales director!',
                                icon: 'success',
                                button: {
                                    text: "Close",
                                    value: true,
                                    visible: true,
                                    className: "btn btn-primary"
                                }
                            }).then(() => {
                                $('#table-student').DataTable().ajax.reload();
                            });
                        } else {
                            swal({
                                title: response,
                                icon: 'warning',
                                button: {
                                    text: "Close",
                                    value: true,
                                    visible: true,
                                    className: "btn btn-danger"
                                }
                            });
                        }
                    },
                    error: function (error) {
                        console.error(error);
                    }
                });
            }
        });
    });
    
});