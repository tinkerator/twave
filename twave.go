// Program twave displays a dump.vcd file as an ascii text dump. The
// format is the list of signals at the top and the traces down the
// page.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	path    = flag.String("file", "", "pathname for dump.vcd file")
	base    = flag.String("time", "", "override the time of the start for trace")
	syms    = flag.String("syms", "", "comma separated symbols to list(default to all symbols)")
	waved   = flag.Bool("wavy", false, "output in wavy text format")
	rawTime = flag.Bool("raw", false, "output time in ticks since start")
	start   = flag.Int("start", -1, "start of displayed trace from this tick")
	end     = flag.Int("end", -1, "truncate trace at this tick")
)

func defaultTime(val string) time.Time {
	var t0 time.Time
	var err error
	if t0, err = time.Parse("2006-01-02 15:04:05.999999999 07:00", val); err == nil {
		return t0
	}
	if t0, err = time.Parse("2006-01-02 15:04:05 07:00", val); err == nil {
		return t0
	}
	if t0, err = time.Parse(time.ANSIC, val); err == nil {
		return t0
	}
	log.Fatalf("failed to parse time, %q: %v", val, err)
	return t0
}

type Signal struct {
	Label string
	Bits  int
	Value string
	Alias []string
}

type Valued struct {
	Now   int
	Value string
}

type ParserState struct {
	// What a time unit of 1 in the trace means.
	Timescale       time.Duration
	TimescaleFactor time.Duration
	T0              time.Time
	Scope           []string
	Keys            map[string]int
	Signals         []*Signal
	Symbols         map[string]bool
	Ordered         []string
	LabelMaxLength  int
	Now             int
	Changed         bool
	Waves           map[string][]Valued
	Shortest        int
}

func (s *ParserState) TimeString(n int) string {
	d := time.Duration(n) * s.Timescale
	ans := s.T0.Add(d / time.Duration(s.TimescaleFactor)).Format("2006-01-02 15:04:05.000000000")
	if s.TimescaleFactor == 1 {
		d = 0
	}
	return ans + fmt.Sprintf("%03d", d%1000)
}

func (s *ParserState) Augment(tokens []string) {
	var ts []string
	for _, x := range tokens {
		ts = append(ts, strings.Fields(x)...)
	}
	var err error
	switch ts[0] {
	case "$timescale":
		if strings.HasSuffix(ts[1], "ps") {
			s.Timescale, err = time.ParseDuration(ts[1][:len(ts[1])-2] + "ns")
			s.TimescaleFactor = 1000
		} else {
			s.Timescale, err = time.ParseDuration(ts[1])
			s.TimescaleFactor = 1
		}
		if err != nil {
			log.Fatalf("unable to parse duration %q: %v", ts[1], err)
		}
	case "$date":
		tString := strings.Join(ts[1:len(ts)-1], " ")
		s.T0 = defaultTime(tString)
	case "$scope":
		s.Scope = append(s.Scope, ts[2])
	case "$upscope":
		s.Scope = s.Scope[:len(s.Scope)-1]
	case "$var":
		if s.Keys == nil {
			s.LabelMaxLength = 32
			s.Keys = make(map[string]int)
		}
		if ts[1] == "parameter" {
			return
		}
		label := fmt.Sprintf("%s.%s", strings.Join(s.Scope, "."), ts[4])
		bits := 1
		if n, err := strconv.Atoi(ts[2]); err == nil {
			bits = n
		}
		if bits != 1 {
			label = fmt.Sprint(label, ts[5])
		}
		if l := len(label); l > s.LabelMaxLength {
			s.LabelMaxLength = l
		}

		old, present := s.Keys[ts[3]]
		if present {
			s.Signals[old].Alias = append(s.Signals[old].Alias, label)
		} else {
			s.Keys[ts[3]] = len(s.Signals)
			s.Signals = append(s.Signals, &Signal{
				Label: label,
				Bits:  bits,
				Value: "x",
			})
		}
	default:
		if !*waved {
			fmt.Println(s.Scope, ":", ts)
		}
	}
}

// DumpStateNow displays the state at a single timestamp.
func (s *ParserState) DumpStateNow() {
	if *start != -1 && s.Now < *start {
		return
	}
	if *end != -1 && s.Now > *end {
		return
	}
	if !*waved {
		if *rawTime {
			fmt.Printf(fmt.Sprintf("%%%dd", s.LabelMaxLength), s.Now)
		} else {
			fmt.Printf(fmt.Sprintf("%%%ds", s.LabelMaxLength), s.TimeString(s.Now))
		}
		for _, sym := range s.Ordered {
			var c *Signal
			for _, c = range s.Signals {
				if c.Label == sym {
					break
				}
			}
			value := c.Value
			if c.Bits != 1 && len(value) != c.Bits {
				if value == "x" || value == "z" {
					value = strings.Repeat(value, c.Bits)
				} else {
					value = strings.Repeat("0", c.Bits-len(value)) + value
				}
			}
			fmt.Printf(" %s", value)
		}
		fmt.Println()
	} else {
		if s.Waves == nil {
			s.Waves = make(map[string][]Valued)
		}
		for _, c := range s.Signals {
			if s.Symbols != nil && !s.Symbols[c.Label] {
				continue
			}
			ar := s.Waves[c.Label]
			if n := len(ar); n > 0 {
				// drop duplicates
				if ar[n-1].Value == c.Value {
					continue
				}
			}
			s.Waves[c.Label] = append(ar, Valued{
				Now:   s.Now,
				Value: c.Value,
			})
		}
	}
}

