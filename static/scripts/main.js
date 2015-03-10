init = function() {
  CANVAS = $("#canvas")[0].getContext("2d");
  TILE_WIDTH = 1;
  TILE_HEIGHT = 1;
  MUTATIONS_PER_ITERATION = 6000;
  ITERATIONS = 40;

  var results = loadImage();
  var target = results[0];
  var colorsInPicture = results[1];
  $("#canvas").attr("width", target.w);
  $("#canvas").attr("height", target.h);
  var approxImage = new Rect(target.w, target.h, 0, 255, 255);
  draw(approxImage, CANVAS);
  approximateImage(target, colorsInPicture);
}

approximateImage = function(target, colorsInPicture) {

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
  document.write(timeTaken + " seconds.");
  draw(approxImage, CANVAS);
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
loadImage = function() {
  var img = $('#target-image')[0];
  var canvas = $('#photo')[0];

  // Draw the image into the canvas
  canvas.width = img.width;
  canvas.height = img.height;
  var ctx = canvas.getContext('2d');
  ctx.drawImage(img, 0, 0, img.width, img.height);
  var data = ctx.getImageData(0, 0, $('#photo').attr('width'), $('#photo').attr('height'));

  var colorsInPicture = {};
  colorsObject = [];

  // Set up a rectangle with the pixel data from the image.
  var rect = new Rect(data.width, data.height, 0, 0, 0);
  for(var y = 0; y < img.height; y++) {
    for(var x = 0; x < img.width; x++) {
      var r = data.data[((img.width * y) + x) * 4];
      var g = data.data[((img.width * y) + x) * 4 + 1];
      var b = data.data[((img.width * y) + x) * 4 + 2];
      var color = new Color(r, g, b);
      rect.set(x, y, color);
      if (!colorsInPicture[color.hex]) {
        colorsInPicture[color.hex] = true;
        colorsObject.push(color);
      }
    }
  }
  console.log(colorsInPicture);
  draw(rect, ctx);
  return [rect, colorsObject];
}