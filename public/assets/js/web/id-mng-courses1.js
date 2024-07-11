$(document).ready(function () {
    $.fn.dataTable.ext.errMode = 'none';
    var table_courses = $('#table-courses').DataTable({
        "processing": true,
        "serverSide": true,
        ajax: {
            url: '/managements/mng-courses/api',
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
                    return table_courses.rows().count() > 0 ?
                        meta.row + meta.settings._iDisplayStart + 1 : 0;
                }
            },
            {
                data: 'image',
                render: function (data, type, full, row) {
                    if (type === 'display') {

                        return `<img src="/assets/images/courses/${data}" width="150" height="120">`;


                    }
                    return data;
                }
            },
            {data: 'title'},
            {data: 'price'},
            {data: 'description'},
            {
                data: 'course_id',
                render: function (data, type, full, row) {
                    if (type === 'display') {
                        return `<div class="row" style="margin:0;">
                                          <div class="col-sm-6 t-center">
                                            <a href="/managements/mng-courses/${data}">
                                              <i class="fas fa-edit cl-ed" aria-hidden="true" title="Edit"></i>
                                            </a>
                                          </div>
                                          <div class="col-sm-6 t-center">
                                            <a>
                                              <i class="fas fa-trash cl-del item-delete" data-delete="${data}" aria-hidden="true" title="Delete"></i>
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

    $('#table-courses tbody').on('click', '.item-delete', function() {
        var CourseID = $(this).data('delete');
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
                    url: '/managements/mng-courses/' + CourseID,
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
                            table_courses.ajax.reload();
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