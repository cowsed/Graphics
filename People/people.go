package people

import "../Rendering"

type Person struct{
	Name string
	X,Y,Z int
	Renderer *render.ActorRenderer
}

func (p Person) Render(){
	(*p.Renderer).RenderSprite(p.X,p.Y,p.Z)

}
