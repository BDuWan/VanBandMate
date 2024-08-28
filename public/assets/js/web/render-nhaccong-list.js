
// function setupPagination(totalItems, itemsPerPage, listSelector, paginationSelector) {
//     var $pagination = $(paginationSelector);
//
//     $pagination.empty();
//     var numPages = Math.ceil(totalItems / itemsPerPage);
//
//     function renderPageItem(page) {
//         return $('<div class="page-item"><a href="#" class="page-link"></a></div>')
//             .find('a')
//             .text(page)
//             .end()
//             .appendTo($pagination)
//             .on('click', function(e) {
//                 e.preventDefault();
//                 renderNewsList(page, itemsPerPage, listSelector, paginationSelector);
//                 setupPaginationControls(page);
//             });
//     }
//
//     function renderEllipsis() {
//         $('<div class="page-item"><strong>...</strong></div>').appendTo($pagination);
//     }
//
//     function setupPaginationControls(currentPage) {
//         $pagination.find('a').removeClass('active');
//         $pagination.empty();
//
//         renderPageItem(1);
//
//         if (currentPage > 3) {
//             renderEllipsis();
//         }
//
//         var start = Math.max(2, currentPage - 1);
//         var end = Math.min(numPages - 1, currentPage + 1);
//
//         for (var i = start; i <= end; i++) {
//             renderPageItem(i);
//         }
//
//         if (currentPage < numPages - 2) {
//             renderEllipsis();
//         }
//
//         if (numPages > 1) {
//             renderPageItem(numPages);
//         }
//
//         $pagination.find('a').filter(function() {
//             return $(this).text() == currentPage;
//         }).addClass('active');
//     }
//
//     // Hiển thị trang đầu tiên
//     renderNewsList(1, itemsPerPage, listSelector, paginationSelector);
//     setupPaginationControls(1);
// }
//
// function renderNewsList(page, itemsPerPage, listSelector, paginationSelector) {
//     currentPage = page
//     $.ajax({
//         url: '/news/api',
//         type: 'GET',
//         dataType: 'json',
//         data: {
//             page: page,
//             itemsPerPage: itemsPerPage
//         },
//         success: function(response) {
//             var $list = $(listSelector);
//             var $pagination = $(paginationSelector);
//             var data = response.data;
//             var user_id = response.user_id;
//             var totalItems = response.totalItems; // Tổng số lượng phần tử từ backend
//
//             // Xóa danh sách cũ
//             $list.empty();
//
//             // Render dữ liệu danh sách
//             $.each(data, function(index, item) {
//                 var listItem = renderListItem(item, user_id)
//                 $list.append(listItem);
//             });
//         },
//         error: function(xhr, status, error) {
//             console.error("Error fetching hiring news: ", error);
//         }
//     });
// }

var currentPage = 1
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
        url: '/hiring/api/find',
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
    return `
        <li class="list-group-item">
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
                                  
                    </div>
                </div>
            </div>
        </li>
    `;
}

$(document).on('click', '.apply-item', function(e) {
    e.preventDefault();
    let contractId = $(this).data('id');

    $.ajax({
        url: '/news/apply',
        type: 'POST',
        data: { id: contractId },
        success: function(response) {
            renderNewsList(currentPage, itemsPerPage,  '#list-news', '#pagination-news')
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
            alert('Đã xảy ra lỗi khi gửi yêu cầu ứng tuyển.');
        }
    });
});

$(document).on('click', '.cancel-item', function(e) {
    e.preventDefault();
    let contractId = $(this).data('id');

    $.ajax({
        url: '/news/cancel-apply',
        type: 'POST',
        data: { id: contractId },
        success: function(response) {
            renderNewsList(currentPage, itemsPerPage,  '#list-news', '#pagination-news')
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
            alert('Đã xảy ra lỗi khi gửi yêu cầu ứng tuyển.');
        }
    });
});