$(document).ready(function () {
    var table_students = $('#table-students').DataTable({
        ajax: {
            url: '/managements/mng-payments/api/students',
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
            {data:'phone_number'},
            {data:'address'},
            {
                data: 'verify',
                render:function (data){
                    return StatusVerify(data);
                }
            },
            {
                data: 'paid',
                render: function (data, type, full, meta) {
                    if (data == 0) {
                        return 'Unpaid';
                    } else {
                        return 'Has Paid';
                    }
                }
            },
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
    var table_periods = $('#table-periods').DataTable({
        ajax: {
            url: '/managements/mng-payments/api/periods',
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
                    return table_periods.rows().count() > 0 ?
                        meta.row + meta.settings._iDisplayStart + 1 : 0;
                }
            },
            {data: 'number_student'},
            {data: 'student_paid'},
            {data: 'commission'},
            {
                data: 'period_start',
                render: function (data, type, full, meta) {
                    return moment(data).format('YYYY-MM-DD');
                }
            },
            {
                data: 'period_end',
                render: function (data, type, full, meta) {
                    return moment(data).format('YYYY-MM-DD');
                }
            }              
        ],
        "paging": true, 
        "searching": true,
        "ordering": true,
        "order": [[4, 'desc']],
        "drawCallback": function (settings) {
            var api = this.api();
            api.column(0, {
                order: 'applied'
            }).nodes().each(function (cell, i) {
                cell.innerHTML = i + 1;
            });
        }
    });

    var table_history_pay = $('#table-history-pay').DataTable({
        ajax: {
            url: '/managements/mng-payments/api/history-pay',
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
                    return table_history_pay.rows().count() > 0 ?
                        meta.row + meta.settings._iDisplayStart + 1 : 0;
                }
            },
            {data: 'commission_total'},
            {data: 'commission_paid'},
            {data: 'commission_debt'},
            {data: 'description'},               
            {
                data: 'created_at',
                render:function (data){
                    return formatTime(data);
                },
            },
        ],

        "paging": true, 
        "searching": true,
        "ordering": true,
        "order": [[5, 'desc']],
        "drawCallback": function (settings) {
            var api = this.api();
            api.column(0, {
                order: 'applied'
            }).nodes().each(function (cell, i) {
                cell.innerHTML = i + 1;
            });
        }
    });


    function formatTime(dateString) {
        const date = new Date(dateString); // Assuming the dateString is in a format recognized by Date.parse()

        const day = String(date.getDate()).padStart(2, '0');
        const month = String(date.getMonth() + 1).padStart(2, '0'); // Months are 0-based
        const year = String(date.getFullYear()); // Get last 2 digits of year

        const hours = String(date.getHours()).padStart(2, '0');
        const minutes = String(date.getMinutes()).padStart(2, '0');
        const seconds = String(date.getSeconds()).padStart(2, '0');

        return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
    }

    function StatusVerify(verify){
        if(verify === true){
            return "Verified";
        }
        return "Not verified";
    }
});