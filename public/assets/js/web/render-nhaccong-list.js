var currentPage = 1
var jsonDataFilter = {}
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
    jsonDataFilter = { ...jsonData }
    var extendedJsonData = { ...jsonData };
    extendedJsonData.page = page;
    $.ajax({
        url: '/hiring/api/find',
        method: 'POST',
        contentType: "application/json",
        data: JSON.stringify(extendedJsonData),
        success: function(response) {
            var $list = $(listSelector);
            var $pagination = $(paginationSelector);
            var data = response.data;
            var hiringNewsId = response.hiringNewsId;
            var totalItems = response.totalItems; // Tổng số lượng phần tử từ backend

            // Xóa danh sách cũ
            $list.empty();

            // Render dữ liệu danh sách
            $.each(data, function(index, item) {
                var listItem = renderListItem(item, hiringNewsId)
                $list.append(listItem);
            });
        },
        error: function(xhr, status, error) {
            console.error("Error fetching hiring news: ", error);
        }
    });
}

function renderListItem(item, hiringNewsId) {
    let backgroundColor = 'background-color: #FFFFFF;';
    let statusMessage = null;
    if (item.invitation_status === 0) {
        backgroundColor = 'background-color: #FFCCFF;';
        statusMessage = "Đã gửi lời mời"
    } else if(item.invitation_status === 1){
        backgroundColor = 'background-color: #99FF99;';
        statusMessage = "Đã được chấp nhận"
    } else if(item.invitation_status === 3){
        backgroundColor = 'background-color: #FFFF99;';
        statusMessage = "Bị từ chối"
    }
    let statusMessageHtml = statusMessage ? `<p class="list-group-item-text">${statusMessage}</p>` : '';

    const inviteButtonHtml = item.invitation_status === 2 || item.invitation_status === 3 || item.invitation_status === 4
        ? `<a href="#" class="dropdown-item invite-item" data-hiring-id="${hiringNewsId}" data-nhaccong-id="${item.user_id}">Gửi lời mời</a>`
        : '';
    const cancelInviteButtonHtml = item.invitation_status === 0
        ? `<a href="#" class="dropdown-item cancel-invite-item" data-hiring-id="${hiringNewsId}" data-nhaccong-id="${item.user_id}">Thu hồi lời mời</a>`
        : '';
    return `
        <li class="list-group-item" style="${backgroundColor}">
            <div class="list-group-item-figure">
                <a href="/info/profile/${item.user_id}" class="user-avatar">
                    <div class="avatar">
                        <img src="/assets/img/avatar/${item.image}" alt="..." class="avatar-img rounded-circle">
                    </div>
                </a>
            </div>
            <div class="list-group-item-body pl-3 pl-md-4">
                <div class="row">
                    <div class="col-12 col-lg-10">
                        <strong class="list-group-item-text break-word">${item.last_name} ${item.first_name}</strong>
                        <br>
                      
                        <p class="list-group-item-text break-word">
                            ${item.ward.full_name}, ${item.district.full_name}, ${item.province.name}
                        </p>
                        <p class="list-group-item-text break-word">
                            SDT: ${item.phone_number}
                        </p>
                    </div>
                    <div class="col-12 col-lg-2 text-lg-right">         
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
                        <a href="/info/profile/${item.user_id}" class="dropdown-item view-detail">Xem trang cá nhân</a>
                        ${inviteButtonHtml}   
                        ${cancelInviteButtonHtml}                                    
                    </div>
                </div>
            </div>
        </li>
    `;
}

$(document).on('click', '.invite-item', function(e) {
    e.preventDefault();
    let nhaccongId = $(this).data('nhaccong-id');
    let hiringNewsId = parseInt($(this).data('hiring-id'),10);

    $.ajax({
        url: '/hiring/invite',
        type: 'POST',
        data: {
            nhaccongId: nhaccongId,
            hiringNewsId: hiringNewsId
        },
        success: function(response) {
            renderListFilter(currentPage, itemsPerPage,  '#list-nhaccong', '#pagination-nhaccong', jsonDataFilter)
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
            alert('Đã xảy ra lỗi khi gửi lời mời.');
        }
    });
});

$(document).on('click', '.cancel-invite-item', function(e) {
    e.preventDefault();
    let nhaccongId = $(this).data('nhaccong-id');
    let hiringNewsId = parseInt($(this).data('hiring-id'),10);

    $.ajax({
        url: '/hiring/cancel-invite',
        type: 'POST',
        data: {
            nhaccongId: nhaccongId,
            hiringNewsId: hiringNewsId
        },
        success: function(response) {
            renderListFilter(currentPage, itemsPerPage,  '#list-nhaccong', '#pagination-nhaccong', jsonDataFilter)
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
            alert('Đã xảy ra lỗi khi thu hồi lời mời.');
        }
    });
});