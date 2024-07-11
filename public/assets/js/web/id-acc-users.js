$(document).ready(function () {
    $.fn.dataTable.ext.errMode = 'none';
    var table_admin = $('#table-admin').DataTable({
        "processing": true,
        "serverSide": true,
        ajax: {
            url: '/accounts/users/api/admin',
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
                    return table_admin.rows().count() > 0 ?
                        meta.row + meta.settings._iDisplayStart + 1 : 0;
                }
            },
            {data: 'email'},
            {data: 'first_name'},
            {data: 'last_name'},
            {data:'Role.name'},
            {
                data: 'verify',
                render:function (data){
                    return StatusVerify(data);
                }
            },
            {
                data: 'state',
                render: function (data, type, full, row) {
                    if (type === 'display') {
                        if(full.user_id > 1){
                            return `<div class="check-box">
                                            <input type="checkbox" class="item-state" data-state="${full.user_id}" ${data ? 'checked' : ''}>
                                        </div>`;
                        }else{
                            return '<td></td>';
                        }

                    }
                    return data;
                }
            },
            {
                data: 'user_id',
                render: function (data, type, full, row) {
                    if (type === 'display') {
                        if(full.user_id > 1){
                            return `<div class="row" style="margin:0;">
                                          <div class="col-sm-6 t-center">
                                            <a href="/accounts/users/${data}">
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

    var table_sales = $('#table-sales').DataTable({
        "processing": true,
        "serverSide": true,
        ajax: {
            url: '/accounts/users/api/sales',
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
                    return table_sales.rows().count() > 0 ?
                        meta.row + meta.settings._iDisplayStart + 1 : 0;
                }
            },
            {data: 'email'},
            {data: 'first_name'},
            {data: 'last_name'},
            {data:'Role.name'},
            {
                data: 'verify',
                render:function (data){
                    return StatusVerify(data);
                }
            },
            {
                data: 'state',
                render: function (data, type, full, row) {
                    if (type === 'display') {
                        if(full.user_id > 1){
                            return `<div class="check-box">
                                            <input type="checkbox" class="item-state" data-state="${full.user_id}" ${data ? 'checked' : ''}>
                                        </div>`;
                        }else{
                            return '<td></td>';
                        }

                    }
                    return data;
                }
            },
            {
                data: 'user_id',
                render: function (data, type, full, row) {
                    if (type === 'display') {
                        if(full.user_id > 1){
                            return `<div class="row" style="margin:0;">
                                          <div class="col-sm-6 t-center">
                                            <a href="/accounts/users/${data}">
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

    var table_business = $('#table-business').DataTable({
        "processing": true,
        "serverSide": true,
        ajax: {
            url: '/accounts/users/api/business',
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
                    return table_business.rows().count() > 0 ?
                        meta.row + meta.settings._iDisplayStart + 1 : 0;
                }
            },
            {data: 'email'},
            {data: 'name_business'},
            {data: 'full_name_representative'},
            {data:'Role.name'},
            {
                data: 'verify',
                render:function (data){
                    return StatusVerify(data);
                }
            },
            {
                data: 'state',
                render: function (data, type, full, row) {
                    if (type === 'display') {
                        if(full.user_id > 1){
                            return `<div class="check-box">
                                            <input type="checkbox" class="item-state" data-state="${full.user_id}" ${data ? 'checked' : ''}>
                                        </div>`;
                        }else{
                            return '<td></td>';
                        }

                    }
                    return data;
                }
            },
            {
                data: 'user_id',
                render: function (data, type, full, row) {
                    if (type === 'display') {
                        if(full.user_id > 1){
                            return `<div class="row" style="margin:0;">
                                          <div class="col-sm-6 t-center">
                                            <a href="/accounts/users/${data}">
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

    var table_instructors = $('#table-instructors').DataTable({
        "processing": true,
        "serverSide": true,
        ajax: {
            url: '/accounts/users/api/instructors',
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
                    return table_instructors.rows().count() > 0 ?
                        meta.row + meta.settings._iDisplayStart + 1 : 0;
                }
            },
            {data: 'email'},
            {data: 'first_name'},
            {data: 'last_name'},
            {data:'Role.name'},
            {
                data: 'verify',
                render:function (data){
                    return StatusVerify(data);
                }
            },
            {
                data: 'state',
                render: function (data, type, full, row) {
                    if (type === 'display') {
                        if(full.user_id > 1){
                            return `<div class="check-box">
                                            <input type="checkbox" class="item-state" data-state="${full.user_id}" ${data ? 'checked' : ''}>
                                        </div>`;
                        }else{
                            return '<td></td>';
                        }

                    }
                    return data;
                }
            },
            {
                data: 'user_id',
                render: function (data, type, full, row) {
                    if (type === 'display') {
                        if(full.user_id > 1){
                            return `<div class="row" style="margin:0;">
                                          <div class="col-sm-6 t-center">
                                            <a href="/accounts/users/${data}">
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

    var table_students = $('#table-students').DataTable({
        "processing": true,
        "serverSide": true,
        ajax: {
            url: '/accounts/users/api/students',
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
                    return table_students.rows().count() > 0 ?
                        meta.row + meta.settings._iDisplayStart + 1 : 0;
                }
            },
            {data: 'email'},
            {data: 'first_name'},
            {data: 'last_name'},
            {data:'Role.name'},
            {
                data: 'verify',
                render:function (data){
                    return StatusVerify(data);
                }
            },
            {
                data: 'state',
                render: function (data, type, full, row) {
                    if (type === 'display') {
                        if(full.user_id > 1){
                            return `<div class="check-box">
                                            <input type="checkbox" class="item-state" data-state="${full.user_id}" ${data ? 'checked' : ''}>
                                        </div>`;
                        }else{
                            return '<td></td>';
                        }

                    }
                    return data;
                }
            },
            {
                data: 'user_id',
                render: function (data, type, full, row) {
                    if (type === 'display') {
                        if(full.user_id > 1){
                            return `<div class="row" style="margin:0;">
                                          <div class="col-sm-6 t-center">
                                            <a href="/accounts/users/${data}">
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

    $('#table-admin tbody').on('click', '.item-delete', function() {
        var UserID = $(this).data('delete');
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
                    url: '/accounts/users/' + UserID,
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
                            table_admin.ajax.reload();
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

    $('#table-sales tbody').on('click', '.item-delete', function() {
        var UserID = $(this).data('delete');
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
                    url: '/accounts/users/' + UserID,
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
                            table_sales.ajax.reload();
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

    $('#table-business tbody').on('click', '.item-delete', function() {
        var UserID = $(this).data('delete');
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
                    url: '/accounts/users/' + UserID,
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
                            table_business.ajax.reload();
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

    $('#table-instructors tbody').on('click', '.item-delete', function() {
        var UserID = $(this).data('delete');
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
                    url: '/accounts/users/' + UserID,
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
                            table_instructors.ajax.reload();
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

    $('#table-students tbody').on('click', '.item-delete', function() {
        var UserID = $(this).data('delete');
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
                    url: '/accounts/users/' + UserID,
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
                            table_students.ajax.reload();
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

    $('#table-admin tbody').on('change', '.item-state', function () {

        var isChecked = $(this).prop('checked');
        var id = $(this).data("state");
        $.ajax({
            url: '/accounts/users/state/' + id,
            method: 'PUT',
            dataType:'json',
            data: {
                id: id,
                state: isChecked
            },
            success: function (response) {
                if (response === "Success") {
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
                    })
                }
            },
            error: function (error) {
                console.error(error);
            }
        });
    });

    $('#table-sales tbody').on('change', '.item-state', function () {

        var isChecked = $(this).prop('checked');
        var id = $(this).data("state");
        $.ajax({
            url: '/accounts/users/state/' + id,
            method: 'PUT',
            dataType:'json',
            data: {
                id: id,
                state: isChecked
            },
            success: function (response) {
                if (response === "Success") {
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
                    })
                }
            },
            error: function (error) {
                console.error(error);
            }
        });
    });

    $('#table-business tbody').on('change', '.item-state', function () {

        var isChecked = $(this).prop('checked');
        var id = $(this).data("state");
        $.ajax({
            url: '/accounts/users/state/' + id,
            method: 'PUT',
            dataType:'json',
            data: {
                id: id,
                state: isChecked
            },
            success: function (response) {
                if (response === "Success") {
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
                    })
                }
            },
            error: function (error) {
                console.error(error);
            }
        });
    });

    $('#table-instructors tbody').on('change', '.item-state', function () {

        var isChecked = $(this).prop('checked');
        var id = $(this).data("state");
        $.ajax({
            url: '/accounts/users/state/' + id,
            method: 'PUT',
            dataType:'json',
            data: {
                id: id,
                state: isChecked
            },
            success: function (response) {
                if (response === "Success") {
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
                    })
                }
            },
            error: function (error) {
                console.error(error);
            }
        });
    });

    $('#table-students tbody').on('change', '.item-state', function () {

        var isChecked = $(this).prop('checked');
        var id = $(this).data("state");
        $.ajax({
            url: '/accounts/users/state/' + id,
            method: 'PUT',
            dataType:'json',
            data: {
                id: id,
                state: isChecked
            },
            success: function (response) {
                if (response === "Success") {
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
                    })
                }
            },
            error: function (error) {
                console.error(error);
            }
        });
    });

    function StatusVerify(verify){
        if(verify === true){
            return "Verified";
        }
        return "Not verified";
    }

});