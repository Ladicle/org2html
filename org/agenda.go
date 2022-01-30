package org

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type AgendaKey string

const (
	AgendaClosed    AgendaKey = "CLOSED"
	AgendaDeadline  AgendaKey = "DEADLINE"
	AgendaScheduled AgendaKey = "SCHEDULED"
)

var _ Node = Agenda{}

// Agenda is a Node to describe Headline agenda such as closed date, deadline and start time.
type Agenda struct {
	Logs map[AgendaKey]Timestamp
}

func (a Agenda) Write(w io.Writer) error {
	return errors.New("not implemented yet")
}

var agendaRegexp = regexp.MustCompile(fmt.Sprintf(
	`(%v|%v|%v):\s(?:\[|<)(\d{4}-\d{2}-\d{2}) ?([A-Za-z]+)? ?(\d{2}:\d{2})? ?(\+\d+[dwmy])?(?:]|>)`,
	AgendaClosed, AgendaDeadline, AgendaScheduled))

func LexAgenda(line string) (Token, bool) {
	// 2nd argument `3` indicates the number of AgendaKeys.
	m := agendaRegexp.FindAllStringSubmatch(line, 3)
	if len(m) > 0 {
		var ret []string
		for _, v := range m {
			ret = append(ret, v[1:]...)
		}
		return NewToken(KindAgenda, len(m), ret), true
	}
	return Token{}, false
}

func ParseAgenda(p *Parser, i int) (int, Node, error) {
	// (0: key, 1: date, 2: week, 3: time, 4: interval)...
	const colNum = 5
	var (
		itemNum = p.tokens[i].num
		agenda  = Agenda{Logs: make(map[AgendaKey]Timestamp, itemNum)}
		vals    = p.tokens[i].vals
	)
	// validate the number of items and values
	if itemNum*colNum != len(vals) {
		return 0, nil, fmt.Errorf("agenda item number and its values are unmatched: num=%v, vals=%#v",
			itemNum, vals)
	}
	// parse each item: if items have the same agendaKey, the latter item overwrites the previous one.
	for j := 0; j < itemNum; j++ {
		var (
			idx      = j * colNum
			key      = AgendaKey(strings.ToUpper(vals[idx]))
			interval = vals[idx+4]
		)
		if vals[idx+3] != "" {
			// time format (YYYY-MM-DD WW HH:MM)
			t, err := ParseTimestamp(strings.Join(vals[idx+1:idx+4], " "), interval)
			if err != nil {
				return 0, nil, err
			}
			agenda.Logs[key] = t
		} else {
			// date format (YYYY-MM-DD WW)
			t, err := ParseDatestamp(strings.Join(vals[idx+1:idx+3], " "), interval)
			if err != nil {
				return 0, nil, nil
			}
			agenda.Logs[key] = t
		}
	}
	return 1, agenda, nil
}
