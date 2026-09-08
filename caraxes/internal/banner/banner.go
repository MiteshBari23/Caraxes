// Package banner renders the glowing "CARAXES" logo shown after `caraxes init`.
package banner

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
)

// 5x7 dot-matrix glyphs. '#' is lit, '.' is dark.
var glyphs = map[rune][]string{
	'C': {
		".###.",
		"#...#",
		"#....",
		"#....",
		"#....",
		"#...#",
		".###.",
	},
	'A': {
		".###.",
		"#...#",
		"#...#",
		"#####",
		"#...#",
		"#...#",
		"#...#",
	},
	'R': {
		"####.",
		"#...#",
		"#...#",
		"####.",
		"#.#..",
		"#..#.",
		"#...#",
	},
	'X': {
		"#...#",
		"#...#",
		".#.#.",
		"..#..",
		".#.#.",
		"#...#",
		"#...#",
	},
	'E': {
		"#####",
		"#....",
		"#....",
		"####.",
		"#....",
		"#....",
		"#####",
	},
	'S': {
		".####",
		"#....",
		"#....",
		".###.",
		"....#",
		"....#",
		"####.",
	},
}

const word = "CARAXES"
const glyphHeight = 7

// fireRamp goes from ember red (top of the letters) to a hot yellow-white
// (bottom), giving the impression of heat rising off the text.
var fireRamp = []int{196, 202, 208, 214, 220, 226, 230}

// shadowRamp is a dark, smoldering echo of fireRamp used for the drop
// shadow cast behind the letters.
var shadowRamp = []int{52, 52, 58, 58, 94, 94, 100}

// embers are faint background particles printed above/below the logo,
// echoing the dotted texture of a hot, glowing surface.
var emberChars = []rune(".·:'`,")

const (
	shadowDX = 2 // columns the shadow is offset to the right
	shadowDY = 1 // rows the shadow is offset downward
)

func noColor() bool {
	return os.Getenv("NO_COLOR") != ""
}

func fg256(code int) string {
	return fmt.Sprintf("\x1b[38;5;%dm", code)
}

const reset = "\x1b[0m"
const bold = "\x1b[1m"
const dim = "\x1b[2m"

type cell struct {
	set   bool
	ch    rune
	color int
	bold  bool
	dim   bool
}

func emberLine(width int, colorCode int, plain bool) string {
	var b strings.Builder
	for i := 0; i < width; i++ {
		if rand.Intn(6) == 0 {
			ch := emberChars[rand.Intn(len(emberChars))]
			if plain {
				b.WriteRune(ch)
			} else {
				b.WriteString(dim)
				b.WriteString(fg256(colorCode))
				b.WriteRune(ch)
				b.WriteString(reset)
			}
		} else {
			b.WriteByte(' ')
		}
	}
	return b.String()
}

// Print writes the glowing, drop-shadowed CARAXES banner to stdout.
func Print() {
	plain := noColor()

	// Build the source dot-matrix grid for the whole word.
	rows := make([]string, glyphHeight)
	for _, r := range word {
		glyph, ok := glyphs[r]
		if !ok {
			continue
		}
		for i := 0; i < glyphHeight; i++ {
			rows[i] += glyph[i] + "."
		}
	}
	srcCols := len(rows[0])
	pixelWidth := srcCols * 2

	canvasWidth := pixelWidth + shadowDX
	canvasHeight := glyphHeight + shadowDY

	canvas := make([][]cell, canvasHeight)
	for i := range canvas {
		canvas[i] = make([]cell, canvasWidth)
	}

	// Pass 1: cast the shadow, offset down-right of the glyphs.
	for r := 0; r < glyphHeight; r++ {
		for c := 0; c < srcCols; c++ {
			if rows[r][c] != '#' {
				continue
			}
			sr, sc := r+shadowDY, c*2+shadowDX
			shColor := shadowRamp[r]
			canvas[sr][sc] = cell{set: true, ch: '█', color: shColor, dim: true}
			canvas[sr][sc+1] = cell{set: true, ch: '█', color: shColor, dim: true}
		}
	}

	// Pass 2: draw the glowing glyphs on top, overwriting the shadow
	// wherever they overlap so only the fringe of the shadow peeks out.
	for r := 0; r < glyphHeight; r++ {
		for c := 0; c < srcCols; c++ {
			if rows[r][c] != '#' {
				continue
			}
			color := fireRamp[r]
			if r+1 < len(fireRamp) && rand.Intn(3) == 0 {
				color = fireRamp[r+1]
			}
			cl := cell{set: true, ch: '█', color: color, bold: true}
			canvas[r][c*2] = cl
			canvas[r][c*2+1] = cl
		}
	}

	fmt.Println()
	fmt.Println(emberLine(canvasWidth, fireRamp[0], plain))

	for _, row := range canvas {
		var b strings.Builder
		for _, cl := range row {
			if !cl.set {
				b.WriteByte(' ')
				continue
			}
			if plain {
				b.WriteRune(cl.ch)
				continue
			}
			if cl.bold {
				b.WriteString(bold)
			} else if cl.dim {
				b.WriteString(dim)
			}
			b.WriteString(fg256(cl.color))
			b.WriteRune(cl.ch)
			b.WriteString(reset)
		}
		fmt.Println(b.String())
	}

	fmt.Println(emberLine(canvasWidth, fireRamp[len(fireRamp)-1], plain))
	fmt.Println()
}
