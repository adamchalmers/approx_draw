function Rect(w, h, r, g, b) {
  if (w < 1 || h < 1) {
    throw "Can't have non-positive dimensions.";
  }
  this.w = w;
  this.h = h;
  this.pixels = new Array(w);
  for (var i = 0; i < w; i++) {
    this.pixels[i] = new Array(h);
    for (var j = 0; j < h; j++) {
      this.pixels[i][j] = new Color(r, g, b);
    }
  }
}

Rect.prototype.toString = function() {
  return "Rect (" + this.w + ", " + this.h + ")";
}

// Getter for rect color at a point.
Rect.prototype.get = function(x, y) {
  if (x >= this.w || y >= this.h || x < 0 || y < 0) {
    throw "Out of bounds access to " + this.toString();
  }
  return this.pixels[x][y];
}

// Setter for rect color at a point.
Rect.prototype.set = function(x, y, color) {
  if (x >= this.w || y >= this.h || x < 0 || y < 0) {
    throw "Out of bounds access to " + this.toString();
  }
  this.pixels[x][y] = color;
}

/**
 * Overlays the other rectangle on top of this one, from (x,y).
 * Returns the combined rectangle.
 */
Rect.prototype.addAs = function(x, y, w, h, color) {
  if (x + w > this.w || y + h > this.h) {
    throw "Other rect is too large!";
  } else {

    var newRect = new Rect(this.w, this.h, 0, 0, 0);

    // Copy this rect into the new rect
    for (var i = 0; i < this.w; i++) {
      for (var j = 0; j < this.h; j++) {
        newRect.set(i, j, this.get(i,j));
      }
    }

    // Overwrite new rect with the other rect's colors
    for (var i = 0; i < w; i++) {
      for (var j = 0; j < h; j++) {
        newRect.set(i + x, j + y, color);
      }
    }

    return newRect;
  }
}

Rect.prototype.distFrom = function(other) {
  if (this.w != other.w || this.h != other.h) {
    throw "Can only find distance between equal-dimension rectangles."
  }
  var dist = 0;
  for (var i = 0; i < this.w; i++) {
    for (var j = 0; j < this.h; j++) {
      dist += this.get(i,j).distFrom(other.get(i,j));
    }
  }
  return dist;
}

function Color(r, g, b) {
    this.r = r;
    this.g = g;
    this.b = b;
    this.hex = "#" + ((1 << 24) + (r << 16) + (g << 8) + b).toString(16).slice(1);
}

Color.rnd = function() {
  var c = new Color(rnd(0, 255), rnd(0, 255), rnd(0, 255));
  return c;
}

Color.prototype.distFrom = function(other) {
  return Math.abs(this.r - other.r) + Math.abs(this.g - other.g) + Math.abs(this.b - other.b);
}

function rnd(low, high) {
  return Math.floor(Math.random()*(high-low) + low);
}