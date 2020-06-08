import * as auth from "./js/auth.handlers"

var app = new Vue({
    el: '#app',
    data: {
        isAuthenticated: false
    },
    methods: {
        Authenticate: auth.signIn(),
        CreateUser: auth.createNewUser(),
        SignOut: auth.signOut()
    }
})

$("#get-messages").on("click", async function getMessages() {
    var url = "https://slack.api.tristanmacelli.com/v1/channels/5edbf9c61d5a4d4adcee0f3b"
    sessionToken = localStorage.getItem('auth')

    // send a get request with the above data
    let resp = await fetch(url, {
        method: 'GET',
        headers: {
            "Authorization": sessionToken
        }
    })
    if (!resp.ok) {
        alert("Error: ", resp.status)
    }
    response = await resp.json()
    console.log(response)
})

$("#get-channels").on("click", async function getChannels() {
    console.log("Getting Channels:")
    var url = "https://slack.api.tristanmacelli.com/v1/channels"
    sessionToken = localStorage.getItem('auth')

    // send a get request with the above data
    resp = await fetch(url, {
        method: 'GET',
        headers: {
            "Authorization": sessionToken
        }
    });
    if (!resp.ok) {
        alert(resp.status)
        throw new Error(resp.status)
    }
    let channels = await resp.json();
    console.log(channels)
})

// when the form is submitted
$("#summary").submit(async function getSummary(e) {
    e.preventDefault();

    let urlToSummarize = document.getElementById("summaryUrl").innerText
    let url = new URL("https://slack.api.tristanmacelli.com/v1/summary")
    url.searchParams.append("url", urlToSummarize)

    // send a get request with the above data
    let resp = await fetch(url, {
        method: 'GET'
    })
    if (!resp.ok) {
        alert("Error: ", resp.status)
    }
    response = await resp.json()
    console.log(response)
    // display_results(result, "#summaryResult")
});

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

$("#home-page").ready(homePageLoad)

// Loads user info into home page after page is loaded
function homePageLoad() {
    sessionToken = localStorage.getItem('auth')
    if (sessionToken) {
        auth.request_user(auth.display_user_first_name, sessionToken)
        new WebSocket("wss://slack.api.tristanmacelli.com/v1/ws?auth=" + sessionToken)
    }
}

$("#account-page").ready(accountPageLoad)

function accountPageLoad() {
    sessionToken = localStorage.getItem('auth')
    if (sessionToken) {
        auth.request_user(auth.display_user, sessionToken)
    }
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
