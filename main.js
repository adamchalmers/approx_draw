$(function() {
    const socket = new WebSocket(
        "ws://" +
            location.hostname +
            (location.port ? ':' + location.port : '') +
            "/status/"
    );

    socket.onmessage = function(evt) {
        const event = JSON.parse(evt.data);

        switch (event.type) {
        case "log":
            const log = $("#serverlog");
            log.append("<li>" + event.properties.text + "</li>");
            log.scrollTop(log[0].scrollHeight);
            break;
        case "mutation":
            const progress = parseFloat(event.properties.progress) * 100.00 / parseFloat(event.properties.total - 1);
            $("#progress .filled").css("width", progress.toString() + "%");
            break;
        }
    };

    socket.onclose = function(evt) {
        console.error(evt);
    }
});

i = 0; 
$("#start").on("click", function() {
    // Load the target and approximation images.
    const targetPrefix = "/remote/img?url=";
    const approxPrefix = "/approx/img?url=";
    const imgUrl = $("#imgUrl").val();
    $("#target-image").attr("src", targetPrefix + imgUrl).css("height", "auto").css("width", "auto");
    $("#approx-image").attr("src", approxPrefix + imgUrl + "&" + i++).css("height", "auto").css("width", "auto");
    console.log("Rendering", imgUrl);
    const startTime = (new Date()).getTime();


    // Once the target has loaded, we can display it
    $("#target-image").on("load", function() {
        $("#error").text("");
        $(this).show();
        $("#approx-image").show();
    })


    // Once the approximation's loaded, hide the loading placeholder and log stats to console.
    $("#approx-image").on("load", function() {
        $.get("/stats", function(data) {
            console.log(data);
            console.log("Took " + Math.round(((new Date()).getTime()-startTime)/1000) + " seconds.");
        });
    });
});
