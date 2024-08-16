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
                renderList(page, itemsPerPage, listSelector, paginationSelector);
                setupPaginationControls(page);
            });
    }

    function renderEllipsis() {
        $('<div class="page-item"><strong>...</strong></div>').appendTo($pagination);
    }

    renderPageItem(1);
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
    renderList(1, itemsPerPage, listSelector, paginationSelector);
    setupPaginationControls(1);
}

function renderList(page, itemsPerPage, listSelector, paginationSelector) {
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
            var totalItems = response.totalItems; // Tổng số lượng phần tử từ backend

            // Xóa danh sách cũ
            $list.empty();

            // Render dữ liệu danh sách
            $.each(data, function(index, item) {
                var formatedDate = formatDate(item.date);
                var formatedPrice = formatPrice(item.price);
                var formatedRelativeDate = formatRelativeDate(item.created_at)
                const backgroundColor = item.hiring_enough == 0 ? 'background-color: #FFFFFF;' : 'background-color: #99FFFF;';
                var listItem = `
                    <li class="list-group-item" style="${backgroundColor}">
                        <div class="list-group-item-figure">
                            <a href="conversations.html" class="user-avatar">
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
                                    <a href="#" class="dropdown-item">Xem chi tiết</a>
                                    <a href="#" class="dropdown-item">Chỉnh sửa</a>
                                    <a href="#" class="dropdown-item">Xem danh sách ứng tuyển</a>
                                    <a href="#" class="dropdown-item">Xóa</a>
                                </div>
                            </div>
                        </div>
                    </li>
                `;
                $list.append(listItem);
            });
        },
        error: function(xhr, status, error) {
            console.error("Error fetching hiring news: ", error);
        }
    });
}