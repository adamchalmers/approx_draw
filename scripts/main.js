init = function(ITERATIONS, MUTATIONS_PER_ITERATION, IMG, url) {
  CANVAS = $("#canvas")[0].getContext("2d");
  TILE_WIDTH = 1;
  TILE_HEIGHT = 1;

  var results = loadImage(IMG, url, function(target, colorsArray) {
    $("#canvas").attr("width", target.w);
    $("#canvas").attr("height", target.h);
    $("#target-image").show();
    var approxImage = new Rect(target.w, target.h, 0, 0, 0);
    draw(approxImage, CANVAS);
    window.setTimeout(function() {
      approximateImage(target, colorsArray);
    }, 1000);
  });
}

