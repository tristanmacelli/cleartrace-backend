// Creating a new session based on the form values
$('#getChannels').submit(function signIn(e) {
    e.preventDefault();

    var form = $(this);
    var url = form.attr('action');
    var method = form.attr('method')
    var param = form.serialize()

    // send a GET request with the above data
    $.ajax({
        type: method, // GET
        url: url, // endpoint: v1/channels
        contentType: "application/json"
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

// Creating a new session based on the form values
$('#createChannel').submit(function signIn(e) {
    e.preventDefault();

    var form = $(this);
    var url = form.attr('action');
    var method = form.attr('method')
    var param = form.serialize()

    var b = {
        "name": $('#name').val(),
        "description": $('#description').val(),
        "private": $('#private').val(),
        // "members": "ENTER THE CURRENT USER'S ID HERE",
        // "createdAt": "ENTER A RANDOM DATE HERE",
        // "creator": "ENTER THE CURRENT USER'S ID HERE",
        // "editedAt": "ENTER A RANDOM DATE HERE",
    }

    // send a POST request with the above data
    $.ajax({
        type: method, // POST
        url: url, // endpoint: v1/channels
        contentType: "application/json",
        dataType: "json",
        body: b,
        crossDomain: true,
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
