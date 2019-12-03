
// when the form is submitted
$("#myform").submit(function (e) {

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
            // $("#result").html("<strong>" + result + "</strong>")
        }
    });
});

$('#creds').submit(function (e) {
    e.preventDefault();

    var form = $(this);
    var url = form.attr('action');
    var param = form.serialize()

    var b = {
        "Email":        $('#email').val(),
        "Password":     $('#pass').val(),
        "PasswordConf": $('#pass').val(),
        "UserName":     "user",
        "FirstName":    "First",
        "LastName":     "Name"
    }

    //console.log(b)

    // send a get request with the above data
    $.ajax({
        type: "POST",
        url: url,
        data: JSON.stringify(b),
        contentType:"application/json",
        success: function (result) {
            console.log(result)
            // $("#result").html("<strong>" + result + "</strong>")
        }
    });
})

$('#credsSess').submit(function (e) {
    e.preventDefault();

    var form = $(this);
    var url = form.attr('action');
    var param = form.serialize()

    var b = {
        "Email":        $('#emailSess').val(),
        "Password":     $('#passSess').val(),
    }

    //console.log(b)

    // send a get request with the above data
    $.ajax({
        type: "POST",
        url: url,
        data: JSON.stringify(b),
        contentType:"application/json"
        // success: function (result) {
        //     console.log(result)
        //     console.log(result.getResponseHeader("authorization"))
        //     // $("#result").html("<strong>" + result + "</strong>")
        // }
    }).done(function (data, textStatus, xhr) { 
        console.log(data)
        console.log(xhr.getResponseHeader('authorization')); 
    });
})

$('#signOut').submit(function (e) {
    e.preventDefault();

    var form = $(this);
    var url = form.attr('action');
    var param = form.serialize()

    //console.log(b)

    // send a get request with the above data
    $.ajax({
        type: "DELETE",
        url: url,
        success: function (result) {
            console.log(result)
            // $("#result").html("<strong>" + result + "</strong>")
        }
    });
})

// $('#credsId').submit(function (e) {
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

