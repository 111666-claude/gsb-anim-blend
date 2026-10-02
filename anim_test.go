package anim

import "testing"

func TestEmptyMixerHasNoWeights(t *testing.T) {
	if len(New("idle").Sample()) != 0 {
		t.Fatal("空混合器不应该有权重")
	}
}

func TestSingleClipKeepsItsWeight(t *testing.T) {
	mixer := New("idle")
	mixer.AddClip("run")
	layer := mixer.AddLayer()
	mixer.Play(layer, "run", 1)
	if got := mixer.Sample()["run"]; got != 1 {
		t.Fatalf("单剪辑权重应该是 1：%v", got)
	}
}

func TestUnknownLayerIsIgnored(t *testing.T) {
	mixer := New("idle")
	mixer.AddClip("run")
	mixer.Play(5, "run", 1)
	if len(mixer.Sample()) != 0 {
		t.Fatal("不存在的层不应该有权重")
	}
}

func TestMissingStartsAtZero(t *testing.T) {
	if New("idle").Missing() != 0 {
		t.Fatal("Missing 初始应该是 0")
	}
}

func TestScannedStartsAtZero(t *testing.T) {
	if New("idle").Scanned() != 0 {
		t.Fatal("Scanned 初始应该是 0")
	}
}
