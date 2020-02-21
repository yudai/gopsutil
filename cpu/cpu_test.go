package cpu

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCpu_times(t *testing.T) {
	v, err := Times(false)
	if err != nil {
		t.Errorf("error %v", err)
	}
	if len(v) == 0 {
		t.Error("could not get CPUs ", err)
	}
	empty := TimesStat{}
	for _, vv := range v {
		if vv == empty {
			t.Errorf("could not get CPU User: %v", vv)
		}
	}

	// test sum of per cpu stats is within margin of error for cpu total stats
	cpuTotal, err := Times(false)
	if err != nil {
		t.Errorf("error %v", err)
	}
	if len(cpuTotal) == 0 {
		t.Error("could not get CPUs ", err)
	}
	perCPU, err := Times(true)
	if err != nil {
		t.Errorf("error %v", err)
	}
	if len(perCPU) == 0 {
		t.Error("could not get CPUs ", err)
	}
	var perCPUUserTimeSum float64
	var perCPUSystemTimeSum float64
	var perCPUIdleTimeSum float64
	for _, pc := range perCPU {
		perCPUUserTimeSum += pc.User
		perCPUSystemTimeSum += pc.System
		perCPUIdleTimeSum += pc.Idle
	}
	margin := 2.0
	assert.InEpsilon(t, cpuTotal[0].User, perCPUUserTimeSum, margin)
	assert.InEpsilon(t, cpuTotal[0].System, perCPUSystemTimeSum, margin)
	assert.InEpsilon(t, cpuTotal[0].Idle, perCPUIdleTimeSum, margin)
}

func TestCpu_counts(t *testing.T) {
	v, err := Counts(true)
	if err != nil {
		t.Errorf("error %v", err)
	}
	if v == 0 {
		t.Errorf("could not get CPU counts: %v", v)
	}
}

func TestCPUTimeStat_String(t *testing.T) {
	v := TimesStat{
		CPU:    "cpu0",
		User:   100.1,
		System: 200.1,
		Idle:   300.1,
	}
	e := `{"cpu":"cpu0","user":100.1,"system":200.1,"idle":300.1,"nice":0.0,"iowait":0.0,"irq":0.0,"softirq":0.0,"steal":0.0,"guest":0.0,"guestNice":0.0}`
	if e != fmt.Sprintf("%v", v) {
		t.Errorf("CPUTimesStat string is invalid: %v", v)
	}
}

func TestCpuInfo(t *testing.T) {
	v, err := Info()
	if err != nil {
		t.Errorf("error %v", err)
	}
	if len(v) == 0 {
		t.Errorf("could not get CPU Info")
	}
	for _, vv := range v {
		if vv.ModelName == "" {
			t.Errorf("could not get CPU Info: %v", vv)
		}
	}
}

func testCPUPercent(t *testing.T, percpu bool) {
	numcpu := runtime.NumCPU()
	testCount := 3

	if runtime.GOOS != "windows" {
		testCount = 100
		v, err := Percent(time.Millisecond, percpu)
		if err != nil {
			t.Errorf("error %v", err)
		}
		// Skip CircleCI which CPU num is different
		if os.Getenv("CIRCLECI") != "true" {
			if (percpu && len(v) != numcpu) || (!percpu && len(v) != 1) {
				t.Fatalf("wrong number of entries from CPUPercent: %v", v)
			}
		}
	}
	for i := 0; i < testCount; i++ {
		duration := time.Duration(10) * time.Microsecond
		v, err := Percent(duration, percpu)
		if err != nil {
			t.Errorf("error %v", err)
		}
		for _, percent := range v {
			// Check for slightly greater then 100% to account for any rounding issues.
			if percent < 0.0 || percent > 100.0001*float64(numcpu) {
				t.Fatalf("CPUPercent value is invalid: %f", percent)
			}
		}
	}
}

func testCPUPercentLastUsed(t *testing.T, percpu bool) {
	numcpu := runtime.NumCPU()
	testCount := 10

	if runtime.GOOS != "windows" {
		testCount = 2
		v, err := Percent(time.Millisecond, percpu)
		if err != nil {
			t.Errorf("error %v", err)
		}
		// Skip CircleCI which CPU num is different
		if os.Getenv("CIRCLECI") != "true" {
			if (percpu && len(v) != numcpu) || (!percpu && len(v) != 1) {
				t.Fatalf("wrong number of entries from CPUPercent: %v", v)
			}
		}
	}
	for i := 0; i < testCount; i++ {
		v, err := Percent(0, percpu)
		if err != nil {
			t.Errorf("error %v", err)
		}
		time.Sleep(1 * time.Millisecond)
		for _, percent := range v {
			// Check for slightly greater then 100% to account for any rounding issues.
			if percent < 0.0 || percent > 100.0001*float64(numcpu) {
				t.Fatalf("CPUPercent value is invalid: %f", percent)
			}
		}
	}

}

