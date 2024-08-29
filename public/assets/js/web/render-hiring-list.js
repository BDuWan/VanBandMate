var currentPage = 1
var jsonDataFilter = {}
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
                renderHiringList(page, itemsPerPage, listSelector, paginationSelector);
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
    renderHiringList(1, itemsPerPage, listSelector, paginationSelector);
    setupPaginationControls(1);
}

function renderHiringList(page, itemsPerPage, listSelector, paginationSelector) {
    currentPage = page
    $.ajax({
        url: '/hiring/api',
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
            console.error("Error fetching hiring news: ", error);
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
    jsonDataFilter = { ...jsonData }
    var extendedJsonData = { ...jsonData };
    extendedJsonData.page = page;
    $.ajax({
        url: '/hiring/api/filter',
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
            console.error("Error fetching hiring news: ", error);
        }
    });
}

function renderListItem(item, user_id) {
    var formatedDate = formatDate(item.date);
    var formatedPrice = formatPrice(item.price);
    var formatedRelativeDate = formatRelativeDate(item.created_at);
    const backgroundColor = item.hiring_enough == 0 ? 'background-color: #FFFFFF;' : 'background-color: #99FFFF;';
    let enoughMessageHtml = ''
    if(item.hiring_enough == 1){
        enoughMessageHtml = `<p class="list-group-item-text">Đã tuyển đủ</p>`
    }
    const deleteButtonHtml = user_id === item.chuloadai_id
        ? `<a href="#" class="dropdown-item delete-item" data-id="${item.hiring_news_id}">Xóa</a>`
        : '';
    const editButtonHtml = user_id === item.chuloadai_id
        ? `<a href="#" class="dropdown-item edit-item" data-id="${item.hiring_news_id}">Chỉnh sửa</a>`
        : '';
    const showButtonHtml = user_id === item.chuloadai_id
        ? `<a href="#" class="dropdown-item show-item" data-id="${item.hiring_news_id}">Xem danh sách ứng tuyển</a>`
        : '';
    const inviteButtonHtml = user_id === item.chuloadai_id
        ? `<a href="/hiring/invite/${item.hiring_news_id}" class="dropdown-item">Tìm kiếm nhạc công</a>`
        : '';

    return `
        <li class="list-group-item" style="${backgroundColor}">
            <div class="list-group-item-figure">
                <a href="/info/profile/${item.chuloadai_id}" class="user-avatar">
                    <div class="avatar">
                        <img src="/assets/img/avatar/${item.User.image}" alt="..." class="avatar-img rounded-circle">
                    </div>
                </a>
            </div>
            <div class="list-group-item-body pl-3 pl-md-4">
                <div class="row">
                    <div class="col-12 col-lg-10">
                        <strong class="list-group-item-text break-word">${item.User.last_name} ${item.User.first_name}</strong>
                        <br>
                        <strong class="list-group-item-text">${formatedDate}</strong>
                        <p class="list-group-item-text break-word">
                            ${item.Ward.full_name}, ${item.District.full_name}, ${item.Province.name}
                        </p>
                        <strong class="list-group-item-text">
                            ${formatedPrice}
                        </strong>
                    </div>
                    <div class="col-12 col-lg-2 text-lg-right">
                        <p class="list-group-item-text">Thời gian tạo: ${formatedRelativeDate} </p>
                        <strong>${enoughMessageHtml}</strong>
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
                        <a class="dropdown-item view-detail" data-id="${item.hiring_news_id}">Xem chi tiết</a>
                        ${showButtonHtml}
                        ${inviteButtonHtml}
                        ${editButtonHtml}                       
                        ${deleteButtonHtml}
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
        url: '/hiring/api/detail/'+hiringNewsId,
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
            // $('#detailHiringNewsModal #name').text(hiringNews.User.last_name + " " + hiringNews.User.first_name);
            // $('#detailHiringNewsModal #role').text(hiringNews.User.last_name);
            $('#detailHiringNewsModal #link-facebook').attr('href', hiringNews.User.link_facebook);

            $('#detailHiringNewsModal .modal-body .avatar-img').attr('src', '/assets/img/avatar/' + hiringNews.User.image);
            $('#detailHiringNewsModal #detail-name').text(hiringNews.User.last_name + " " + hiringNews.User.first_name);
            $('#detailHiringNewsModal #detail-role').text(hiringNews.User.role.name);
            $('#detailHiringNewsModal #detail-phone').text("SDT: " + hiringNews.User.phone_number);
            $('#detailHiringNewsModal #detail-email').text("Email: " + hiringNews.User.email);

            $('#detailHiringNewsModal #detail-create-at').text(formattedCreatedAt);
            $('#detailHiringNewsModal #detail-update-at').text(formattedUpdatedAt);
            $('#detailHiringNewsModal #detail-date').text(formattedDate);
            $('#detailHiringNewsModal #detail-price').text(formattedPrice);

            $('#detailHiringNewsModal #detail-address').text(
                hiringNews.address_detail + ", " + hiringNews.Ward.full_name + ", " + hiringNews.District.full_name + ", " + hiringNews.Province.name
            );
            $('#detailHiringNewsModal #detail-describe').text(hiringNews.describe);

            $('html').removeClass('topbar_open');
            $('#detailHiringNewsModal').modal('show');

        }
    });
});

$(document).on('click', '.show-item', function(event) {
    event.preventDefault();
    var hiringNewsId = $(this).data('id');
    $.ajax({
        url: '/hiring/api/list-apply/' + hiringNewsId,
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
            var listApply = response.data
            var hiringEnough = response.hiring_enough
            if(hiringEnough){
                $('#showListApplyModal #hiringEnoughCheckbox').prop('checked', true);
            }
            $.each(listApply, function(index, item) {
                var listItem = renderListApplyItem(item)
                $('#showListApplyModal #list-apply').append(listItem);
            });
            $('#saveApplyBtn').data('id', hiringNewsId);
            $('html').removeClass('topbar_open');
            $('#showListApplyModal').modal('show');
            $('#showListApplyModal').on('hidden.bs.modal', function () {
                $('#showListApplyModal #list-apply').empty();
            });
        }
    });
});

$(document).on('click', '.delete-item', function(event) {
    event.preventDefault();
    var hiringNewsId = $(this).data('id');
    alert(hiringNewsId)
});

function renderListApplyItem(item) {
    var formatedApplyAt = formatRelativeDate(item.apply_at);
    const backgroundColor = item.status == 0 ? 'background-color: #FFFFFF;' : 'background-color: #99FFFF;';
    const checkboxHtml = item.status == 0
        ? `<input type="checkbox" class="approve-checkbox" style="width: 25px; height: 25px" data-id="${item.user_hiring_news_id}" data-nhaccong-id="${item.nhaccong_id}" />`
        : '';
    return `
        <li class="list-group-item" style="${backgroundColor}">
            <div class="list-group-item-figure">
                <a href="/info/profile/${item.nhaccong_id}" class="user-avatar">
                    <div class="avatar">
                        <img src="/assets/img/avatar/${item.User.image}" alt="..." class="avatar-img rounded-circle">
                    </div>
                </a>
            </div>
            <div class="list-group-item-body pl-3 pl-md-4">
                <div class="row">
                    <div class="col-12 col-lg-10">
                        <strong class="list-group-item-text break-word">${item.User.last_name} ${item.User.first_name}</strong>            
                        <p class="list-group-item-text">${formatedApplyAt} </p>         
                    </div>
                 
                </div>
            </div>
            <div class="list-group-item-figure">
                ${checkboxHtml}
            </div>
        </li>
    `;
}

$(document).on('click', '#saveApplyBtn', function(event) {
    event.preventDefault();
    var hiringNewsId = $(this).data('id');

    var hiringEnough = $('#showListApplyModal #hiringEnoughCheckbox').is(':checked')
    var selectedItems = [];
    $('.approve-checkbox').each(function() {
        var isChecked = $(this).is(':checked');
        if(isChecked) {
            var userHiringNewsId = $(this).data('id');
            var nhaccongId = $(this).data('nhaccong-id');

            selectedItems.push({
                user_hiring_news_id: userHiringNewsId,
                nhaccong_id: nhaccongId
            });
        }
    });

    var formData = {
        hiring_news_id: hiringNewsId,
        hiring_enough: hiringEnough,
        selected_items: selectedItems
    };

    // Gửi dữ liệu qua AJAX
    $.ajax({
        url: '/hiring/api/save-apply',
        method: 'POST',
        contentType: 'application/json',
        data: JSON.stringify(formData),
        success: function(response) {
            if(response !== "success"){
                swal("", response, {
                    icon : "error",
                    buttons: {
                        confirm: {
                            className : 'btn btn-danger'
                        }
                    },
                });
                return false;
            }
            swal("", "Thành công", {
                icon: "success",
                buttons: {
                    confirm: {
                        className: 'btn btn-danger'
                    }
                }
            }).then(function() {
                $('#showListApplyModal').modal('hide');
                renderListFilter(currentPage, itemsPerPage,  '#list-hiring', '#pagination-hiring', jsonDataFilter)
            });
        },
        error: function() {
            swal("Error", "Có lỗi xảy ra khi lưu dữ liệu.", "error");
        }
    });
});

