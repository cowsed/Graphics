# Architecture
A description of the architecture of this isometric visualiztion system 





## `Assets/`
Holds all the art (tilesets) now all home made (poorly)

## `Materials/`
Holds indices for tiles in the form of consts
(may be phased out soon once the renderer update that allows loading tiles from a file)


## `Rendering/`
Holds all that is responsible for rendering to the screen.

Chunks for holding the environment
ActorRenderer for holding characters or props that are not the environment

Way to think about it is you can stand in the same unit block as an actor renderer (a tree, an umbrella, a person) but not a tile from a chunk (dirt cube, stone cube) - this is important for pathfinding

- *Camera.go* - Files responsible for moving and viewing the world

- *Renderer.go* - Responsible for initializing the rendering by loading tilesets, occupying chunks, and creating arrays that need it

- *Chunk.go* - Holds code responsible for rendering and optimizing the drawing of a chunk's tiles and sprites

- *ActorRenderer.go* - Holds code responsible for the communication from an ai or a player with the renderer (sorta a WIP for now. its not the most functional or prtty )
- *UI.go* Holds code responsible for drawing the debug information and soon general information
- *general.go* holds code that is mostly responsible for changing from non-visual game space (arrays of arrays of arrays) to the isometric view of the world

## `Screenshots/`
What it says on the tin. Screenshots of the visualization