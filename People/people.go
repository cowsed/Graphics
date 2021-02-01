package people

import "fmt"
import "../Rendering"

type Person struct{
	Name string
	X,Y,Z int
	Renderer *render.ActorRenderer
}


func (p Person) UpdateRenderAll(enabled bool){ //, sprite offsett

        (p.Renderer).UpdatePos(p.X,p.Y,p.Z)

		//Something weird is happening where passing enabled doesnt update it but just saying false does
    	(*p.Renderer).UpdateVisibility(enabled)
		fmt.Println("Second", (p.Renderer).X,(p.Renderer).Y,(p.Renderer).Z)
		fmt.Println("Second", (*p.Renderer).Visible)
}

//Updates the rendered position to be equal to the persons position
func (p Person) UpdateRenderPos(){ //, sprite offsett
        p.Renderer.UpdatePos(p.X,p.Y,p.Z)
}

func (p Person) UpdateVisibility(visible bool){
	//fmt.Println("To: ",visible)
    (p.Renderer).UpdateVisibility(visible)
	//fmt.Println("Outcome: ",(*p.Renderer))

}
