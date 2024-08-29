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
    currentPage = page
    jsonDataFilter = { ...jsonData }
    var extendedJsonData = { ...jsonData };
    extendedJsonData.page = page;
    $.ajax({
        url: '/received-inv/api/filter',
        method: 'POST',
        contentType: "application/json",
        data: JSON.stringify(extendedJsonData),
        success: function(response) {
            var $list = $(listSelector);
            var $pagination = $(paginationSelector);
            var data = response.data;
            var totalItems = response.totalItems; // Tổng số lượng phần tử từ backend

            // Xóa danh sách cũ
            $list.empty();

            // Render dữ liệu danh sách
            $.each(data, function(index, item) {
                var listItem = renderListItem(item)
                $list.append(listItem);
            });
        },
        error: function(xhr, status, error) {
            console.error("Error fetching hiring news: ", error);
        }
    });
}

function renderListItem(item) {
    var formatedDate = formatDate(item.HiringNews.date);
    var formatedPrice = formatPrice(item.HiringNews.price);
    var formatedRelativeDate = formatRelativeDate(item.invite_at);
    let backgroundColor = 'background-color: #FFFFFF;';
    let statusMessage = null;
    if(item.status === 1){
        backgroundColor = 'background-color: #99FF99;';
        statusMessage = "Đã chấp nhận"
    } else if(item.status === 2){
        backgroundColor = 'background-color: #FFCCFF;';
        statusMessage = "Đã từ chối"
    } else if(item.status === 5){
        backgroundColor = 'background-color: #FFFF99;';
        statusMessage = "Trùng thời gian"
    }
    let statusMessageHtml = statusMessage ? `<p class="list-group-item-text">${statusMessage}</p>` : '';
    const acceptInviteButtonHtml = item.status === 0
        ? `<a href="#" class="dropdown-item accept-invite-item" data-id="${item.invitation_id}">Chấp nhận</a>`
        : '';
    const refuseInviteButtonHtml = item.status === 0
        ? `<a href="#" class="dropdown-item refuse-invite-item" data-id="${item.invitation_id}">Từ chối</a>`
        : '';

    return `
        <li class="list-group-item" style="${backgroundColor}">
            <div class="list-group-item-figure">
                <a href="/info/profile/${item.chuloadai_id}" class="user-avatar">
                    <div class="avatar">
                        <img src="/assets/img/avatar/${item.HiringNews.User.image}" alt="..." class="avatar-img rounded-circle">
                    </div>
                </a>
            </div>
            <div class="list-group-item-body pl-3 pl-md-4">
                <div class="row">
                    <div class="col-12 col-lg-10">
                        <strong class="list-group-item-text break-word">
                            ${item.HiringNews.User.last_name} ${item.HiringNews.User.first_name}
                        </strong>
                        <br>
                        <strong class="list-group-item-text">${formatedDate}</strong>
                        <p class="list-group-item-text break-word">
                            ${item.HiringNews.Ward.full_name}, ${item.HiringNews.District.full_name}, ${item.HiringNews.Province.name}
                        </p>
                        <strong class="list-group-item-text">
                            ${formatedPrice}
                        </strong>
                    </div>
                    <div class="col-12 col-lg-2 text-lg-right">
                        <p class="list-group-item-text">Thời gian gửi: ${formatedRelativeDate} </p> 
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
                        <a class="dropdown-item view-detail" data-id="${item.HiringNews.hiring_news_id}">Xem chi tiết</a>
                        ${acceptInviteButtonHtml}   
                        ${refuseInviteButtonHtml}                     
                    </div>
                </div>
            </div>
        </li>
    `;
}

$(document).on('click', '.view-detail', function(event) {
    event.preventDefault();
    var hiringNewsId = $(this).data('id');
    $.ajax({
        url: '/news/api/detail/'+hiringNewsId,
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
            var hiringNews = response.data;
            var formattedCreatedAt = moment(hiringNews.created_at, "YYYY-MM-DD HH:mm:ss").format("DD/MM/YYYY HH:mm:ss");
            var formattedUpdatedAt = moment(hiringNews.updated_at, "YYYY-MM-DD HH:mm:ss").format("DD/MM/YYYY HH:mm:ss");
            var formattedDate = moment(hiringNews.date, "YYYY-MM-DD").format("DD/MM/YYYY");
            var formattedPrice = formatPrice(hiringNews.price); // Assuming you have a function to format price

            // Gán giá trị vào các trường trong modal
            // $('#detailNewsModal #name').text(hiringNews.User.last_name + " " + hiringNews.User.first_name);
            // $('#detailNewsModal #role').text(hiringNews.User.last_name);
            $('#detailNewsModal #link-facebook').attr('href', hiringNews.User.link_facebook);

            $('#detailNewsModal .modal-body .avatar-img').attr('src', '/assets/img/avatar/' + hiringNews.User.image);
            $('#detailNewsModal #detail-name').text(hiringNews.User.last_name + " " + hiringNews.User.first_name);
            $('#detailNewsModal #detail-role').text(hiringNews.User.role.name);
            $('#detailNewsModal #detail-phone').text("SDT: " + hiringNews.User.phone_number);
            $('#detailNewsModal #detail-email').text("Email: " + hiringNews.User.email);

            $('#detailNewsModal #detail-create-at').text(formattedCreatedAt);
            $('#detailNewsModal #detail-update-at').text(formattedUpdatedAt);
            $('#detailNewsModal #detail-date').text(formattedDate);
            $('#detailNewsModal #detail-price').text(formattedPrice);

            $('#detailNewsModal #detail-address').text(
                hiringNews.address_detail + ", " + hiringNews.Ward.full_name + ", " + hiringNews.District.full_name + ", " + hiringNews.Province.name
            );
            $('#detailNewsModal #detail-describe').text(hiringNews.describe);

            $('html').removeClass('topbar_open');
            $('#detailNewsModal').modal('show');

        }
    });
});

$(document).on('click', '.accept-invite-item', function(e) {
    e.preventDefault();
    let invitationId = $(this).data('id');

    $.ajax({
        url: '/received-inv/accept',
        type: 'POST',
        data: { id: invitationId },
        success: function(response) {
            renderListFilter(currentPage, itemsPerPage,  '#list-invitation', '#pagination-invitation', jsonDataFilter)
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
            alert('Đã xảy ra lỗi khi chấp nhận lời mời.');
        }
    });
});

$(document).on('click', '.refuse-invite-item', function(e) {
    e.preventDefault();
    let contractId = $(this).data('id');

    $.ajax({
        url: '/received-inv/refuse',
        type: 'POST',
        data: { id: contractId },
        success: function(response) {
            renderListFilter(currentPage, itemsPerPage,  '#list-invitation', '#pagination-invitation', jsonDataFilter)
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
            alert('Đã xảy ra lỗi khi từ chối lời mời.');
        }
    });
});