func checkPercentStat(t *testing.T, ps PercentStat, numcpu int) {
	t.Helper()

	// Check for slightly greater then 100% to account for any rounding issues.
	if ps.User < 0.0 || ps.User > 100.0001*float64(numcpu) {
		t.Fatalf("CPUPercent value User is invalid: %f", ps.User)
	}
	if ps.System < 0.0 || ps.System > 100.0001*float64(numcpu) {
		t.Fatalf("CPUPercent value System is invalid: %f", ps.System)
	}
	if ps.Idle < 0.0 || ps.Idle > 100.0001*float64(numcpu) {
		t.Fatalf("CPUPercent value Idle is invalid: %f", ps.Idle)
	}
	if ps.Nice < 0.0 || ps.Nice > 100.0001*float64(numcpu) {
		t.Fatalf("CPUPercent value Nice is invalid: %f", ps.Nice)
	}
	if ps.Iowait < 0.0 || ps.Iowait > 100.0001*float64(numcpu) {
		t.Fatalf("CPUPercent value Iowait is invalid: %f", ps.Iowait)
	}
	if ps.Irq < 0.0 || ps.Irq > 100.0001*float64(numcpu) {
		t.Fatalf("CPUPercent value Irq is invalid: %f", ps.Irq)
	}
	if ps.Softirq < 0.0 || ps.Softirq > 100.0001*float64(numcpu) {
		t.Fatalf("CPUPercent value Softirq is invalid: %f", ps.Softirq)
	}
	if ps.Steal < 0.0 || ps.Steal > 100.0001*float64(numcpu) {
		t.Fatalf("CPUPercent value Steal is invalid: %f", ps.Steal)
	}

	total := ps.User + ps.System + ps.Idle + ps.Nice + ps.Iowait +
		ps.Irq + ps.Softirq + ps.Steal
	if math.Round(total) != 100 {
		t.Fatalf("CPUPercent total is invalid: %f", total)
	}

	if ps.Guest < 0.0 || ps.Guest > 100.0001*float64(numcpu) {
		t.Fatalf("CPUPercent value Guest is invalid: %f", ps.Guest)
	}
	if ps.GuestNice < 0.0 || ps.GuestNice > 100.0001*float64(numcpu) {
		t.Fatalf("CPUPercent value GuestNice is invalid: %f", ps.GuestNice)
	}
}

func testCPUItemizedPercent(t *testing.T, percpu bool) {
	numcpu := runtime.NumCPU()
	testCount := 10

	if runtime.GOOS != "windows" {
		testCount = 2
		v, err := ItemizedPercent(time.Millisecond, percpu)
		if err != nil {
			t.Errorf("error %v", err)
		}
		// Skip CircleCI which CPU num is different
		if os.Getenv("CIRCLECI") != "true" {
			if (percpu && len(v) != numcpu) || (!percpu && len(v) != 1) {
				t.Fatalf("wrong number of entries from CPUPercent: %v", v)
			}
		}
	}
	for i := 0; i < testCount; i++ {
		duration := time.Duration(10) * time.Millisecond
		v, err := ItemizedPercent(duration, percpu)
		if err != nil {
			t.Errorf("error %v", err)
		}
		for _, ps := range v {
			checkPercentStat(t, ps, numcpu)
		}
	}
}

func testCPUItemizedPercentLastUsed(t *testing.T, percpu bool) {
	numcpu := runtime.NumCPU()
	testCount := 10

	if runtime.GOOS != "windows" {
		testCount = 2
		v, err := ItemizedPercent(time.Millisecond, percpu)
		if err != nil {
			t.Errorf("error %v", err)
		}
		// Skip CircleCI which CPU num is different
		if os.Getenv("CIRCLECI") != "true" {
			if (percpu && len(v) != numcpu) || (!percpu && len(v) != 1) {
				t.Fatalf("wrong number of entries from CPUPercent: %v", v)
			}
		}
	}
	for i := 0; i < testCount; i++ {
		v, err := ItemizedPercent(0, percpu)
		if err != nil {
			t.Errorf("error %v", err)
		}
		time.Sleep(10 * time.Millisecond)
		for _, ps := range v {
			checkPercentStat(t, ps, numcpu)
		}
	}
}

func TestCPUPercent(t *testing.T) {
	testCPUPercent(t, false)
}

func TestCPUPercentPerCpu(t *testing.T) {
	testCPUPercent(t, true)
}

func TestCPUPercentIntervalZero(t *testing.T) {
	testCPUPercentLastUsed(t, false)
}

func TestCPUPercentIntervalZeroPerCPU(t *testing.T) {
	testCPUPercentLastUsed(t, true)
}

func TestCPUItemizedPercent(t *testing.T) {
	testCPUItemizedPercent(t, false)
}

func TestCPUItemizedPercentPerCpu(t *testing.T) {
	testCPUItemizedPercent(t, true)
}

func TestCPUItemizedPercentIntervalZero(t *testing.T) {
	testCPUItemizedPercentLastUsed(t, false)
}

func TestCPUItemizedPercentIntervalZeroPerCPU(t *testing.T) {
	testCPUItemizedPercentLastUsed(t, true)
}
