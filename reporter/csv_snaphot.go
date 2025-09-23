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

// SnapshotCSV reporter.
//
// Writes a CSV file per step.
type SnapshotCSV struct {
	Observer       observer.Table // Observer to get data from.
	FilePattern    string         // File path and pattern for output files, like out/foo-%06d.csv
	Sep            string         // Column separator. Default ",".
	UpdateInterval int            // Update interval in model ticks.
	Final          bool           // Whether Callback should be called on finalization only, instead of on every tick.
	header         []string
	builder        strings.Builder
	step           int64
}

// Initialize the system
func (s *SnapshotCSV) Initialize(w *ecs.World) {
	s.Observer.Initialize(w)
	s.header = s.Observer.Header()
	if s.UpdateInterval == 0 {
		s.UpdateInterval = 1
	}

	if s.Sep == "" {
		s.Sep = ","
	}

	err := os.MkdirAll(filepath.Dir(fmt.Sprintf(s.FilePattern, 1)), os.ModePerm)
	if err != nil {
		panic(err)
	}

	s.step = 0
}

// Update the system
func (s *SnapshotCSV) Update(w *ecs.World) {
	s.Observer.Update(w)
	if !s.Final && (s.UpdateInterval == 0 || s.step%int64(s.UpdateInterval) == 0) {
		if err := s.writeToFile(w); err != nil {
			panic(err)
		}
	}
	s.step++
}

// Finalize the system
func (s *SnapshotCSV) Finalize(w *ecs.World) {
	if !s.Final {
		return
	}
	if err := s.writeToFile(w); err != nil {
		panic(err)
	}
}

func (s *SnapshotCSV) writeToFile(w *ecs.World) error {
	file, err := os.Create(s.createFileName())
	if err != nil {
		return err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}()

	writer := bufio.NewWriterSize(file, 4096)

	_, err = fmt.Fprintf(writer, "%s\n", strings.Join(s.header, s.Sep))
	if err != nil {
		return err
	}

	values := s.Observer.Values(w)
	s.builder.Reset()
	for _, row := range values {
		for i, v := range row {
			fmt.Fprint(&s.builder, strconv.FormatFloat(v, 'f', -1, 64))
			if i < len(row)-1 {
				fmt.Fprint(&s.builder, s.Sep)
			}
		}
		fmt.Fprint(&s.builder, "\n")
	}
	_, err = fmt.Fprint(writer, s.builder.String())
	if err != nil {
		return err
	}
	return writer.Flush()
}

func (s *SnapshotCSV) createFileName() string {
	var filename string
	if strings.Contains(s.FilePattern, "%") {
		filename = fmt.Sprintf(s.FilePattern, s.step)
	} else {
		filename = s.FilePattern
	}
	return filename
}
