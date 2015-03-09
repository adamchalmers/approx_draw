QUnit.test("Rect basics", function(assert) {
  var r = new Rect(4, 2, 0, 0, 0);

  // Check the rectangle has expected dimensions.
  assert.ok(r.w == 4);
  assert.ok(r.h == 2);
  assert.equal(r.pixels.length, 4);
  assert.equal(r.pixels[0].length, 2);
  assert.equal(r.pixels[1].length, 2);
  assert.equal(r.pixels[2].length, 2);
  assert.equal(r.toString(), "Rect (4, 2)");
  assert.equal(r.get(0, 0).hex, "#000000");

  // Check that set and get have the right direction.
  r.set(2, 1, new Color(255, 255, 255));
  assert.equal(r.get(2, 1).hex, "#ffffff");
  assert.equal(r.get(1, 1).hex, "#000000");
  assert.equal(r.get(3, 1).hex, "#000000");
});

QUnit.test("Color basics", function(assert) {
  var c = new Color(255, 255, 255);
  assert.equal(c.hex, "#ffffff");
  var d = new Color(255, 0, 0);
  assert.equal(d.hex, "#ff0000");
  var e = new Color(0, 255, 0);
  assert.equal(e.hex, "#00ff00");
  assert.equal(d.distFrom(e), 510);
  assert.equal(e.distFrom(d), 510);
  assert.equal(c.distFrom(d), 510);
  var f = new Color(0, 0, 0);
  var g = new Color(0, 1, 1);
  assert.equal(g.distFrom(f), 2);
});

QUnit.test("Rect intersection", function(assert) {
  var smallRect = new Rect(10, 10, 0, 0, 0);
  var largeRect = new Rect(20, 20, 255, 255, 255);
  assert.throws(
    function() {
      smallRect.add(0, 0, largeRect);
    },
    /too large/,
    "Correctly refuses to add a large rect to a small one."
  );
  var newRect = largeRect.add(0, 0, smallRect);
  assert.equal(newRect.w, largeRect.w);
  assert.equal(newRect.h, largeRect.h);
  assert.equal(newRect.get(19, 19).hex, "#ffffff", "Original color is preserved.");
  assert.equal(newRect.get(0, 0).hex, "#000000", "New color correctly changed.");
});

QUnit.test("Random builders", function(assert) {
  assert.ok(Color.rnd().hex);
  assert.ok(Rect.rnd(10, 10).get(0, 0).hex);
})

QUnit.test("Rect distance", function(assert) {
  var black = new Rect(4, 4, 0, 0, 0);
  var nearlyBlack = new Rect(4, 4, 0, 10, 0);
  assert.equal(black.distFrom(nearlyBlack), 160);
  var white = new Rect(4, 4, 255, 255, 255);
  assert.equal(black.distFrom(white), 12240);
  var large = new Rect(5, 5, 0, 0, 0);
  assert.throws(
    function() {
      black.distFrom(large);
    },
    /dimension/);
})