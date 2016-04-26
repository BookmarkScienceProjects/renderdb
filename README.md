# Introduction

RenderDB is a HTTP server hosting a REST API for performing efficient geometric queries
for static 3D geometry and metadata.

It's based on R-trees for fast lookups and efficiently handles frequent queries, with 
infrequent updates. 

The API will support queries based on camera viewport, e.g. to yield all visible objects
from a certain viewpoint. To achieve this occlusion culling heurestic will be applied in
addition to view frustum culling. The server will support simple level of detail schemes
to reduce the amount of geometry.


# API

The API is layered into 'worlds', 'layers', 'scenes' and 'objects'.

- A 'world' typically defines a separate project, e.g. a building site.
  Each world can have many layers.
- A 'layer' typically defines a group of data, e.g. all the plumbing in the building.
  Each layer can have many scenes.
- A 'scene' is the smallest distingushable feature and can e.g. represents
  all plumbing in floor 2. Each scene can have many 'objects'.
- 'Objects' are geometric 3D entities that can be rendered. Each object
  can have JSON metadata which can be used for dynamic filtering.
  There is no API to add single objects or to query for objects based
  on IDs. To add objects a new 'scene' must be added.

## Data management endpoints

- POST 	/worlds
  Adds a new world
- GET  	/worlds
  Returns metadata for all known worlds
- GET  	/worlds/{id}
  Returns metadata the world with the given ID
- DELETE 	/worlds/{id} 	(Not implemented yet)
  Deletes the world with the given ID. Deletes all
  layers in the world.
- POST 	/worlds/{id}/layers
  Adds a new layer to the given world
- GET 		/worlds/{id}/layers
  Returns metadata for all layers in the world
- GET 		/worlds/{id}/layers/{id}
  Returns metadata for the layer with the given ID
- DELETE   /worlds/{id}/layers/{id} 	(Not implemented yet)
  Deletes the layer with the given ID. Deletes all
  scenes in the layer.
- POST 	/worlds/{id}/layers/{id}/scenes
  Adds a new scene to the given layer. Scenes are specified
  using the Wavefront OBJ-format and each group
  in the file is considered to be a separate object.
- PUT 		/worlds/{id}/layers/{id}/scenes/{id}	(Not implemented yet)
  Replaces all geometry in a scene. Supports the same
  formats as the POST request.
- GET 		/worlds/{id}/layers/{id}/scenes
  Returns metadata for all scenes in the layer.
- GET 		/worlds/{id}/layers/{id}/scenes/{id}
  Returns metadata for the given scene.
- DELETE 	/worlds/{id}/layers/{id}/scenes/{id}	(Not implemented yet)
  Deletes the scene with the given ID and all the objects
  in the scene.

## Geometry query endpoints
- GET /world/{id}/geometry?{filter}&{options}	(Not implemented yet)
  Gets all geometry in the world that matches the filter.
- GET /world/{id}/layers/{id}/geometry?{filter}&{options}	(Not implemented yet)
  Gets all geometry in the layer that matches the filter.

Filters is used to filter away unwanted data, e.g. based on location or distance
to camera.
Options are used to e.g. sort the results by distance to a camera, or
restrict the number of returned triangles.