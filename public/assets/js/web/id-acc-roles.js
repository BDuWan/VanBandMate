$(document).ready(function () {
    $.fn.dataTable.ext.errMode = 'none';
    var table_roles = $('#table-roles').DataTable({
        "processing": true,
        "serverSide": true,
        ajax: {
            url: '/accounts/roles/api',
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
                    return table_roles.rows().count() > 0 ?
                        meta.row + meta.settings._iDisplayStart + 1 : 0;
                }
            },
            {data: 'name'},
            {
                data: 'role_id',
                render: function (data, type, full, row) {
                    if (type === 'display') {
                        if(full.role_id > 5){
                            return `<div class="row" style="margin:0;">
                                          <div class="col-sm-6 t-center">
                                            <a href="/accounts/roles/${data}">
                                              <i class="fas fa-edit cl-ed" aria-hidden="true" title="Edit"></i>
                                            </a>
                                          </div>
                                          <div class="col-sm-6 t-center">
                                            <a>
                                              <i class="fas fa-trash cl-del item-delete" data-delete="${data}" aria-hidden="true" title="Delete"></i>
                                            </a>
                                          </div>
                                        </div>`;
                        }else{
                            return '<td></td>';
                        }

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

    $('#table-roles tbody').on('click', '.item-delete', function() {
        var RoleID = $(this).data('delete');
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
                    url: '/accounts/roles/' + RoleID,
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
                            table_roles.ajax.reload();
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
});