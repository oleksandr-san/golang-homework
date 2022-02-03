package logger

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type mockPrinter struct {
	printCalls  []interface{}
	printfCalls []interface{}
}

func (p *mockPrinter) Println(v ...interface{}) {
	p.printCalls = append(p.printCalls, v)
}

func (p *mockPrinter) Printf(fmt string, v ...interface{}) {
	args := []interface{}{fmt}
	args = append(args, v...)
	p.printfCalls = append(p.printfCalls, args)
}

func TestLogger(t *testing.T) {
	t.Run("INFO log level is considered", func(t *testing.T) {
		m := mockPrinter{}
		l := Logger{printer: &m, level: INFO}

		l.Debug("debug")
		l.Info("info")
		l.Warning("warning")
		l.Error("error")

		require.Equal(t, 0, len(m.printfCalls))
		require.Equal(t, []interface{}{
			[]interface{}{"info"},
			[]interface{}{"warning"},
			[]interface{}{"error"},
		},
			m.printCalls)
	})

	t.Run("ERROR logf level is considered", func(t *testing.T) {
		m := mockPrinter{}
		l := Logger{printer: &m, level: ERROR}

		l.Debugf("debug: %s", "A")
		l.Infof("info: %s", "A")
		l.Warningf("warning: %s", "A")
		l.Errorf("error: %s", "A")

		require.Equal(t, 0, len(m.printCalls))
		require.Equal(t, []interface{}{
			[]interface{}{"error: %s", "A"},
		}, m.printfCalls)
	})

	t.Run("INFO is default log level", func(t *testing.T) {
		l, err := New("", "")

		require.Nil(t, err)
		require.Nil(t, l.closer)
		require.Equal(t, l.level, INFO)
		require.Equal(t, fmt.Sprintf("%T", l.printer), "*log.Logger")
	})

	t.Run("INFO is default log level", func(t *testing.T) {
		l, err := New("", "")

		require.NoError(t, err)
		require.Nil(t, l.closer)
		require.Equal(t, l.level, INFO)
		require.Equal(t, fmt.Sprintf("%T", l.printer), "*log.Logger")
	})

	t.Run("Log file is used when filePath specified", func(t *testing.T) {
		file, err := ioutil.TempFile("", "logs")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(file.Name())

		l, err := New("WARNING", file.Name())
		require.NoError(t, err)
		require.NotNil(t, l.closer)

		defer l.closer.Close()
		require.Equal(t, l.level, WARNING)
		require.Equal(t, fmt.Sprintf("%T", l.printer), "*log.Logger")
	})
}
