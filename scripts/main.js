init = function() {
  document.write("Starting");
  CANVAS = $("#canvas")[0].getContext("2d");
  TILE_WIDTH = 20;
  TILE_HEIGHT = 20;

  bg = new Rect(20, 20, 255, 255, 255);
  for (var i = 18; i >= 1; i-=2) {
    var rect = Rect.rnd(i, i);
    bg = bg.add(rnd(0, 20-rect.w), rnd(0, 20-rect.h), rect);
  }
  draw(bg);
};

draw = function(rect) {
  for (var x = 0; x < rect.w; x++) {
    for (var y = 0; y < rect.h; y++) {
      CANVAS.fillStyle = rect.get(x, y).hex;
      CANVAS.fillRect(
        x*TILE_WIDTH,
        y*TILE_HEIGHT,
        TILE_WIDTH,
        TILE_HEIGHT);
    }
  }
};