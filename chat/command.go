package chat

import (
	"fmt"
	"lin/analysis"
	"strconv"
	"time"
)

// command 命令格式
type command interface { //参数形式
	name() string
	run(param ...string) string
}

var commands = []command{
	&help{},
	&popular{},
	&stats{},
}

func init() {

}

// help
type help struct{}

func (c *help) name() string {
	return "/help"
}

func (c *help) run(param ...string) string {
	out := ""
	for _, v := range commands {
		if v.name() != c.name() {
			out += v.name() + "\r\n"
		}
	}
	return out
}

// popular
type popular struct{}

func (c *popular) name() string {
	return "/popular"
}

func (c *popular) run(param ...string) string {
	out := ""
	for _, v := range analysis.GetWorldsFrequencyStatistics() {
		out += v + "\r\n"
	}
	return out
}

// stats
type stats struct{}

func (c *stats) name() string {
	return "/stats"
}

func formatOutTime(dut int64) string {
	d := dut / int64(time.Hour*24/time.Second)
	h := dut/int64(time.Hour/time.Second) - d*24
	m := dut/int64(time.Minute/time.Second) - h*60 - d*24*60
	s := dut - m*60 - h*60*60 - d*24*60*60
	out := fmt.Sprintf("%02dd %02dh %02dm %02ds", d, h, m, s)
	return out
}

func (c *stats) run(param ...string) string {
	if len(param) != 1 {
		return "param error, param len != 1"
	}
	uid, err := strconv.Atoi(param[0])
	if err != nil {
		return err.Error()
	}
	u, ok := WorldRoom.users[UUID(uid)]
	if !ok {
		return "user not found"
	}
	return formatOutTime(time.Now().Unix() - u.loginTIme)
}
