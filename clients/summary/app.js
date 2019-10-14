
// when the form is submitted
$("#myform").submit(function (e) {

    e.preventDefault();

    var form = $(this);
    var url = form.attr('action');
    var param = form.serialize()

    alert("requesting " + param.slice(5))
    // send a get request with the above data
    $.ajax({
        type: "GET",
        url: url + '?url=' + param.slice(5),

        success: function (result) {
            // show response from the php script.
            $("#result").html("<strong>" + result + "</strong>")
        }
    });
});

