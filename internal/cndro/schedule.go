package cndro

import (
	"bufio"
	_ "embed"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
)

//go:embed day1_schedule
var day1Data string

//go:embed day2_schedule
var day2Data string

//go:embed day1_workshop_schedule
var day1WorkshopData string

//go:embed day2_workshop_schedule
var day2WorkshopData string

// scheduleEntry is one row: start/end clock, talk title, room, optional speaker.
type scheduleEntry struct {
	startH, startM int
	endH, endM     int // -1 when only HH MM (no end) was given
	title          string
	room           string
	speaker        string
	workshop       bool
}

// parseScheduleLine parses:
//   - Optional speaker suffix: " ... @ room # Speaker Name"
//   - "HH MM EH EM title @ room"
//   - "HH MM title @ room"
func parseScheduleLine(line string) (scheduleEntry, bool) {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "#") {
		return scheduleEntry{}, false
	}
	var speaker string
	if i := strings.LastIndex(line, " # "); i >= 0 {
		speaker = strings.TrimSpace(line[i+len(" # "):])
		line = strings.TrimSpace(line[:i])
	}
	parts := strings.SplitN(line, " @ ", 2)
	if len(parts) != 2 {
		return scheduleEntry{}, false
	}
	left := strings.Fields(parts[0])
	room := strings.TrimSpace(parts[1])
	if len(left) < 3 {
		return scheduleEntry{}, false
	}
	sh, err1 := strconv.Atoi(left[0])
	sm, err2 := strconv.Atoi(left[1])
	if err1 != nil || err2 != nil {
		return scheduleEntry{}, false
	}
	eh, err3 := strconv.Atoi(left[2])
	em, err4 := strconv.Atoi(left[3])
	if err3 == nil && err4 == nil && len(left) >= 5 {
		title := strings.Join(left[4:], " ")
		return scheduleEntry{
			startH: sh, startM: sm,
			endH: eh, endM: em,
			title: title, room: room, speaker: speaker,
		}, true
	}
	title := strings.Join(left[2:], " ")
	return scheduleEntry{
		startH: sh, startM: sm,
		endH: -1, endM: -1,
		title: title, room: room, speaker: speaker,
	}, true
}

func formatClock(h, m int) string {
	return fmt.Sprintf("%02d:%02d", h, m)
}

func formatTimeRange(e scheduleEntry) string {
	if e.endH < 0 {
		return formatClock(e.startH, e.startM)
	}
	return formatClock(e.startH, e.startM) + "–" + formatClock(e.endH, e.endM)
}

func parseScheduleData(data string) []scheduleEntry {
	var out []scheduleEntry
	sc := bufio.NewScanner(strings.NewReader(data))
	for sc.Scan() {
		e, ok := parseScheduleLine(sc.Text())
		if ok {
			out = append(out, e)
		}
	}
	return out
}

func parseWorkshopSchedule(data string) []scheduleEntry {
	var out []scheduleEntry
	sc := bufio.NewScanner(strings.NewReader(data))
	for sc.Scan() {
		e, ok := parseScheduleLine(sc.Text())
		if ok {
			e.workshop = true
			out = append(out, e)
		}
	}
	return out
}

func endSortKey(e scheduleEntry) (int, int) {
	if e.endH < 0 {
		return e.startH, e.startM
	}
	return e.endH, e.endM
}

// mergeScheduleEntries concatenates a and b, then sorts by start time, end time (shorter first),
// then sessions before workshops, then room.
func mergeScheduleEntries(a, b []scheduleEntry) []scheduleEntry {
	out := make([]scheduleEntry, 0, len(a)+len(b))
	out = append(out, a...)
	out = append(out, b...)
	sort.SliceStable(out, func(i, j int) bool {
		ai, aj := out[i], out[j]
		if ai.startH != aj.startH {
			return ai.startH < aj.startH
		}
		if ai.startM != aj.startM {
			return ai.startM < aj.startM
		}
		ehi, emi := endSortKey(ai)
		ehj, emj := endSortKey(aj)
		if ehi != ehj {
			return ehi < ehj
		}
		if emi != emj {
			return emi < emj
		}
		if ai.workshop != aj.workshop {
			return !ai.workshop
		}
		return ai.room < aj.room
	})
	return out
}

func writeScheduleSection(w *tabwriter.Writer, entries []scheduleEntry) {
	fmt.Fprintln(w, "⏱️ TIME\t ⚒️ TYPE\t 💡TITLE\t☸️ SPEAKER\t🏢ROOM")
	prevRange := ""
	for _, e := range entries {
		tr := formatTimeRange(e)
		timeCol := tr
		if tr == prevRange {
			timeCol = ""
		} else {
			prevRange = tr
		}
		typ := "Session"
		if e.workshop {
			typ = "Workshop"
		}
		r := e.room
		if r == "" || r == "-" {
			r = "—"
		}
		sp := e.speaker
		if sp == "" {
			sp = "—"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", timeCol, typ, e.title, sp, r)
	}
}

// WriteDay1 prints the Day 1 schedule table to w (sessions and workshops merged).
func WriteDay1(w io.Writer) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	writeScheduleSection(tw, mergeScheduleEntries(parseScheduleData(day1Data), parseWorkshopSchedule(day1WorkshopData)))
	return tw.Flush()
}

// WriteDay2 prints the Day 2 schedule table to w (sessions and workshops merged).
func WriteDay2(w io.Writer) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	writeScheduleSection(tw, mergeScheduleEntries(parseScheduleData(day2Data), parseWorkshopSchedule(day2WorkshopData)))
	return tw.Flush()
}
