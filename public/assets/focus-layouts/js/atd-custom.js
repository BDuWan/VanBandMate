var scrollPosition;

$(document).on("pjax:send", function () {
  scrollPosition = $(window).scrollTop();
  $("#table-atd-loading").addClass("show");
});

$(document).on("pjax:success", function () {
  $(window).scrollTop(scrollPosition);
  $("#table-atd-loading").removeClass("show");
 
});
$(document).pjax("a.atd-list-tab, a.page-link", "#pjax-container", {
  // push: false,
});

$(document).on("click", ".btn-search .icon-search", function (e) {
  var form = $("form[pjax-container]");
  form.trigger("submit");
});

$(document).on("click", ".btn-search .icon-cancel", function (e) {
  $('input[name="s"]').val("");
  var form = $("form[pjax-container]");
  form.trigger("submit");
});

$(document).on("click", "#btn-reset", function (e) {
  // $('input[name="s"]').val("");
  var form = $("form[pjax-container]");
  form.find("input, select, textarea").each(function () {
    var elementType = this.nodeName.toLowerCase();
    $(this).val("");
    if (elementType == "select") {
      $(this).trigger("change");
    }
  });
  form.trigger("submit");
});

$(function () {
  var form = $("form[pjax-container]");
  var flagMoreFiltering = false;
  form.find("input, select, textarea").each(function () {
    var elementName = $(this).attr("name");
    var elementValue = $(this).val();
    if (elementValue) {
      if (elementName !== "s") {
        flagMoreFiltering = true;
      }
    }
  });

  if (flagMoreFiltering) {
    $("#btn-more-pulse").addClass("show");
  } else {
    $("#btn-more-pulse").removeClass("show");
  }
});

$(document).on("submit", "form[pjax-container]", function (event) {
  //$.pjax.submit(event, '#pjax-container')
  event.preventDefault(); // Prevent the form from submitting traditionally

  var form = $(this);
  // var query = form.find('input[name="s"]').val();

  // Get the current URL
  var url = new URL(window.location.href);

  // Remove all parameters except "tab" and "query"
  var params = url.searchParams;

  var keys = Array.from(params.keys());

  keys.forEach(function (key) {
    if (key !== "tab") {
      params.delete(key);
    }
  });

  // Append the updated query parameter
  // params.set('s', query);
  var flagMoreFiltering = false;
  form.find("input, select, textarea").each(function () {
    var elementValue = $(this).val();
    var elementName = $(this).attr("name");
    var elementValue = $(this).val();
    if (elementValue) {
      if (elementName !== "s") {
        flagMoreFiltering = true;
      }
      params.set(elementName, elementValue);
    }
  });

  if (flagMoreFiltering) {
    $("#btn-more-pulse").addClass("show");
  } else {
    $("#btn-more-pulse").removeClass("show");
  }

  // Perform the PJAX request
  $.pjax({
    url: url.href,
    container: "#pjax-container", // Replace with your container element ID
  });
});

$(document).on("click", function (event) {
  var $target = $(event.target);
  // console.log($target.closest(".select2-container--open"));
  if (
    !$target.closest(".btn-more-panel.show").length &&
    !$target.closest(".btn-search").length &&
    !$target.closest(".btn-more").length &&
    !$target.closest(".select2-container--open").length
  ) {
    // Click is outside of #myElement
    $("#btn-more a i").removeClass("active");
    $("#btn-more-panel").removeClass("show");
  }
});

$(".btn-app-right-sidebar-toggle").on("click", function (e) {
  $(".btn-app-right-sidebar-toggle i").toggleClass("active");
  $("#app-right-sidebar").toggleClass("show");
});

$(".btn-more .btn-press").on("click", function (e) {
  $("#btn-more a i").toggleClass("active");
  $("#btn-more-panel").toggleClass("show");
});

$(document).on("click", ".atd-list-tab", function () {
  $(".atd-list-tab").removeClass("active");
  $(this).addClass("active");
});
$(document).on("click", ".page-item", function () {
  $(".page-item").removeClass("active");
  $(this).addClass("active");
});
$(document).ready(function () {
  // does current browser support PJAX
  if ($.support.pjax) {
    $.pjax.defaults.timeout = 1000; // time in milliseconds
  }
});

(function ($) {
  $(".select2").select2();
  $(".select2-multiple").select2({
    closeOnSelect: false,
  });
  $(".type-select").select2({
    templateSelection: formatSelect,
    templateResult: formatSelect,
  });

  // $(".type-select-readonly").select2({
  //   templateSelection: formatSelect,
  //   templateResult: formatSelect,
  //   disabled: "readonly",
  // });

  $(".type-model").select2({
    templateSelection: formatModel,
    templateResult: formatModel,
  });
  function formatSelect(state) {
    if (!state.id) {
      return state.text;
    }
    var $state = $("<span><span></span></span>");

    // Use .text() instead of HTML string concatenation to avoid script injection issues
    $state.find("span").text(state.text).addClass(state.element.dataset.class);
    $state.find("img").attr("src", state.element.dataset.image);

    return $state;
  }

  function formatModel(state) {
    if (!state.id) {
      return state.text;
    }
    var $state = $(
      '<span class="d-flex align-items-center"><img width="24" height="24" class="type-model-img round-img mr-2" /><span></span></span>'
    );

    // Use .text() instead of HTML string concatenation to avoid script injection issues
    $state.find("span").text(state.text);
    $state.find("img").attr("src", state.element.dataset.image);

    return $state;
  }

  // $('.sidebar-right-trigger').on('click', function() {
  //     $('.sidebar-right').toggleClass('show');
  // });

  // $(".btn-create").attr("href", location.pathname + "/create");
})(jQuery);

