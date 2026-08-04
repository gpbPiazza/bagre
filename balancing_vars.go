package main

import "image/color"

// All game balancing knobs live here, grouped by context. Tweak freely: tests
// derive expected values from these constants instead of hardcoding them.

// Screen / clock
const (
	screenWidth, screenHeight = 950, 550
	ticksPerSecond            = 60 // one game.Update() call = one tick
)

// Wes
const (
	wesSpeedRation  = 2.0
	wesStartingLife = 3
	wesBorderMargin = 10 // px wes keeps from the screen edge; 10 px is visually good

	wesAttackTicksPerPose = 6
	wesAttackFrameCount   = 6
	wesWalkTicksPerPose   = 5
	wesWalkFrameCount     = 4

	wesHurtTicksPerPose = 20
	wesHurtFrameCount   = 2

	// how long wes stays in hurt state
	wesHurtDurationTicks = 2 * ticksPerSecond
)

// JellyFish flocking
const (
	// How far a jelly senses neighbors. Keep it LOCAL: too big and every jelly
	// sees the whole population, so cohesion's "average position" becomes the
	// screen center and all of them collapse into one static blob. Small radius
	// = many sub-groups that drift, split, and re-form (a living school).
	// Lower bound: must exceed the drawn sprite size (48*0.5 = 24px) so a jelly
	// senses a neighbor BEFORE their images overlap.
	jellyViewRadius       = 20
	jellyAttackViewRadius = 25

	adjustRateAligment      = 0.15
	adjustRateCohesion      = 0.03
	adjustRateSeparation    = 0.20
	adjustRateSeparationWes = 0.45
	jellysCount             = 400

	// Cruising speed window: floor = never stall into a frozen blob,
	// ceiling = never blast off across the screen.
	jellyMinSpeed = 0.5
	jellyMaxSpeed = 1.5

	// Cap for the border bounce push: at pos <= 0 a raw 1/pos flips sign (or
	// hits +Inf), so distances are clamped to this before dividing.
	borderPushCap = 0.5
)

// JellyFish animation
const (
	jellyWalkTicksPerPose = 5
	jellyWalkFrameCount   = 4

	jellyDeathTicksPerPose = 10
	jellyDeathFrameCount   = 6

	jellyAttackTicksPerPose = 5
	jellyAttackFrameCount   = 4
)

// Stages
const (
	// jellys wes must eat before the electric jelly stage begins
	startEletricJellyAt = 20
	// how often a new random group of jellys is set to attack
	eletricAttackIntervalTicks = 5 * ticksPerSecond
	// how long a jelly stays in attack state
	jellyAttackDurationTicks = 4 * ticksPerSecond
	// hard cap of simultaneous attacking jellys
	maxAttackingJellys = 10
	// the smack is split into this many chunks; one random chunk attacks
	smackChunkCount = 4
)

// UI
const (
	counterFontSize          = 24
	counterPosX, counterPosY = 150, 10
	wesLifePosX, wesLifePosY = 10, 10
)

// Palette
var (
	darkGrey    = color.RGBA{R: 40, G: 45, B: 60, A: 255} //nolint:mnd // RGB palette values
	hitBoxGreen = color.RGBA{R: 0, G: 255, B: 0, A: 255}  //nolint:mnd // RGB palette values
)
