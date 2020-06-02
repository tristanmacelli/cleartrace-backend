
// when the form is submitted
$("#summary").submit(function getSummary(e) {

    e.preventDefault();

    var form = $(this);
    var url = form.attr('action');
    var param = form.serialize()

    // send a get request with the above data
    $.ajax({
        type: "GET",
        url: url + '?url=' + param.slice(5),

        success: function (result) {
            display_results(result)
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
        body: b,
        contentType: "application/json",
        dataType: "json",
        crossDomain: true,
        success: function (result) {
            console.log(result)
        }
    }).done(function (_, _, xhr) {
        console.log(xhr.getResponseHeader('authorization'));
        // https://developer.mozilla.org/en-US/docs/Web/API/Window/sessionStorage
        sessionStorage.setItem('auth', xhr.getResponseHeader('authorization'))
        // window.location.replace("https://slack.client.tristanmacelli.com/home.html");
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
            // $("#result").html("<strong>" + result + "</strong>")
        }
    }).done(function (_, _, xhr) {
        console.log(xhr.getResponseHeader('authorization'));
        sessionStorage.setItem('auth', xhr.getResponseHeader('authorization'))
        window.location.replace("https://slack.client.tristanmacelli.com/home.html");
    });
})

// Removing a session based on the form values
$('#signOut').submit(function signOut(e) {
    console.log(sessionStorage.getItem('auth'))
    e.preventDefault();

    var form = $(this);
    var url = form.attr('action');

    // send a get request with the above data
    $.ajax({
        type: "DELETE",
        url: url,
        headers: {
            Authorization: sessionStorage.getItem('auth'),
        },        
        success: function (result) {
            console.log(result)
        }
    }).done(function () {
        window.location.replace("https://slack.client.tristanmacelli.com/index.html");
    });
})

// $('#getUser').submit(function (e) {
//     e.preventDefault();

//     var form = $(this);
//     var url = form.attr('action') + $('#userId').val();
//     var param = form.serialize()

//     // var b = {
//     //     "Email":        $('#emailSess').val(),
//     //     "Password":     $('#passSess').val(),
//     // }

//     //console.log(b)

//     // send a get request with the above data
//     $.ajax({
//         type: "GET",
//         url: url,
//         contentType:"application/json",
//         success: function (result) {
//             console.log(result)
//             // $("#result").html("<strong>" + result + "</strong>")
//         }
//     });
// })

function display_results(result) {
    json_obj = JSON.parse(result)
    var final_html = "<h2> Summary of " + json_obj.url + "</h2>"
    final_html += "<h4> Title : " + json_obj.title + " </h4>"
    final_html += "<h4> Description : " + json_obj.description + " </h4>"
    image_div = "<div>"

    for (i = 0; i < json_obj.images.length; i++) {
        image_div += "<img src=\"" + json_obj.images[i].url + "\">"
    }
    image_div += "</div>"
    final_html += image_div
    $("#result").html(final_html)
}

