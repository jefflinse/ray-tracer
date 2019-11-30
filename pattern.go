package rt

import (
	"math"
)

// Color constants
var black = NewColor(0, 0, 0)
var white = NewColor(1, 1, 1)

// A Pattern is a pattern that can be applied to a surface.
type Pattern interface {
	At(point Tuple) Color
	AtObject(object Shape, point Tuple) Color
	GetTransform() Transformation
	SetTransform(transform Transformation)
}

// PatternProps contains properties common to all patterns.
type PatternProps struct {
	A         Pattern
	B         Pattern
	Transform Transformation
	p         Pattern
}

// NewPatternProps creates a new PatternProps.
func NewPatternProps(a Pattern, b Pattern) PatternProps {
	return PatternProps{
		A:         a,
		B:         b,
		Transform: NewTransform(),
	}
}

// GetTransform returns the pattern's transformation.
func (props *PatternProps) GetTransform() Transformation {
	return props.Transform
}

// SetTransform sets the pattern's transformation.
func (props *PatternProps) SetTransform(transform Transformation) {
	props.Transform = transform
}

// AtObject returns the pattern color on the specified object at the specified point.
func (props *PatternProps) AtObject(object Shape, worldPoint Tuple) Color {
	localPoint := object.GetTransform().Inverse().ApplyTo(worldPoint)
	patternPoint := props.p.GetTransform().Inverse().ApplyTo(localPoint)
	return props.p.At(patternPoint)
}

// A SolidPattern is just a single color.
type SolidPattern struct {
	PatternProps
	color Color
}

// NewSolidPattern creates a new SolidPattern.
func NewSolidPattern(color Color) *SolidPattern {
	pattern := &SolidPattern{NewPatternProps(nil, nil), color}
	pattern.p = pattern
	return pattern
}

// At returns the pattern color at the given point.
func (p *SolidPattern) At(point Tuple) Color {
	return p.color
}

// A BlendedPattern is blended combination of two other patterns.
type BlendedPattern struct {
	PatternProps
	patternA Pattern
	patternB Pattern
}

// NewBlendedPattern creates a new BlendedPatter.
func NewBlendedPattern(a Pattern, b Pattern) *BlendedPattern {
	pattern := &BlendedPattern{NewPatternProps(nil, nil), a, b}
	pattern.p = pattern
	return pattern
}

// AtObject returns the pattern color on the specified object at the specified point.
func (p *BlendedPattern) AtObject(object Shape, worldPoint Tuple) Color {
	localPoint := object.GetTransform().Inverse().ApplyTo(worldPoint)
	patternPointA := p.patternA.GetTransform().Inverse().ApplyTo(localPoint)
	colorA := p.patternA.At(patternPointA)
	patternPointB := p.GetTransform().Inverse().ApplyTo(localPoint)
	colorB := p.patternA.At(patternPointB)
	return colorA.AverageBlend(colorB)
}

// At returns the pattern color at the given point.
func (p *BlendedPattern) At(point Tuple) Color {
	return nil
}

// A StripePattern is a pattern of colors alternates in the X axis.
type StripePattern struct {
	PatternProps
}

// NewStripePattern creates a new StripePattern.
func NewStripePattern(a Pattern, b Pattern) *StripePattern {
	pattern := &StripePattern{NewPatternProps(a, b)}
	pattern.p = pattern
	return pattern
}

// At returns the pattern color at the given point.
func (p *StripePattern) At(point Tuple) Color {
	if int(math.Floor(point.X()))%2 == 0 {
		return p.A.At(point)
	}

	return p.B.At(point)
}

// A GradientPattern is a pattern that fades from one color to another.
type GradientPattern struct {
	PatternProps
}

// NewGradientPattern creates a new GradientPattern.
func NewGradientPattern(a Pattern, b Pattern) *GradientPattern {
	pattern := &GradientPattern{NewPatternProps(a, b)}
	pattern.p = pattern
	return pattern
}

// At returns the pattern color at the given point.
func (p *GradientPattern) At(point Tuple) Color {
	distance := p.B.At(point).Subtract(p.A.At(point))
	fraction := point.X() - math.Floor(point.X())
	return p.A.At(point).Add(distance.Multiply(fraction))
}

// A RingPattern is a pattern of alternating rings of color.
type RingPattern struct {
	PatternProps
}

// NewRingPattern creates a new RingPattern.
func NewRingPattern(a Pattern, b Pattern) *RingPattern {
	pattern := &RingPattern{NewPatternProps(a, b)}
	pattern.p = pattern
	return pattern
}

// At returns the pattern color at the given point.
func (p *RingPattern) At(point Tuple) Color {
	if int(math.Floor(math.Sqrt(math.Pow(point.X(), 2)+math.Pow(point.Z(), 2))))%2 == 0 {
		return p.A.At(point)
	}

	return p.B.At(point)
}

// A CheckerPattern is a pattern of alternating colors in all dimensions.
type CheckerPattern struct {
	PatternProps
}

// NewCheckerPattern creates a new RingPattern.
func NewCheckerPattern(a Pattern, b Pattern) *CheckerPattern {
	pattern := &CheckerPattern{NewPatternProps(a, b)}
	pattern.p = pattern
	return pattern
}

// At returns the pattern color at the given point.
func (p *CheckerPattern) At(point Tuple) Color {
	if (int(math.Floor(point.X()))+int(math.Floor(point.Y()))+int(math.Floor(point.Z())))%2 == 0 {
		return p.A.At(point)
	}

	return p.B.At(point)
}
