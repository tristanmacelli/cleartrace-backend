
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

function display_results(result) {
    if (typeof result == 'object') {
        json_obj = result
    }else {
        json_obj = JSON.parse(result)
    }
    var final_html = "<h2> Summary of " + json_obj.url + "</h2>"
    final_html += "<h4> Title : " + json_obj.title + " </h4>"
    final_html += "<h4> Description : " + json_obj.description + " </h4>"
    image_div = "<div>"
    if (typeof json_obj.images != "undefined"){
        for (i = 0; i < json_obj.images.length; i++) {
            image_div += "<img src=\"" + json_obj.images[i].url + "\">"
        }
    }
    
    image_div += "</div>"
    final_html += image_div
    $("#result").html(final_html)

}

