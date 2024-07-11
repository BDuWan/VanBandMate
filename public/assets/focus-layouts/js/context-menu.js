if ($(".table-atd-list").length > 0) {
  $.contextMenu({
    selector: ".table-atd-list tbody tr",
    build: function ($trigger) {
      var table = $trigger.closest("table");
      var actions = table.data("actions");
      var prefix = table.data("prefix");
      var name = table.data("name");
      var options = {
        callback: function (key, options) {
          // var m = "clicked: " + key;
          // window.console && console.log(m) || alert(m);
          // console.log(key, options);
          switch (key) {
            case "detail":
              showDetailPage($(this).attr("row-id"), prefix, name);
              break;
            case "goto":
              goDetailPage($(this).attr("row-id"), prefix, name);
              break;
            case "edit":
              goEditPage($(this).attr("row-id"), prefix, name);
              break;
            case "delete":
              doDestroy($(this).attr("row-id"), prefix, name);
              break;

            default:
          }
        },
        items: {
          // detail: { name: "Detail", icon: "fa-search" },
          // edit: { name: "Edit", icon: "fa-edit" },
          // sep1: "---------",
          // delete: { name: "Delete", icon: "fa-trash" },
          // sep2: "---------",
          // quit: { name: "Quit", icon: "fa-remove" },
        },
      };
      if (typeof actions["update"] !== "undefined" && actions["detailQuick"]) {
        options.items.detail = {
          name: "Detail",
          icon: "fa-search",
        };
      }
      options.items.goto = {
        name: "Go to",
        icon: "fa-arrow-right",
      };
      if (typeof actions["update"] !== "undefined" && actions["update"]) {
        options.items.edit = {
          name: "Edit",
          icon: "fa-edit",
        };
      }
      if (typeof actions["delete"] !== "undefined" && actions["delete"]) {
        options.items.sep2 = "---------";
        options.items.delete = {
          name: "Delete",
          icon: "fa-trash",
        };
      }

      options.items.sep3 = "---------";
      options.items.quit = {
        name: "Quit",
        icon: "fa-remove",
      };
      return options;
    },
  });
}

// if ($(".table-atd-link-1n").length > 0) {
$.contextMenu({
  selector: ".table-atd-link-1n tbody tr",
  build: function ($trigger) {
    var table = $trigger.closest("table");
    var actions = table.data("actions");
    var prefix = table.data("prefix");
    var name = table.data("name");
    var currentUrl = window.location.href;
    var currentUrlEncode = encodeURIComponent(currentUrl);
    var params = {};
    params["relatedFrom"] = currentUrlEncode;
    var options = {
      callback: function (key, options) {
        // var m = "clicked: " + key;
        // window.console && console.log(m) || alert(m);
        // console.log(key, options);
        switch (key) {
          case "detail":
            showDetailPage($(this).attr("row-id"), prefix, name);
            break;
          case "goto":
            goDetailPage($(this).attr("row-id"), prefix, name, params);
            break;
          case "edit":
            goEditPage($(this).attr("row-id"), prefix, name, params);
            break;
          case "delete":
            doDestroy($(this).attr("row-id"), prefix, name, params);
            break;

          default:
        }
      },
      items: {},
    };
    options.items.detail = {
      name: "Detail",
      icon: "fa-search",
    };
    options.items.goto = {
      name: "Go to",
      icon: "fa-arrow-right",
    };
    if (typeof actions["update"] !== "undefined" && actions["update"]) {
      options.items.edit = {
        name: "Edit",
        icon: "fa-edit",
      };
    }
    // if (typeof actions["delete"] !== "undefined" && actions["delete"]) {
    //   options.items.sep2 = "---------";
    //   options.items.delete = {
    //     name: "Delete",
    //     icon: "fa-trash",
    //   };
    // }

    options.items.sep3 = "---------";
    options.items.quit = {
      name: "Quit",
      icon: "fa-remove",
    };
    return options;
  },
});
// }

// if ($(".table-atd-link-n1").length > 0) {
$.contextMenu({
  selector: ".table-atd-link-n1 tbody tr",
  build: function ($trigger) {
    var table = $trigger.closest("table");
    var actions = table.data("actions");
    var prefix = table.data("prefix");
    var name = table.data("name");
    var currentUrl = window.location.href;
    var currentUrlEncode = encodeURIComponent(currentUrl);
    var params = {};
    params["relatedFrom"] = currentUrlEncode;
    var options = {
      callback: function (key, options) {
        // var m = "clicked: " + key;
        // window.console && console.log(m) || alert(m);
        // console.log(key, options);
        switch (key) {
          case "detail":
            showDetailPage($(this).attr("row-id"), prefix, name);
            break;
          case "goto":
            goDetailPage($(this).attr("row-id"), prefix, name, params);
            break;
          // case "edit":
          //   goEditPage($(this).attr("row-id"), prefix, name, params);
          //   break;
          // case "delete":
          //   doDestroy($(this).attr("row-id"), prefix, name, params);
          //   break;

          default:
        }
      },
      items: {},
    };
    options.items.detail = {
      name: "Detail",
      icon: "fa-search",
    };
    options.items.goto = {
      name: "Go to",
      icon: "fa-arrow-right",
    };
    // if (typeof actions["update"] !== "undefined" && actions["update"]) {
    //   options.items.edit = {
    //     name: "Edit",
    //     icon: "fa-edit",
    //   };
    // }
    // if (typeof actions["delete"] !== "undefined" && actions["delete"]) {
    //   options.items.sep2 = "---------";
    //   options.items.delete = {
    //     name: "Delete",
    //     icon: "fa-trash",
    //   };
    // }

    options.items.sep3 = "---------";
    options.items.quit = {
      name: "Quit",
      icon: "fa-remove",
    };
    return options;
  },
});
// }
