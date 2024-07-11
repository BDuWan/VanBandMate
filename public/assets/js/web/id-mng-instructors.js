$(document).ready(function () {
    $.fn.dataTable.ext.errMode = 'none';
    var table_instructor = $('#table-instructor').DataTable({
        "processing": true,
        "serverSide": true,
        ajax: {
            url: '/managements/instructors/api',
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
                    return table_instructor.rows().count() > 0 ?
                        meta.row + meta.settings._iDisplayStart + 1 : 0;
                }
            },
            {data: 'first_name'},
            {data: 'last_name'},
            {data: 'email'},
            {data:'phone_number'},
            {data: 'address'},
            {
                data: 'user_id',
                render: function (data, type, full, row) {
                    if (type === 'display') {

                        return `<div class="row" style="margin:0;">
                                          <div class="col-sm-12 t-center">
                                            <a href="/managements/instructors/${data}">
                                              <i class="fas fa-info-circle cl-ed" aria-hidden="true" title="Information"></i>
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
        lengthMenu: [[25, 50, 100], ['25', '50', '100']],
        // language: {
        //     lengthMenu: "{{.Lag.Show}}" + '_MENU_' + "{{.Lag.Index}}",
        //     infoFiltered: "({{.Lag.InTotal}} _MAX_ {{.Lag.Index}})",
        //     info: "{{.Lag.Show}} _START_ {{.Lag.Arrive}} _END_ {{.Lag.BelongTo}} _TOTAL_ {{.Lag.Index}}",
        //     infoEmpty: "{{.Lag.Show}} 0 {{.Lag.Arrive}} 0 {{.Lag.BelongTo}} 0 {{.Lag.Index}}",
        //     search: "{{.Lag.Search}}",
        //     paginate: {
        //         next: '{{.Lag.Next}}',
        //         previous: '{{.Lag.Previous}}'
        //     },
        // }
    });
});