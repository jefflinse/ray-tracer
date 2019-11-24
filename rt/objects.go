package rt

import (
	"math"
)

// A Sphere represents a sphere.
type Sphere struct {
	Transform Matrix
}

// NewSphere creates a new Sphere.
func NewSphere() *Sphere {
	return &Sphere{
		Transform: NewTransform(),
	}
}

// Intersect returns a set of points where a ray intersects the sphere.
func (s *Sphere) Intersect(ray *Ray) IntersectionSet {
	r := ray.Transform(s.Transform.Inverse())

	sphereToRay := r.Origin.Subtract(Origin())
	a := r.Direction.Dot(r.Direction)
	b := 2 * r.Direction.Dot(sphereToRay)
	c := sphereToRay.Dot(sphereToRay) - 1

	discriminant := math.Pow(b, 2) - 4*a*c
	if discriminant < 0 {
		return IntersectionSet{}
	}

	return IntersectionSet{
		NewIntersection((-b-math.Sqrt(discriminant))/(2*a), s),
		NewIntersection((-b+math.Sqrt(discriminant))/(2*a), s),
	}
}

// NormalAt returns the normal vector from the sphere for a point p.
func (s *Sphere) NormalAt(worldPoint Tuple) Tuple {
	objectPoint := s.Transform.Inverse().MultiplyTuple(worldPoint)
	objectNormal := objectPoint.Subtract(Origin())
	worldNormal := s.Transform.Inverse().Transpose().MultiplyTuple(objectNormal)
	worldNormal[3] = 0
	return worldNormal.Normalize()
}
