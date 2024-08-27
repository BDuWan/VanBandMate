var currentPage = 1
function setupPagination(totalItems, itemsPerPage, listSelector, paginationSelector) {
    var $pagination = $(paginationSelector);

    $pagination.empty();
    var numPages = Math.ceil(totalItems / itemsPerPage);

    function renderPageItem(page) {
        return $('<div class="page-item"><a href="#" class="page-link"></a></div>')
            .find('a')
            .text(page)
            .end()
            .appendTo($pagination)
            .on('click', function(e) {
                e.preventDefault();
                renderContractList(page, itemsPerPage, listSelector, paginationSelector);
                setupPaginationControls(page);
            });
    }

    function renderEllipsis() {
        $('<div class="page-item"><strong>...</strong></div>').appendTo($pagination);
    }

    // renderPageItem(1);
    function setupPaginationControls(currentPage) {
        $pagination.find('a').removeClass('active');
        $pagination.empty();

        renderPageItem(1);

        if (currentPage > 3) {
            renderEllipsis();
        }

        var start = Math.max(2, currentPage - 1);
        var end = Math.min(numPages - 1, currentPage + 1);

        for (var i = start; i <= end; i++) {
            renderPageItem(i);
        }

        if (currentPage < numPages - 2) {
            renderEllipsis();
        }

        if (numPages > 1) {
            renderPageItem(numPages);
        }

        $pagination.find('a').filter(function() {
            return $(this).text() == currentPage;
        }).addClass('active');
    }

    // Hiển thị trang đầu tiên
    renderContractList(1, itemsPerPage, listSelector, paginationSelector);
    setupPaginationControls(1);
}

function renderContractList(page, itemsPerPage, listSelector, paginationSelector) {
    currentPage = page
    $.ajax({
        url: '/contract/api',
        type: 'GET',
        dataType: 'json',
        data: {
            page: page,
            itemsPerPage: itemsPerPage
        },
        success: function(response) {
            var $list = $(listSelector);
            var $pagination = $(paginationSelector);
            var data = response.data;
            var user_id = response.user_id;
            var totalItems = response.totalItems; // Tổng số lượng phần tử từ backend

            // Xóa danh sách cũ
            $list.empty();

            // Render dữ liệu danh sách
            $.each(data, function(index, item) {
                var listItem = renderListItem(item, user_id)
                $list.append(listItem);
            });
        },
        error: function(xhr, status, error) {
            console.error("Đã xảy ra lỗi ", error);
        }
    });
}

function setupPaginationFilter(totalItems, itemsPerPage, listSelector, paginationSelector, jsonData) {
    var $pagination = $(paginationSelector);

    $pagination.empty();
    var numPages = Math.ceil(totalItems / itemsPerPage);

    function renderPageFilterItem(page) {
        return $('<div class="page-item"><a href="#" class="page-link"></a></div>')
            .find('a')
            .text(page)
            .end()
            .appendTo($pagination)
            .on('click', function(e) {
                e.preventDefault();
                renderListFilter(page, itemsPerPage, listSelector, paginationSelector, jsonData);
                setupPaginationFilterControls(page);
            });
    }

    function renderFilterEllipsis() {
        $('<div class="page-item"><strong>...</strong></div>').appendTo($pagination);
    }

    // renderPageFilterItem(1);
    function setupPaginationFilterControls(currentPage) {
        $pagination.find('a').removeClass('active');
        $pagination.empty();

        renderPageFilterItem(1);

        if (currentPage > 3) {
            renderFilterEllipsis();
        }

        var start = Math.max(2, currentPage - 1);
        var end = Math.min(numPages - 1, currentPage + 1);

        for (var i = start; i <= end; i++) {
            renderPageFilterItem(i);
        }

        if (currentPage < numPages - 2) {
            renderFilterEllipsis();
        }

        if (numPages > 1) {
            renderPageFilterItem(numPages);
        }

        $pagination.find('a').filter(function() {
            return $(this).text() == currentPage;
        }).addClass('active');
    }

    // Hiển thị trang đầu tiên
    renderListFilter(1, itemsPerPage, listSelector, paginationSelector, jsonData);
    setupPaginationFilterControls(1);
}

