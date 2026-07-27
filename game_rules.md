# Game Rules

This file lists all the rules the game should respect. Write a rule here before implementing it, and make sure it's tested to confirm it holds.

Rules stay broad — state *what* must happen, not *how*. e.g. "the player can jump," not what triggers the jump or how it's calculated.

Rules stay free of balancing numbers (lives, cooldowns, thresholds, caps, etc.) — those values will change as the game gets tuned, but the rule itself shouldn't.

## Game Description

Bagre Wes is a 2D game where Wes the catfish (the player) must eat jellyfish to win. Wes must survive the environment and other creatures encountered along the way while eating as many jellyfish as possible.

The game is a path of maps, each with multiple stages. Eating jellyfish is what advances Wes through stages — each stage can raise the difficulty and also make the player stronger to face it.

## Rules

### Wes (Player)

1. Wes can move around with W, A, S, D and cannot pass beyond the screen limits.
2. Wes can attack by pressing space; the attack has a triangle hit box, and he can't start a new attack until the current one finishes.
3. Wes has a limited number of lives and loses one each time an attacking jellyfish touches him.
4. Wes eats every jellyfish caught in his attack hit box; each one eaten counts toward his total.

### Jellyfish (Smack)

1. Jellyfish move together as a flocking school — grouping up, matching direction, and keeping distance from close neighbors.
2. Jellyfish flee from Wes when he gets close.
3. Jellyfish never stop moving, and never move faster than their cruising speed.
4. Jellyfish can enter an attack state; while attacking, touching Wes damages him.
5. A jellyfish's attack state lasts a few seconds before it reverts to normal.
6. Eaten jellyfish play a death animation before disappearing from the game.
7. Only a limited number of jellyfish can be attacking at the same time.

### Stages

1. Stage effects are cumulative — effects from stage 1 remain active in later stages.
2. Stage 1 (electric jellyfish) begins once Wes has eaten enough jellyfish.
3. Once stage 1 begins, a random group of jellyfish is set to attack on a recurring timer.

