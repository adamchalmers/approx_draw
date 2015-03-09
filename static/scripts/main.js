init = function() {
  CANVAS = $("#canvas")[0].getContext("2d");
  TILE_WIDTH = 1;
  TILE_HEIGHT = 1;
  MUTATIONS_PER_ITERATION = 60;
  ITERATIONS = 40;

  target = loadImage();
  $("#canvas").attr("width", target.w);
  $("#canvas").attr("height", target.h);

  // Start off our approximation with a white rectangle.
  var approxImage = new Rect(target.w, target.h, 255, 255, 255);
  var min = approxImage.distFrom(target);
  var bestCanvas = approxImage;
  var start = Date.now();
  console.log(min);
  console.log("-----");

  for (var j = 0; j < ITERATIONS; j++) {

    // Try MUTATIONS_PER_ITERATION different mutations, keep the best.
    for (var i = 0; i < MUTATIONS_PER_ITERATION; i++) {
      var rect = Rect.rnd(approxImage.w, approxImage.h);
      var mutation = approxImage.add(rnd(0, target.w-rect.w), rnd(0, target.h-rect.h), rect);
      var score = mutation.distFrom(target);
      if (score < min) {
        min = score;
        bestCanvas = mutation;
      }
    }
    approxImage = bestCanvas;
    console.log(min);
  }
  var timeTaken = (Date.now() - start)/1000;
  document.write(timeTaken + " seconds.");
  draw(bestCanvas, CANVAS);
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

  // Set up a rectangle with the pixel data from the image.
  var rect = new Rect(data.width, data.height, 0, 0, 0);
  for(var y = 0; y < img.height; y++) {
    for(var x = 0; x < img.width; x++) {
      var r = data.data[((img.width * y) + x) * 4];
      var g = data.data[((img.width * y) + x) * 4 + 1];
      var b = data.data[((img.width * y) + x) * 4 + 2];
      rect.set(x, y, new Color(r, g, b));
    }
  }
  draw(rect, ctx);
  return rect;
}