function renderListFilter(page, itemsPerPage, listSelector, paginationSelector, jsonData) {
    var extendedJsonData = { ...jsonData };
    extendedJsonData.page = page;
    $.ajax({
        url: '/contract/api/filter',
        method: 'POST',
        contentType: "application/json",
        data: JSON.stringify(extendedJsonData),
        success: function(response) {
            var $list = $(listSelector);
            var $pagination = $(paginationSelector);
            var data = response.data;
            var user_id = response.user_id;
            var totalItems = response.totalItems; // Tổng số lượng phần tử từ backend

            // Xóa danh sách cũ
            $list.empty();

            // Render dữ liệu danh sách
            $.each(data, function(index, item) {
                var listItem = renderListItem(item, user_id)
                $list.append(listItem);
            });
        },
        error: function(xhr, status, error) {
            console.error("Error fetching contract news: ", error);
        }
    });
}

function renderListItem(item, user_id) {
    var formatedDate = formatDate(item.date);
    var formatedPrice = formatPrice(item.price);
    var formatedRelativeDate = formatRelativeDate(item.created_at);
    const backgroundColor = item.status === 1
        ? 'background-color: #FFFFFF;'
        : item.status === 0
            ? 'background-color: #99FFFF;'
            : 'background-color: #FFFF99;';
    let statusMessage = null;
    let btnHtml = `<a href="#" class="dropdown-item request-delete-item" data-id="${item.contract_id}">Yêu cầu hủy hợp đồng</a>`

    if (item.status === 2) {
        if (item.request_delete_by !== user_id) {
            statusMessage = "Đối tác yêu cầu hủy";
            btnHtml = `<a href="#" class="dropdown-item confirm-delete-item" data-id="${item.contract_id}">Chấp nhận hủy hợp đồng</a>`
        } else if (item.request_delete_by === user_id) {
            statusMessage = "Bạn đã yêu cầu hủy";
            btnHtml = `<a href="#" class="dropdown-item cancel-delete-item" data-id="${item.contract_id}">Thu hồi yêu cầu hủy</a>`
        }
    } else if (item.status === 0){
        statusMessage = "Đã hoàn thành";
        btnHtml =""
    }
    let statusMessageHtml = statusMessage ? `<p class="list-group-item-text">${statusMessage}</p>` : '';

    return `
        <li class="list-group-item" style="${backgroundColor}">
            <div class="list-group-item-figure d-flex flex-column">                
                <div class="avatar">
                    <img src="/assets/img/avatar/${item.ChuLoaDai.image}" alt="..." class="avatar-img rounded-circle">
                </div>                       
                <div class="avatar">
                    <img src="/assets/img/avatar/${item.NhacCong.image}" alt="..." class="avatar-img rounded-circle">
                </div>              
            </div>
            <div class="list-group-item-body pl-3 pl-md-4">
                <div class="row">
                    <div class="col-12 col-lg-10">
                        <strong class="list-group-item-text">${formatedDate}</strong>
                        <p class="list-group-item-text break-word">
                            Chủ loa đài: <strong>${item.ChuLoaDai.last_name} ${item.ChuLoaDai.first_name}</strong>
                        </p>   
                        <p class="list-group-item-text break-word">
                            Nhạc công: <strong>${item.NhacCong.last_name} ${item.NhacCong.first_name}</strong>
                        </p>                       
                        <p class="list-group-item-text break-word">
                            ${item.ward.full_name}, ${item.district.full_name}, ${item.province.name}
                        </p>
                        <strong class="list-group-item-text">${formatedPrice}</strong>
                    </div>
                    <div class="col-12 col-lg-2 text-lg-right">
                        <p class="list-group-item-text">Thời gian tạo: ${formatedRelativeDate}</p>                     
                        <strong>${statusMessageHtml}</strong>
                    </div>
                </div>
            </div>
            <div class="list-group-item-figure">
                <div class="dropdown">
                    <button class="btn-dropdown" data-toggle="dropdown">
                        <i class="fa fa-ellipsis-v"></i>
                    </button>
                    <div class="dropdown-arrow"></div>
                    <div class="dropdown-menu dropdown-menu-right">
                        <a class="dropdown-item view-detail" data-id="${item.contract_id}">Xem chi tiết</a>
                        ${btnHtml}                       
                    </div>
                </div>
            </div>
        </li>
    `;
}

