package save

import "encoding/json"

type Data struct {
	BestLevel        map[string]int  `json:"best_level"`
	BestRun          int             `json:"best_run,omitempty"`
	Secrets          map[string]bool `json:"secrets"`
	AudioOffsetMs    int             `json:"audio_offset_ms,omitempty"`
	Language         string          `json:"language,omitempty"`
	PressureUnlocked bool            `json:"pressure_unlocked,omitempty"`
	PressureCleared  int             `json:"pressure_cleared,omitempty"`
}

type Storage interface {
	Load() (Data, error)
	Save(Data) error
}

func (d Data) Normalized() Data {
	if d.BestLevel == nil {
		d.BestLevel = map[string]int{}
	}
	if d.Secrets == nil {
		d.Secrets = map[string]bool{}
	}
	return d
}

func decode(raw []byte) (Data, error) {
	var d Data
	if len(raw) == 0 {
		return d.Normalized(), nil
	}
	err := json.Unmarshal(raw, &d)
	return d.Normalized(), err
}

func encode(d Data) ([]byte, error) {
	return json.Marshal(d.Normalized())
}

type Memory struct {
	Stored Data
	Saves  int
}

func (m *Memory) Load() (Data, error) {
	raw, err := encode(m.Stored)
	if err != nil {
		return Data{}, err
	}
	return decode(raw)
}

func (m *Memory) Save(d Data) error {
	raw, err := encode(d)
	if err != nil {
		return err
	}
	m.Saves++
	m.Stored, err = decode(raw)
	return err
}
