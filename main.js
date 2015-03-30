rendering = false;
ITERATIONS = 10;
MUTATIONS_PER_ITERATION = 400;
$("#start").on("click", function() {
  if (!rendering) {
    rendering = true;
    var targetPrefix = "/remote/img?url=";
    var approxPrefix = "/approx/img?url=";
    var imgUrl = $("#imgUrl").val();
    $("#target-image").attr("src", targetPrefix + imgUrl).css("height", "auto").css("width", "auto");
    $("#approx-image").attr("src", approxPrefix + imgUrl).css("height", "auto").css("width", "auto");
    console.log("Rendering", imgUrl);
    $.get("/stats", function(data) {
        console.log(data);
    });
  }
});