function appShowFormDetail(id) {
  appLoadFormDetail(id);
  $(".form-modal-detail").modal({
    // backdrop: "static",
    // keyboard: false,
  });
}

function appLoadFormDetail(id) {
  var url = `${location.pathname}/${id}/api-show`;
  $.ajax({
    url: url,
    method: "GET",
    data: {
      // _token: "{{csrf_token()}}",
      _token: window.Laravel.csrfToken
    },
    dataType: "json",
    success: function (response) {
      $.fn.appSetContent(response);
    },
    error: function (xhr, status, error) {},
  });
}

$.fn.appSetContent = function (data) {
  var prefix = "d-";
  $.each(data, function (index, item) {
    var element = $(`#${prefix}${index}`);
    if (element.length) {
      element.html(item);
    }
  });
};

function appHideFormDetail() {
  $(".form-modal-detail").modal("hide");
}

function goEditPage(id, prefix, name, params = {}) {
  // window.location.href = `${location.pathname}/${id}/edit`;
  var paramString = "";
  var spliter = "?";
  $.each(params, function (index, item) {
    paramString += spliter + index + "=" + item;
    spliter = "&";
  });
  var href = `/${prefix}/${name}/${id}/edit`;
  if (paramString.length > 0) {
    href += paramString;
  }
  window.location.href = href;
}

function showDetailPage(id, prefix, name, params = {}) {
  params["page-no-layout"] = "true";
  var paramString = "";
  var spliter = "?";
  $.each(params, function (index, item) {
    paramString += spliter + index + "=" + item;
    spliter = "&";
  });
  var href = `/${prefix}/${name}/${id}`;
  if (paramString.length > 0) {
    href += paramString;
  }
  // Load data from page
  $(".modal-atd-loading").addClass("show");
  // $("#modal-detail-body-content").empty();
  $("#modal-detail-body-content").load(href, function (response, status, xhr) {
    if (status == "success") {
      // Code to execute on successful load
      $(".modal-atd-loading").removeClass("show");
    } else {
      // Code to execute on error
      $(".modal-atd-loading").removeClass("show");
    }
  });

  // appShowFormDetail(id);
  $("#modal-detail").modal({
    // backdrop: "static",
    // keyboard: false,
  });
}

function goDetailPage(id, prefix, name, params = {}) {
  // window.location.href = `${location.pathname}/${id}`;
  var paramString = "";
  var spliter = "?";
  $.each(params, function (index, item) {
    paramString += spliter + index + "=" + item;
    spliter = "&";
  });
  var href = `/${prefix}/${name}/${id}`;
  if (paramString.length > 0) {
    href += paramString;
  }
  window.location.href = href;
}

function doDestroy(id, prefix, name, params = {}) {
  // var paramString = params.join("&");
  swal({
    title: "Are you sure to delete ?",
    text: "You will not be able to recover this!!",
    type: "warning",
    showCancelButton: !0,
    confirmButtonColor: "#DD6B55",
    confirmButtonText: "Yes, delete it !!",
    closeOnConfirm: !1,
  }).then(function (result) {
    if (result.value) {
      submitDynamicForm(
        params,
        // addAllCurrentParams(`${location.pathname}/${id}`),
        addAllCurrentParams(`/${prefix}/${name}/${id}`),
        "DELETE"
      );
    }
  });
}

function submitDynamicForm(params, path, method = "POST") {
  // Create a form element
  var form = $("<form>");

  // Set form attributes and properties
  form.attr("action", path);
  form.attr("method", "POST");

  // Create input elements and set their attributes
  var input = $("<input>");
  input.attr("type", "hidden");
  input.attr("name", "_method");
  input.val(method);
  form.append(input);

  input = $("<input>");
  input.attr("type", "hidden");
  input.attr("name", "_token");
  input.val(window.Laravel.csrfToken);
  form.append(input);

  $.each(params, function (index, item) {
    input = $("<input>");
    input.attr("type", "hidden");
    input.attr("name", index);
    input.val(item);
    form.append(input);
  });

  // Append the form to the document body
  $("body").append(form);

  // Submit the form programmatically
  form.submit();
}

function addAllCurrentParams(url) {
  // Get all query parameters
  var queryParams = new URLSearchParams(window.location.search);

  // Get all parameters as an object
  var params = {};
  for (var pair of queryParams.entries()) {
    params[pair[0]] = pair[1];
  }

  // Add additional parameters if needed
  // params.newParam = 'example';
  return url + "?" + $.param(params);
}

$(document).on("dblclick", ".table-atd tbody tr", function () {
  var table = $(this).closest("table");
  var prefix = table.data("prefix");
  var name = table.data("name");
  showDetailPage($(this).attr("row-id"), prefix, name);
  // appShowFormDetail($(this).attr("row-id"));
});

$(".datepicker").datepicker({
  dateFormat: "dd/mm/yy",
  uiLibrary: "bootstrap4",
});

$(".input-type-date").datepicker({
  dateFormat: "dd/mm/yy",
  uiLibrary: "bootstrap4",
});

$(".input-type-datetime").datepicker({
  dateFormat: "dd/mm/yy",
  uiLibrary: "bootstrap4",
});

$(".clockpicker").clockpicker({
  placement: 'left',
  minutestep: 5,
  donetext: "Done",
});
