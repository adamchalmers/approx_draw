init = function(ITERATIONS, MUTATIONS_PER_ITERATION, IMG) {
  CANVAS = $("#canvas")[0].getContext("2d");
  TILE_WIDTH = 1;
  TILE_HEIGHT = 1;

  var results = loadImage(IMG, function(target, colorsInPicture) {
    $("#canvas").attr("width", target.w);
    $("#canvas").attr("height", target.h);
    var approxImage = new Rect(target.w, target.h, 0, 0, 0);
    draw(approxImage, CANVAS);
    window.setTimeout(function() {
      approximateImage(target, colorsInPicture);
    }, 1000);
  });
}

approximateImage = function(target, colorsInPicture) {

  console.log(target);

  // Start off our approximation with a white rectangle.
  var approxImage = new Rect(target.w, target.h, 255, 255, 255);
  var min = approxImage.score(target);
  var bestMutation = undefined; // [x, y, w, h, color]
  var start = Date.now();

  for (var j = 0; j < ITERATIONS; j++) {

    var cachedScore = min;

    // Try MUTATIONS_PER_ITERATION different mutations, keep the best.
    for (var i = 0; i < MUTATIONS_PER_ITERATION; i++) {

      // Choose the mutated block's properties
      var w = rnd(0, approxImage.w);
      var h = rnd(0, approxImage.h);
      var x = rnd(0, target.w-w);
      var y = rnd(0, target.h-h);
      var color = colorsInPicture[rnd(0, colorsInPicture.length)];

      // Score the mutation
      var score = approxImage.scoreWithMutation(x, y, w, h, color, target, cachedScore);

      // Compare it to the best
      if (score < min) {
        min = score;
        bestMutation = [x, y, w, h, color];
      }
    }
    approxImage.mutate(bestMutation[0], bestMutation[1], bestMutation[2],
                       bestMutation[3], bestMutation[4]);
  }
  var timeTaken = (Date.now() - start)/1000;

  draw(approxImage, CANVAS);
  $("#time-info").text(timeTaken + " seconds, " + ITERATIONS + " rectangles, " + MUTATIONS_PER_ITERATION + " population, " + target.w + "x" + target.h);
  console.log(min/1000000);
};

// Draw a rectangle using a canvas 2d context.
draw = function(rect, ctx) {
  for (var x = 0; x < rect.w; x++) {
    for (var y = 0; y < rect.h; y++) {
      ctx.fillStyle = rect.get(x, y).hex;
      ctx.fillRect(
        x*TILE_WIDTH,
        y*TILE_HEIGHT,
        TILE_WIDTH,
        TILE_HEIGHT);
    }
  }
};

// Loads the image from the DOM into a Rect object.
loadImage = function(IMG, callback) {
  var canvas = $('#photo')[0];
  var img = IMG[0];

  // Draw the image into the canvas
  canvas.width = img.width;
  canvas.height = img.height;
  var ctx = canvas.getContext('2d');
  ctx.drawImage(img, 0, 0, img.width, img.height);
  var data = ctx.getImageData(0, 0, $('#photo').attr('width'), $('#photo').attr('height'));
  return unpackImage(data, img, ctx, callback);
}

unpackImage = function(data, img, ctx, callback) {
  var colorsInPicture = {};
  var colorsObject = [];

  // Set up a rectangle with the pixel data from the image.
  $.get("image/img?url=img/monalisa.png", function(_data) {
    data = JSON.parse(_data);
    var rgb = data.Rgb
    var w = data.W
    var h = data.H
    var rect = new Rect(w, h, 0, 0, 0);

    // Loop over the rgb array, unpack into Colors.
    for(var y = 0; y < h; y++) {
      for(var x = 0; x < w; x++) {
        var r = rgb[((w * y) + x) * 3];
        var g = rgb[((w * y) + x) * 3 + 1];
        var b = rgb[((w * y) + x) * 3 + 2];
        var color = new Color(r, g, b);
        rect.set(x, y, color);
        if (!colorsInPicture[color.hex]) {
          colorsInPicture[color.hex] = true;
          colorsObject.push(color);
        }
      }
    }
    draw(rect, ctx);
    callback(rect, colorsObject);
  });
}