// Datum consumes a VVD data item and dumps a value if the tokens start with #
func (s *ParserState) Datum(tokens []string) {
	switch tokens[0][0] {
	case '$':
		return
	case '#':
		was := s.Now
		if s.Changed {
			s.DumpStateNow()
		}
		s.Changed = true
		s.Now, _ = strconv.Atoi(tokens[0][1:])
		if delta := s.Now - was; s.Shortest == 0 || was < s.Shortest {
			s.Shortest = delta
		}
		return
	}
	if len(tokens) == 2 {
		key := tokens[1]
		val := tokens[0][1:]
		i, ok := s.Keys[key]
		if !ok {
			return
		}
		s.Signals[i].Value = val
	} else if len(tokens) == 1 {
		key := tokens[0][1:]
		val := tokens[0][:1]
		i, ok := s.Keys[key]
		if !ok {
			return
		}
		s.Signals[i].Value = val
	}
}

// Legend dumps all of the symbol names in a key preface for the text
// dump.
func (s *ParserState) Legend() {
	if *waved {
		return
	}
	var kept []int
	for i, sym := range s.Ordered {
		var c *Signal
		for _, c = range s.Signals {
			if c.Label == sym {
				break
			}
		}
		kept = append(kept, c.Bits)
		fmt.Printf(fmt.Sprintf("%%%ds", s.LabelMaxLength), sym)
		for j := 0; j <= i; j++ {
			ch := "+"
			if j != i {
				ch = "|"
			}
			fmt.Print(strings.Repeat("-", kept[j]), ch)
		}
		fmt.Println()
	}
	fmt.Print(strings.Repeat(" ", s.LabelMaxLength))
	for _, c := range s.Signals {
		if s.Symbols != nil && !s.Symbols[c.Label] {
			continue
		}
		ch := "|"
		fmt.Print(strings.Repeat(" ", c.Bits), ch)
	}
	fmt.Println()
}

func main() {
	flag.Parse()

	if *path == "" {
		log.Fatal("mandatory argument missing: --file=<pathname>")
	}

	state := &ParserState{}
	if *base != "" {
		state.T0 = defaultTime(*base)
	}
	if *syms != "" {
		state.Symbols = make(map[string]bool)
		for _, s := range strings.Split(*syms, ",") {
			state.Symbols[s] = true
			state.Ordered = append(state.Ordered, s)
		}
	}

	file, err := os.Open(*path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var compound []string
	initialized := false
	for scanner.Scan() {
		tokens := scanner.Text()
		if !initialized {
			compound = append(compound, tokens)
			if strings.HasSuffix(tokens, "$end") {
				initialized = strings.HasPrefix(compound[0], "$enddefinitions")
				if !initialized {
					state.Augment(compound)
				} else {
					if state.Symbols == nil {
						state.Symbols = make(map[string]bool)
						for _, sig := range state.Signals {
							state.Symbols[sig.Label] = true
							state.Ordered = append(state.Ordered, sig.Label)
						}
						sort.Strings(state.Ordered)
					}
					state.Legend()
				}
				compound = nil
			}
		} else {
			state.Datum(strings.Fields(tokens))
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if *waved {
		for _, s := range state.Ordered {
			var sig *Signal
			for _, sig = range state.Signals {
				if sig.Label == s {
					break
				}
			}
			vs, ok := state.Waves[s]
			if !ok {
				continue
			}
			var chunks []string
			var last string
			for j, i := 0, 0; i < state.Now; i += state.Shortest {
				if i < *start && *start != -1 {
					continue
				}
				if i > *end && *end != -1 {
					continue
				}
				if vs[j].Now < i && j < len(vs)-1 && vs[j+1].Now == i {
					j++
				}
				v := vs[j].Value
				if v == "z" || v == "x" {
					fmt.Print(v)
				} else if sig.Bits == 1 {
					switch v {
					case "0":
						fmt.Print("_")
					case "1":
						fmt.Print("^")
					default:
						fmt.Print(v)
					}
					last = v
				} else {
					n, ok := new(big.Int).SetString(v, 2)
					if !ok {
						log.Fatalf("unable to parse %q", v)
					}
					t := n.Text(16)
					if (j+1 < len(vs) && i+state.Shortest >= vs[j+1].Now) || (i+state.Shortest > *end && *end != -1) {
						chunks = append(chunks, t)
						fmt.Print(">")
					} else if last == t {
						fmt.Print("-")
					} else {
						fmt.Print("<")
					}
					last = t
				}
			}
			if len(chunks) != 0 {
				fmt.Printf(" %s %s\n", s, strings.Join(chunks, ","))
			} else {
				fmt.Printf(" %s\n", s)
			}
		}
	} else if *rawTime {
		fmt.Printf("minimum period: %d\n", state.Shortest)
	}
}
