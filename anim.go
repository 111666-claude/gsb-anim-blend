// Package anim 是动画分层混合：层内归一化、上层覆盖、缺失剪辑回退、按真实时间过渡。
// 缺陷：层内不归一化、各层平均分配、缺失剪辑被静默忽略、过渡按固定帧长推进、
// 采样时扫描全部剪辑。
package anim

// Mixer 是动画混合器。
type Mixer struct {
	fallback string
	clips    map[string]bool
	layers   []*layer
	trans    *transition
	missing  int
	scanned  int
}

type play struct {
	clip   string
	weight float64
}

type layer struct {
	plays []play
}

type transition struct {
	layer     int
	from      string
	to        string
	durMs     float64
	elapsedMs float64
}

// New 建混合器，fallback 是缺失剪辑的回退名。
func New(fallback string) *Mixer {
	return &Mixer{fallback: fallback, clips: map[string]bool{}}
}

// AddClip 登记一个存在的剪辑。
func (m *Mixer) AddClip(name string) { m.clips[name] = true }

// AddLayer 加一层，返回层号。
func (m *Mixer) AddLayer() int {
	m.layers = append(m.layers, &layer{})
	return len(m.layers) - 1
}

// Play 在某一层放一个剪辑。
// 缺陷：剪辑不存在时直接忽略，不计数也不回退。
func (m *Mixer) Play(layerIndex int, clip string, weight float64) {
	if layerIndex < 0 || layerIndex >= len(m.layers) {
		return
	}
	if !m.clips[clip] {
		return
	}
	m.layers[layerIndex].plays = append(m.layers[layerIndex].plays, play{clip: clip, weight: weight})
}

// BeginTransition 在某层开始一次过渡。
func (m *Mixer) BeginTransition(layerIndex int, from, to string, durMs float64) {
	if layerIndex < 0 || layerIndex >= len(m.layers) {
		return
	}
	m.trans = &transition{layer: layerIndex, from: from, to: to, durMs: durMs}
}

// Tick 推进过渡。
// 缺陷：按固定帧长推进，和 dtMs 无关。
func (m *Mixer) Tick(dtMs float64) {
	if m.trans == nil {
		return
	}
	m.trans.elapsedMs += 16.67
}

func (m *Mixer) playsOf(index int) []play {
	if m.trans != nil && m.trans.layer == index {
		t := 1.0
		if m.trans.durMs > 0 {
			t = m.trans.elapsedMs / m.trans.durMs
		}
		if t > 1 {
			t = 1
		}
		return []play{{clip: m.trans.from, weight: 1 - t}, {clip: m.trans.to, weight: t}}
	}
	return m.layers[index].plays
}

// Sample 返回每个剪辑的有效权重。
// 缺陷：各层平均分配，层内权重也不归一化。
func (m *Mixer) Sample() map[string]float64 {
	m.scanned += len(m.clips)
	out := map[string]float64{}
	if len(m.layers) == 0 {
		return out
	}
	share := 1.0 / float64(len(m.layers))
	for index := range m.layers {
		for _, item := range m.playsOf(index) {
			out[item.clip] += share * item.weight
		}
	}
	return out
}

// Missing 是缺失剪辑被回退的次数。
func (m *Mixer) Missing() int { return m.missing }

// Scanned 是累计扫描量（规模观测）。
func (m *Mixer) Scanned() int { return m.scanned }
