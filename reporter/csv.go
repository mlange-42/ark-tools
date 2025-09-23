package reporter

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mlange-42/ark-tools/observer"
	"github.com/mlange-42/ark/ecs"
)

// CSV reporter.
//
// Writes one row to a CSV file per step.
type CSV struct {
	Observer       observer.Row // Observer to get data from.
	File           string       // Path to the output file.
	Sep            string       // Column separator. Default ",".
	UpdateInterval int          // Update interval in model ticks.
	Final          bool         // Whether Callback should be called on finalization only, instead of on every tick.
	file           *os.File
	writer         *bufio.Writer
	header         []string
	builder        strings.Builder
	step           int64
}

// Initialize the system
func (s *CSV) Initialize(w *ecs.World) {
	s.Observer.Initialize(w)
	s.header = s.Observer.Header()
	if s.UpdateInterval == 0 {
		s.UpdateInterval = 1
	}
	if s.Sep == "" {
		s.Sep = ","
	}

	err := os.MkdirAll(filepath.Dir(s.File), os.ModePerm)
	if err != nil {
		panic(err)
	}

	s.file, err = os.Create(s.File)
	if err != nil {
		panic(err)
	}
	s.writer = bufio.NewWriterSize(s.file, 4096)
	_, err = fmt.Fprintf(s.writer, "t%s%s\n", s.Sep, strings.Join(s.header, s.Sep))
	if err != nil {
		panic(err)
	}

	s.step = 0
}

// Update the system
func (s *CSV) Update(w *ecs.World) {
	s.Observer.Update(w)
	if !s.Final && (s.UpdateInterval == 0 || s.step%int64(s.UpdateInterval) == 0) {
		if err := s.writeToFile(w); err != nil {
			panic(err)
		}
	}
	s.step++
}

// Finalize the system
func (s *CSV) Finalize(w *ecs.World) {
	if s.Final {
		if err := s.writeToFile(w); err != nil {
			panic(err)
		}
	}
	if err := s.writer.Flush(); err != nil {
		panic(err)
	}
	if err := s.file.Close(); err != nil {
		panic(err)
	}
}

func (s *CSV) writeToFile(w *ecs.World) error {
	values := s.Observer.Values(w)
	s.builder.Reset()
	fmt.Fprintf(&s.builder, "%d%s", s.step, s.Sep)
	for i, v := range values {
		fmt.Fprint(&s.builder, strconv.FormatFloat(v, 'f', -1, 64))
		if i < len(values)-1 {
			fmt.Fprint(&s.builder, s.Sep)
		}
	}
	_, err := fmt.Fprintf(s.writer, "%s\n", s.builder.String())
	if err != nil {
		return err
	}
	return nil
}
