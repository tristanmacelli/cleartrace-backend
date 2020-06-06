
// when the form is submitted
$("#summary").submit(function getSummary(e) {

    e.preventDefault();

    var form = $(this);
    var url = form.attr('action');

    // send a get request with the above data
    $.ajax({
        type: "GET",
        url: url,
        data: {
            url: $('#summaryUrl').val()
        },
        crossDomain: false,
        success: function (result) {
            display_results(result, "#summaryResult")
        },
        error: function(result) {
            alert(result)
        }
    });
});

// Creating a new user based on the form values
$('#createUser').submit(function createNewUser(e) {
    e.preventDefault();

    var form = $(this);
    var url = form.attr('action');
    var username = $('#firstname').val() + "." + $('#lastname').val()

    var b = {
        "Email": $('#email').val(),
        "Password": $('#pass').val(),
        "PasswordConf": $('#pass').val(),
        "UserName": username,
        "FirstName": $('#firstname').val(),
        "LastName": $('#lastname').val()
    }

    // send a get request with the above data
    $.ajax({
        type: "POST",
        url: url,
        data: JSON.stringify(b),
        contentType: "application/json",
        crossDomain: true,
        success: function (result) {
            console.log(result)
        },
        error: function(result) {
            // $("#result").html("<strong>" + result + "</strong>")
            message = "Error: " + result.responseText
            alert(message)
        }
    }).done(function (_, _, xhr) {
        // https://developer.mozilla.org/en-US/docs/Web/API/Window/sessionStorage
        sessionToken = xhr.getResponseHeader('authorization')
        localStorage.setItem('auth', sessionToken)
        window.location.replace("https://slack.client.tristanmacelli.com/home.html");
    });
})

// Creating a new session based on the form values
$('#signIn').submit(function signIn(e) {
    e.preventDefault();

    var form = $(this);
    var url = form.attr('action');

    var b = {
        "Email": $('#emailSess').val(),
        "Password": $('#passSess').val(),
    }

    // send a get request with the above data
    $.ajax({
        type: "POST",
        url: url,
        data: JSON.stringify(b),
        contentType: "application/json",
        success: function (result) {
            console.log(result)
        },
        error: function(result) {
            // $("#result").html("<strong>" + result + "</strong>")
            message = "Error: " + result.responseText
            alert(message)
        }
    }).done(function (_, _, xhr) {
        sessionToken = xhr.getResponseHeader('authorization')
        localStorage.setItem('auth', sessionToken)
        window.location.replace("https://slack.client.tristanmacelli.com/home.html");
    });
})

// Removing a session based on the form values
$('#signOut').submit(function signOut(e) {
    e.preventDefault();

    var form = $(this);
    var url = form.attr('action');
    sessionToken = localStorage.getItem('auth')
    localStorage.removeItem('auth')

    // send a get request with the above data
    $.ajax({
        type: "DELETE",
        url: url,
        headers: {
            Authorization: sessionToken,
        },        
        success: function (result) {
            console.log(result)
            localStorage.removeItem('auth')
            window.location.replace("https://slack.client.tristanmacelli.com");
        }
    })
})

// $('#updateUser').submit()

function display_results(result, result_id) {
    if (result) {
        result_str = JSON.stringify(result)
        json_obj = JSON.parse(result_str)

        var final_html = "<h2> Summary of " + json_obj.url + "</h2>"
        final_html += "<h4> Title : " + json_obj.title + " </h4>"
        final_html += "<h4> Description : " + json_obj.description + " </h4>"
        image_div = "<div>"

        if (json_obj.images) {
            for (i = 0; i < json_obj.images.length; i++) {
                image_div += "<img src=\"" + json_obj.images[i].url + "\">"
            }
        }
        
        image_div += "</div>"
        final_html += image_div
        $(result_id).html(final_html)
    }
}

// Loads user info into home page after page is loaded
$(document).ready(
    function homePageLoad() {
        if (window.location.toString().includes("home")) {
            sessionToken = localStorage.getItem('auth')
            if (sessionToken) {
                var url = "https://slack.api.tristanmacelli.com/v1/users/"
            
                // send a get request with the above data
                $.ajax({
                    type: "GET",
                    url: url,
                    contentType:"application/json",
                    headers: {
                        Authorization: sessionToken,
                    },
                    success: function (result) {
                        console.log("Result", result)

                        display_user(result)
                    },
                    error: function (result) {
                        console.log(result)
                    }
                });
            }
        }
    }
);

function display_user(result) {
    var final_html = "<table cellspacing=0 role=\"presentation\"><tbody>"
    final_html += "<tr><td><p>Username: </p></td>"
    final_html += "<td>" + result.UserName + "</td></tr>"
    
    final_html += "<tr><td><p>First name: </p></td>"
    final_html += "<td>" + result.FirstName + "</td>"
    final_html += "<td><p>Last name: </p></td>"
    final_html += "<td>" + result.LastName + "</td></tr>"

    $("#userInfo").html(final_html)
}

// Users returning to the website with an active session get redirected from the log in page
// to the home page
$(document).ready(
    function indexPageLoad() {
        if (window.location.pathname == "/" && localStorage.getItem("auth")) {
            window.location.replace("https://slack.client.tristanmacelli.com/home.html");
        }
    }
)

// Unauthenticated users trying to visit pages requiring authentication will be returned to
// the log in page
$(document).ready(
    function requiredAuthPageLoad() {
        if (window.location.pathname != "/" && !localStorage.getItem("auth")) {
            alert("You are not authenticated: please return to the Log In page")
            window.location.replace("https://slack.client.tristanmacelli.com/");
        }
    }
)
