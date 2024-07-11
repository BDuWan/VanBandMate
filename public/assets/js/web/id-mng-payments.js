$(document).ready(function () {
    // $.fn.dataTable.ext.errMode = 'none';
    // var table_price = $('#table-price').DataTable({
    //     "processing": true,
    //     "serverSide": true,
    //     ajax: {
    //         url: '/managements/mng-payments/api/prices',
    //         type: 'POST',
    //         "contentType": "application/json",
    //         "data": function (d) {
    //             return JSON.stringify(d);
    //         }

    //     },
    //     columns: [
    //         {
    //             "data": null,
    //             "searchable": false,
    //             "orderable": false,
    //             "render": function (data, type, full, meta) {
    //                 return table_price.rows().count() > 0 ?
    //                     meta.row + meta.settings._iDisplayStart + 1 : 0;
    //             }
    //         },
    //         {data: 'price'},
    //         {data: 'commission'},
    //         {data: 'description'},
    //         {
    //             data: 'start_time',
    //             render: function (data, type, full, row) {
    //                 if (type === 'display') {
    //                     if(full.default === false){
    //                         return `<td>`+formatTime(data)+`</td>`;
    //                     }
    //                 }
    //                 return "<td></td>";
    //             }
    //         },
    //         {
    //             data: 'end_time',
    //             render: function (data, type, full, row) {
    //                 if (type === 'display') {
    //                     if(full.default === false){
    //                         return `<td>`+formatTime(data)+`</td>`;
    //                     }

    //                 }
    //                 return "<td></td>";
    //             }
    //         },
    //         {
    //             data: 'price_program_id',
    //             render: function (data, type, full, row) {
    //                 if (type === 'display') {
    //                     return `<div class="row" style="margin:0;">
    //                                       <div class="col-sm-6 t-center">
    //                                         <a href="/managements/mng-payments/prices/${data}">
    //                                           <i class="fas fa-edit cl-ed" aria-hidden="true" title="Edit"></i>
    //                                         </a>
    //                                       </div>
    //                                       <div class="col-sm-6 t-center">
    //                                         <a>
    //                                           <i class="fas fa-trash cl-del item-delete" data-delete="${data}" aria-hidden="true" title="Delete"></i>
    //                                         </a>
    //                                       </div>                                       
    //                                     </div>`;
    //                 }
    //                 return data;
    //             }
    //         },

    //     ],

    //     "drawCallback": function (settings) {
    //         var api = this.api();
    //         var startIndex = api.context[0]._iDisplayStart;
    //         api.column(0, {page: 'current'}).nodes().each(function (cell, i) {
    //             cell.innerHTML = startIndex + i + 1;
    //         });
    //     },
    //     "ordering": true,
    //     lengthMenu: [[5, 10], ['5', '10']],
    // });

    var table_pro_sale = $('#table-pro-sale').DataTable({
        "processing": true,
        "serverSide": true,
        ajax: {
            url: '/managements/mng-payments/api/pro-sales',
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
                    return table_pro_sale.rows().count() > 0 ?
                        meta.row + meta.settings._iDisplayStart + 1 : 0;
                }
            },
            {data: 'User.first_name'},
            {data: 'User.last_name'},
            // {data: 'User.name_business'},
            // {data: 'User.full_name_representative'},
            {data: 'commission_total'},
            {data: 'commission_paid'},
            {data: 'commission_debt'},
            {
                data: 'User.user_id',
                render: function (data, type, full, row) {
                    if (type === 'display') {

                        return `<div class="row" style="margin:0;">
                                       
                                          <div class="col-sm-12 t-center">
                                            <a href="/managements/sale-business/${full.user_id}">
                                              <i class="fas fa-info-circle cl-ed" aria-hidden="true" title="Info"></i>
                                            </a>
                                          </div>                                       
                                        </div>`;
                    }
                    return data;
                }
            }

        ],

        "drawCallback": function (settings) {
            var api = this.api();
            var startIndex = api.context[0]._iDisplayStart;
            api.column(0, {page: 'current'}).nodes().each(function (cell, i) {
                cell.innerHTML = startIndex + i + 1;
            });
        },
        "ordering": true,
        lengthMenu: [[5, 10], ['5', '10']],
    });

    $('#table-price tbody').on('click', '.item-delete', function() {
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
    function formatTime(dateString) {
        const date = new Date(dateString); // Assuming the dateString is in a format recognized by Date.parse()

        const day = String(date.getDate()).padStart(2, '0');
        const month = String(date.getMonth() + 1).padStart(2, '0'); // Months are 0-based
        const year = String(date.getFullYear()); // Get last 2 digits of year

        const hours = String(date.getHours()).padStart(2, '0');
        const minutes = String(date.getMinutes()).padStart(2, '0');
        const seconds = String(date.getSeconds()).padStart(2, '0');

        return `${day}-${month}-${year} ${hours}:${minutes}:${seconds}`;
    }
});