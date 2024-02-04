package main

type Broker struct {
	runevent    chan CellMsg
	subscribe   chan CellSub
	unsubscribe chan CellSub
	grid        *Grid
}

func NewBroker() *Broker {
	return &Broker{
		make(chan CellMsg),
		make(chan CellSub),
		make(chan CellSub),
		nil,
	}
}

func (b *Broker) Run() {
	subs := make(map[*Cell]map[*Cell]struct{})
	for {
		select {
		case sub := <-b.subscribe:
			if subs[sub.subcell] == nil {
				subs[sub.subcell] = map[*Cell]struct{}{}
			}
			subs[sub.subcell][sub.thiscell] = struct{}{}
		case unsub := <-b.unsubscribe:
			delete(subs[unsub.subcell], unsub.thiscell)
		case cellmsg := <-b.runevent:
			for c := range subs[cellmsg.from] {
				c.Exec()
			}
		}
	}
}

type CellSub struct {
	subcell  *Cell
	thiscell *Cell
}

type CellMsg struct {
	from *Cell
}
