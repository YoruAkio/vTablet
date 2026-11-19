package pointer

type ScreenData struct {
	Uid    string `json:"uid"`
	Path   string `json:"path"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Name   string `json:"name,omitempty"`
}

type pointerBackend interface {
	Update(x, y int32, pressure uint32, tiltX, tiltY int32) error
	Destroy()
}

type AbsPointer struct {
	screenUid string
	backend   pointerBackend

	X        int32
	Y        int32
	Pressure uint32
	TiltX    int32
	TiltY    int32
}

func (p *AbsPointer) Update() error {
	if p.backend != nil {
		return p.backend.Update(p.X, p.Y, p.Pressure, p.TiltX, p.TiltY)
	}
	return nil
}

func (p *AbsPointer) Destroy() {
	if p.backend != nil {
		p.backend.Destroy()
	}
}
