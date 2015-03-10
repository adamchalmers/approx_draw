# approx_draw
Web app that recreates an image using rectangles.

This uses a very simple hill-climbing algorithm to randomly place rectangles on a canvas, in a way which iteratively builds up an approximation of an image.


The algorithm is:
```
for NUM_ITERATIONS
  for NUM_RECTANGLES_PER_ITERATION
    Generate a random_rectangle.
    Calculate the score of (current approximation + random_rectangle).
    if score is better than current best score from this iteration:
      best_score = score
      best_rectangle = random_rectangle
  current approximation += best rectangle.
  ```
