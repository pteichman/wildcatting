package game

import "sync/atomic"

type world struct {
	entities
	nameManager
	playerManager
}

type entity uint32

var None entity = 0

type entities struct {
	prev uint32
}

func (m *entities) NewEntity() entity {
	return entity(atomic.AddUint32(&m.prev, 1))
}

type playerManager struct {
	players []entity
}

func (m *playerManager) AddPlayer(e entity) {
	if _, ok := linfind(m.players, e); ok {
		return
	}
	m.players = append(m.players, e)
}

func (m *playerManager) IsPlayer(e entity) bool {
	_, ok := linfind(m.players, e)
	return ok
}

func (m *playerManager) PlayerOne() entity {
	if len(m.players) > 0 {
		return m.players[0]
	}
	return None
}

func (m *playerManager) Players() []entity {
	return m.players
}

type nameManager struct {
	index []entity
	names []string
}

func (m *nameManager) SetName(e entity, name string) {
	if i, ok := linfind(m.index, e); ok {
		m.names[i] = name
		return
	}
	m.index = append(m.index, e)
	m.names = append(m.names, name)
}

func (m *nameManager) Name(e entity) string {
	if i, ok := linfind(m.index, e); ok {
		return m.names[i]
	}
	return ""
}

func (m *nameManager) DelName(e entity) {
	if i, ok := linfind(m.index, e); ok {
		// Swap the deleted entity with the last entity, then reslice.
		last := len(m.names) - 1
		m.names[i], m.names[last] = m.names[last], m.names[i]
		m.index[i], m.index[last] = m.index[last], m.index[i]
		m.names = m.names[:last]
		m.index = m.index[:last]
	}
}

func linfind(index []entity, e entity) (int, bool) {
	for i, cur := range index {
		if cur == e {
			return i, true
		}
	}
	return 0, false
}
