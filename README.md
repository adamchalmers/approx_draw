# approx_draw
Web app that recreates an image using rectangles.

This uses a very simple hill-climbing algorithm to randomly place rectangles on a canvas, in a way which iteratively builds up an approximation of an image.

<h2>Algorithm</h2>
```
current approximation = empty canvas
for NUM_ITERATIONS
  for NUM_RECTANGLES_PER_ITERATION
    Generate a random_rectangle.
    Calculate the score of (current approximation + random_rectangle).
    if score is better than current best score from this iteration:
      best_score = score
      best_rectangle = random_rectangle
  current approximation += best rectangle.
display current approximation
```
The score is calculated like so:
```
score(approximate_img, target_img):
  score = 0
  for pixel, target_pixel in approximate_img, target_img
    score += abs(target_pixel.r - pixel.r)
    score += abs(target_pixel.g - pixel.g)
    score += abs(target_pixel.b - pixel.b)
  return score
```
Its performance is O(I*R*P), where I is iterations, R is rectangles per iteration, and P is pixels in the image. This suggests that if you want a good render, try downscaling your image first! I've gotten good results with images under 200x200 pixels.
<h2>Demo</h2>
![Image showing demo of approximation](https://cloud.githubusercontent.com/assets/5407457/6572281/0d9c1cfe-c765-11e4-8060-257ff2e5d688.jpg)
![Image showing demo of approximation](https://cloud.githubusercontent.com/assets/5407457/6575528/911fc346-c783-11e4-9cce-521c92305616.png)
![Image showing demo of approximation](https://cloud.githubusercontent.com/assets/5407457/6575534/978b74fa-c783-11e4-9071-5e8ac8a7c801.png)
![Image showing demo of approximation](https://cloud.githubusercontent.com/assets/5407457/6575532/93e73046-c783-11e4-93f6-aafc898a4934.png)

<h2>To do</h2>
 - Make GIFs that show the images being built up, rectangle by rectangle
 - Further optimize by cutting off the scoreWithMutation calculation once it's higher than the cached score by enough that the remaining pixels couldn't possibly even it back out.
 - Optimise the server code - why is the JS version twice as fast as the Go version?
 - Resize large images
