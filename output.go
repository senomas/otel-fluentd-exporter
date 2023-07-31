package fluentdexporter

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fluent/fluent-logger-golang/fluent"
)

type Output interface {
	Post(time *time.Time, tag string, data map[string]interface{})

	Close()
}

type Fluentd struct {
	Host   string
	Port   int
	Tag    string
	clog   chan map[string]interface{}
	fluent *fluent.Fluent
}

func NewOutputFluentd(host string, port int, tag string, bufferSize int) Output {
	fluent, err := fluent.New(fluent.Config{FluentHost: host, FluentPort: port})
	if err != nil {
		panic(err)
	}
	clog := make(chan map[string]interface{}, bufferSize)
	go func() {
		for message := range clog {
			err := fluent.Post(tag, message)
			if err != nil {
				fmt.Printf("ERROR %s: %v+\n", err, message)
			}
		}
	}()
	return &Fluentd{
		Host:   host,
		Port:   port,
		Tag:    tag,
		clog:   clog,
		fluent: fluent,
	}
}

func (out *Fluentd) Post(time *time.Time, mtype string, message map[string]interface{}) {
	message["@timestamp"] = time.Format("2006-01-02T15:04:05.000Z07:00")
	message["type"] = mtype
	out.clog <- message
}

func (out *Fluentd) Close() {
	out.fluent.Close()
}

type OutputFile struct {
	Path       string
	FileFormat string
	UID        int
	GID        int
	flog       *os.File
	tnext      *time.Time
}

func NewOutputFile(Path string, FileFormat string, UID int, GID int) Output {
	return &OutputFile{
		Path:       Path,
		FileFormat: FileFormat,
		UID:        UID,
		GID:        GID,
	}
}

func (out *OutputFile) Post(ctime *time.Time, mtype string, message map[string]interface{}) {
	tz := time.Now()
	if out.tnext == nil || tz.After(*out.tnext) {
		flogFlag := os.O_APPEND | os.O_CREATE | os.O_WRONLY
		tnext := tz.Truncate(time.Hour).Add(time.Hour)
		fn := fmt.Sprintf("%s-%s.log", out.Path, tz.Format(out.FileFormat))
		flog, err := os.OpenFile(fn, flogFlag, 0644)
		if err != nil {
			panic(err)
		}
		flog.Chown(out.UID, out.GID)
		out.flog = flog
		out.tnext = &tnext
	}
	fmt.Fprintf(out.flog, "{\"@timestamp\":\"%s\"", ctime.Format("2006-01-02T15:04:05.000Z07:00"))
	message["type"] = mtype
	for k, v := range message {
		vstr, _ := json.Marshal(v)
		fmt.Fprintf(out.flog, ",\"%s\":%s", k, vstr)
	}
	fmt.Fprintln(out.flog, "}")
}

func (out *OutputFile) Close() {
	if out.flog != nil {
		out.flog.Close()
	}
}

type OutputStdOut struct {
}

func NewOutputStdOut() Output {
	return &OutputStdOut{}
}

func (out *OutputStdOut) Post(ctime *time.Time, mtype string, message map[string]interface{}) {
	fmt.Printf("%s ", ctime.Format("2006-01-02T15:04:05.000Z07:00"))
	message["type"] = mtype
	for k, v := range message {
		vstr, _ := json.Marshal(v)
		fmt.Printf(" %s:%s", k, vstr)
	}
}

func (out *OutputStdOut) Close() {
}

type DummyOutputData struct {
	Time    time.Time
	Tag     string
	Message map[string]interface{}
}

type DummyOutput struct {
	data []*DummyOutputData
}

func (out *DummyOutput) Post(ctime *time.Time, mtype string, message map[string]interface{}) {
	// fmt.Printf("LOG [%s]\n", message)
	message["type"] = mtype
	out.data = append(out.data, &DummyOutputData{Time: *ctime, Tag: mtype, Message: message})
}

func (out *DummyOutput) Close() {
}

func (out *DummyOutput) String() string {
	res := []string{}
	for _, v := range out.data {
		str, _ := json.MarshalIndent(v.Message, "", "  ")
		res = append(res, string(str))
	}
	return strings.Join(res, "\n")
}

func (out *DummyOutput) Field(index int, field string) string {
	if index < len(out.data) {
		f := out.data[index].Message[field]
		if f != nil {
			if v, ok := f.(string); ok {
				return v
			}
			str, _ := json.MarshalIndent(f, "", "  ")
			return string(str)
		}
	}
	return ""
}
