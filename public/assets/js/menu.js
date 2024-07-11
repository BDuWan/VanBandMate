function changeBackgroundColor() {
    var currentURL = window.location.href;

    if (currentURL.includes("/home")) {
        document.getElementById("Home").style.backgroundColor = "#fff";
    }

    if (currentURL.includes("/dashboard")) {
        document.getElementById("Dashboard").style.backgroundColor = "#fff";
    }

    if (currentURL.includes("/courses")) {
        document.getElementById("Courses").style.backgroundColor = "#fff";
    }

    if (currentURL.includes("/classes")) {
        document.getElementById("Classes").style.backgroundColor = "#fff";
    }

    if (currentURL.includes("/managements")) {
        document.getElementById("Managements").style.backgroundColor = "#fff";
    }

    if (currentURL.includes("/accounts")) {
        console.log(currentURL);
        document.getElementById("Accounts").style.backgroundColor = "#fff";
    }

}
window.onload = function () {
    changeBackgroundColor();
};