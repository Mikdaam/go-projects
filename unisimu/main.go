package main

import (
	"gioui.org/layout"
	"gioui.org/op/paint"
	"golang.org/x/image/colornames"
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/op"
)

// Body represents a celestial body in the simulation
type Body struct {
	Mass     float64 // Mass of the body
	Radius   float32 // Radius of the body for visualization
	Position Vec2    // Position of the body
	Velocity Vec2    // Velocity of the body
	Color    color.RGBA
}

// Universe represents the N-body simulation environment
type Universe struct {
	G           float64 // Gravitational constant
	Bodies      []Body  // List of bodies in the simulation
	SimulationG *op.Ops // Ops for Gio rendering
}

// Vec2 represents a 2D vector
type Vec2 struct {
	X, Y float64
}

// CalculateForce calculates the gravitational force between two bodies
func CalculateForce(m1, m2 float64, r1, r2 Vec2) Vec2 {
	dx := r2.X - r1.X
	dy := r2.Y - r1.Y
	dist := math.Sqrt(dx*dx + dy*dy)
	f := m1 * m2 / (dist * dist * dist)
	return Vec2{X: f * dx, Y: f * dy}
}

// IncreaseTime updates the positions and velocities of all the bodies in the universe based on the gravitational forces
func (u *Universe) IncreaseTime(dt float64) {
	for i, b1 := range u.Bodies {
		netForce := Vec2{} // Net force acting on the body

		for j, b2 := range u.Bodies {
			if i != j {
				force := CalculateForce(b1.Mass, b2.Mass, b1.Position, b2.Position)
				netForce.X += force.X
				netForce.Y += force.Y
			}
		}

		// Update velocity
		b1.Velocity.X += dt * netForce.X / b1.Mass
		b1.Velocity.Y += dt * netForce.Y / b1.Mass

		// Update position
		b1.Position.X += dt * b1.Velocity.X
		b1.Position.Y += dt * b1.Velocity.Y

		// Update the body
		u.Bodies[i] = b1
	}
}

// SimulationWindow creates and runs the simulation window
func SimulationWindow(universe *Universe) error {
	window := app.NewWindow(
		app.Title("N-Body Simulation"),
		app.Size(800, 600),
	)
	/*th := material.NewTheme(gofont.Collection())
	cp := colorpicker.New()*/
	var err error

	/*gtx := &layout.Context{
		Ops:         new(op.Ops),
		Constraints: layout.Exact(image.Pt(800, 600)),
	}*/
	// Simulation loop
	for e := range window.Events() {
		if e, ok := e.(key.Event); ok && e.State == key.Press {
			if e.Name == "q" {
				/*window.Close()
				return nil*/
				log.Println("Should kill the app")
			}
		}

		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(universe.SimulationG, e)
			universe.SimulationG.Reset()

			for i := range universe.Bodies {
				// Draw the body
				b := universe.Bodies[i]
				pos := image.Pt(int(b.Position.X), int(b.Position.Y))
				op.Offset(pos).Add(gtx.Ops)
				paint.ColorOp{Color: color.NRGBA(b.Color)}.Add(gtx.Ops)
				paint.PaintOp{}.Add(gtx.Ops)
			}

			e.Frame(gtx.Ops)
			break
		}
	}
	return err
}

func main() {
	bodies := []Body{
		{Mass: 100, Radius: 10, Color: colornames.Blue, Position: Vec2{X: 200, Y: 300}, Velocity: Vec2{X: 0, Y: 0}},
		{Mass: 50, Radius: 8, Color: colornames.Red, Position: Vec2{X: 400, Y: 300}, Velocity: Vec2{X: 0, Y: 0}},
		// Add more bodies as needed
	}

	// Add 100 more bodies
	for i := 0; i < 100; i++ {
		mass := rand.Float64() * 100
		radius := rand.Float32() * 5
		bodyColor := color.RGBA{R: uint8(rand.Intn(256)), G: uint8(rand.Intn(256)), B: uint8(rand.Intn(256)), A: 255}
		position := Vec2{X: rand.Float64() * 800, Y: rand.Float64() * 600}
		velocity := Vec2{X: rand.Float64() * 10, Y: rand.Float64() * 10}
		body := Body{Mass: mass, Radius: radius, Color: bodyColor, Position: position, Velocity: velocity}
		bodies = append(bodies, body)
	}

	universe := &Universe{
		G:      6.67430e-11,
		Bodies: bodies,
	}

	if err := SimulationWindow(universe); err != nil {
		log.Fatal(err)
	}
}
