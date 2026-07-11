# Notes

# Game in general
- We are in a loop where every tick of method update is my clock, every update call will refresh the screen
- Image.Geo.Translates is like MOVE image, where positive x goes to the right negative left, positive Y goes down and negative out of screen
- DrawImage pins an image by its top-left corner. that's why negative y puts things up
        -Y (up)
         ↑
-X ←─────┼─────→ +X
(left)   │   (right)
         ↓
        +Y (down)
- To animate a image we will perform a walk across the filmstrip, the image is a sequence of position of the element where we will
"walk" over each frame, so this is look like we animated the element.
- sips -g pixelWidth -g pixelHeight ./assets/jellyfish/Walk.png CLI command to extract image properties
- Each image can have many frames count diffirently and a humans must inspict to know this point.
- It's possiabel to calculate the framwWidht and height by image properties, but the frame count it's a const per image

---

# Tightened reference (drawing & animating a sprite)

Ref example: https://ebitengine.org/en/examples/animation.html

## Game loop
- Backend = reactive (wait for request). Game = a loop: ~60x/sec Ebiten calls `Update()` then `Draw()`.
- Screen is repainted from scratch every frame; nothing persists.
- `Update()` advances the world (`g.tick++`) = the clock. `Draw(screen)` paints one still of the current state.

## Coordinates
- Origin `(0,0)` = top-left. X grows right, Y grows DOWN. So `-X`=left, `-Y`=up.

## Simplest draw
```go
screen.DrawImage(img, nil)   // whole image at (0,0); nil = no transform
```

## Background — Fill (must be first in Draw)
```go
screen.Fill(color.RGBA{R, G, B, A})   // paints every pixel; drawn-after lands on top
```
- RGBA 0–255, `A` = opacity (background wants 255). Filled rect = the game screen; around it = OS window.

## Positioning — op + GeoM.Translate
```go
op := new(ebiten.DrawImageOptions)   // empty = acts like nil
op.GeoM.Translate(dx, dy)            // MOVE dx right, dy down (negatives = left/up)
```
- `GeoM` = geometry matrix; also does Scale + Rotate. Translate is plain addition on position.

### The pin (anchor) + why TWO translates
- Default pin = top-left corner, so the x,y you set lands the CORNER → not centered.
- #1 pin → image center (image-space): `Translate(-fw/2, -fh/2)` — shift up-left by HALF the image's own size. Value is relative to the IMAGE.
- #2 pin → target (screen-space): `Translate(screenWidth/2, screenHeight/2)` = 320,180.
- They ADD. Read: "make center the anchor, THEN put anchor at screen center."
- #1 alone parks it at top-left. Add Rotate/Scale BETWEEN the two so it transforms around center.
- Get size safely: `w, h := img.Bounds().Dx(), img.Bounds().Dy()`.

## Sprite sheet layout
- One PNG = many frames side by side (filmstrip). Walk.png = 192x48 = 4 frames of 48x48, one row.
- `frameOX/frameOY` = pixel where frame 0 starts (grid ORIGIN). Ours 0,0. Non-zero for margins, or when many animations share one file (frameOY picks the row).
- `frameWidth/frameHeight` = one frame's size. `walkFrameCount` = number of frames.

## Frame-picking math
```go
i := (g.tick / 5) % walkFrameCount
```
- `g.tick` = odometer, +1 per Update. `/5` = hold each pose 5 ticks (~12fps vs 60, else a blur).
- `% count` = wrap to 0 after last frame → loops; result always 0..count-1 (never out of bounds).
- e.g. tick=40 → 40/5=8 → 8%4=0 → frame 0.

## Cropping — SubImage
	// SubImage returns a cropped VIEW into the sheet — the exact pixel rectangle
	// (sx,sy)..(sx+48,sy+48), i.e. one 48x48 frame. No pixels are copied; it's a window.
	//
	// The sheet is one row: [frame0, frame1, frame2, frame3]
	// Each Draw call crops just ONE frame (the one `i` points to right now).
	// As g.tick advances over successive Draw calls, `i` steps 0→1→2→3→0…,
	// so the sequence of stills played over time reads as animation.
	//
	// Cut frame i (a 32x32 window) out of the sheet, then draw it with our transform.
## Where the "loop" actually is

- No `for` over frames in Draw — each Draw crops ONE frame. Progression happens ACROSS Draw calls as `tick` advances + the `%` wrap. Stills over time = motion (like film). The illusion lives in TIME.

One sentence: keep one clock (`tick`), use it to index into a filmstrip (`SubImage`), and use a transform (`GeoM`) to place that frame on the freshly-repainted screen — 60x/sec.
