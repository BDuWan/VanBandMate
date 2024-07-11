$(document).ready(function () {
    $.fn.dataTable.ext.errMode = 'none';
    var table_sale_business = $('#table-sale-business').DataTable({
        ajax: {
            url: '/managements/sale-business/api',
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
                    return table_sale_business.rows().count() > 0 ?
                        meta.row + meta.settings._iDisplayStart + 1 : 0;
                }
            },
            {data: 'User.first_name'},
            {data: 'User.last_name'},
            // {data: 'name_business'},
            // {data: 'full_name_representative'},
            {data: 'User.email'},
            {data: 'User.phone_number'},
            {data: 'User.address'},
            {data: 'commission_paid'},
            {data: 'commission_debt'},
            {
                data: 'user_id',
                render: function (data, type, full, row) {
                    if (type === 'display') {

                        return `<div class="row" style="margin:0;">
                                          <div class="col-sm-12 t-center">
                                            <a href="/managements/sale-business/${data}">
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
});