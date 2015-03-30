rendering = false;
MAXSIZE = 300;
$("#start").on("click", function() {
  if (!rendering) {
    rendering = true;

    // Load the target and approximation images.
    var targetPrefix = "/remote/img?url=";
    var approxPrefix = "/approx/img?url=";
    var imgUrl = $("#imgUrl").val();
    $("#target-image").attr("src", targetPrefix + imgUrl).css("height", "auto").css("width", "auto");
    $("#approx-image").attr("src", approxPrefix + imgUrl).css("height", "auto").css("width", "auto");

    // Once the target has loaded, we can check it's within the size limits.
    $("#target-image").on("load", function() {

        if (this.width > MAXSIZE || this.height > MAXSIZE) {
            // If it's not, show an error and hide the images.
            $(this).hide();
            $("#approx-image").hide();
            $("#error").text("Your image is too large. Please choose an image less than " + MAXSIZE + "x" + MAXSIZE + ".");
        } else {
            // If it is, show the images and stats.
            $("#error").text("");
            $(this).show();
            $("#approx-image").show();
            $.get("/stats", function(data) {
                console.log(data);
            });
        }
        rendering = false;
    })
    console.log("Rendering", imgUrl);

  }
});