$(document).on('click', '.view-detail', function(event) {
    event.preventDefault();
    var contractId = $(this).data('id');
    $('#detailContractModal').modal('show');
    $.ajax({
        url: '/contract/api/detail/'+contractId,
        method: 'GET',
        success: function(response) {
            if(response.message !== "success"){
                swal("", response.message, {
                    icon : "error",
                    buttons: {
                        confirm: {
                            className : 'btn btn-danger'
                        }
                    },
                });
                return false;
            }
            var contract = response.data;
            var formattedCreatedAt = moment(contract.created_at, "YYYY-MM-DD HH:mm:ss").format("DD/MM/YYYY HH:mm:ss");
            var formattedDate = moment(contract.date, "YYYY-MM-DD").format("DD/MM/YYYY");
            var formattedPrice = formatPrice(contract.price);

            // Gán giá trị vào các trường trong modal
            // $('#detailContractModal #name').text(hiringNews.User.last_name + " " + hiringNews.User.first_name);
            // $('#detailContractModal #role').text(hiringNews.User.last_name);

            $('#detailContractModal #detail-name-chuloadai').text(contract.ChuLoaDai.last_name + " " + contract.ChuLoaDai.first_name);
            $('#detailContractModal #detail-phone-chuloadai').text("SDT: " + contract.ChuLoaDai.phone_number);
            $('#detailContractModal #detail-email-chuloadai').text("Email: " + contract.ChuLoaDai.email);

            $('#detailContractModal #detail-name-nhaccong').text(contract.NhacCong.last_name + " " + contract.NhacCong.first_name);
            $('#detailContractModal #detail-phone-nhaccong').text("SDT: " + contract.NhacCong.phone_number);
            $('#detailContractModal #detail-email-nhaccong').text("Email: " + contract.NhacCong.email);

            $('#detailContractModal #detail-create-at').text(formattedCreatedAt);
            $('#detailContractModal #detail-date').text(formattedDate);
            $('#detailContractModal #detail-price').text(formattedPrice);

            $('#detailContractModal #detail-address').text(
                contract.address_detail + ", " + contract.ward.full_name + ", " + contract.district.full_name + ", " + contract.province.name
            );

            $('html').removeClass('topbar_open');
            $('#detailContractModal').modal('show');

        }
    });
});

$(document).on('click', '.request-delete-item', function(e) {
    e.preventDefault();
    let contractId = $(this).data('id');

    $.ajax({
        url: '/contract/request-delete',
        type: 'POST',
        data: { id: contractId },
        success: function(response) {
            renderContractList(currentPage, itemsPerPage,  '#list-contract', '#pagination-contract')
            swal("", response.message, {
                icon : response.icon,
                buttons: {
                    confirm: {
                        className : 'btn btn-danger'
                    }
                },
            });
        },
        error: function(xhr, status, error) {
            // Xử lý lỗi
            alert('Đã xảy ra lỗi khi yêu cầu hủy hợp đồng.');
        }
    });
});

$(document).on('click', '.confirm-delete-item', function(e) {
    e.preventDefault();
    let contractId = $(this).data('id');

    $.ajax({
        url: '/contract/confirm-delete',
        type: 'POST',
        data: { id: contractId },
        success: function(response) {
            renderContractList(currentPage, itemsPerPage,  '#list-contract', '#pagination-contract')
            swal("", response.message, {
                icon : response.icon,
                buttons: {
                    confirm: {
                        className : 'btn btn-danger'
                    }
                },
            });
        },
        error: function(xhr, status, error) {
            // Xử lý lỗi
            alert('Đã xảy ra lỗi khi chấp nhận hủy hợp đồng.');
        }
    });
});

$(document).on('click', '.cancel-delete-item', function(e) {
    e.preventDefault();
    let contractId = $(this).data('id');

    $.ajax({
        url: '/contract/cancel-delete',
        type: 'POST',
        data: { id: contractId },
        success: function(response) {
            renderContractList(currentPage, itemsPerPage,  '#list-contract', '#pagination-contract')
            swal("", response.message, {
                icon : response.icon,
                buttons: {
                    confirm: {
                        className : 'btn btn-danger'
                    }
                },
            });
        },
        error: function(xhr, status, error) {
            // Xử lý lỗi
            alert('Đã xảy ra lỗi khi thu hồi yêu cầu hủy hợp đồng.');
        }
    });
});

