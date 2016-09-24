package report

//Renderer ...
type Renderer func(node UrnComposite, times []int64) string

//RendererType ...
type RendererType string

//RendererRegistry ...
type RendererRegistry struct {
	renderers map[RendererType]Renderer
}

//NewRendererRegistry ...
func NewRendererRegistry() RendererRegistry {
	return RendererRegistry{
		renderers: map[RendererType]Renderer{},
	}
}

//Add ...
func (instance RendererRegistry) Add(rendererType RendererType, renderer Renderer) {
	instance.renderers[rendererType] = renderer
}

//Get ...
func (instance RendererRegistry) Get(rendererType RendererType) Renderer {
	return instance.renderers[rendererType]
}
