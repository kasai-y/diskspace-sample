package main

import (
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

func main() {

	fmt.Println(runtime.GOOS)

	if err := action(); err != nil {
		fmt.Println(fmt.Sprintf("%+v", err))
	}
}

func action() error {
	flag.Parse()
	path := flag.Arg(0)
	if len(path) == 0 {
		return nil
	}

	var ds *diskSpace
	switch runtime.GOOS {
	case "windows":
		d, err := diskUsageWin(path)
		if err != nil {
			return errors.WithStack(err)
		}
		ds = &d
	default:
		d, err := diskUsage(path)
		if err != nil {
			return errors.WithStack(err)
		}
		ds = &d
	}

	fmt.Println(fmt.Sprintf("All  : %d", ds.All))
	fmt.Println(fmt.Sprintf("Used : %d", ds.Used))
	fmt.Println(fmt.Sprintf("Free : %d", ds.Free))
	return nil
}

type diskSpace struct {
	All  int
	Used int
	Free int
}

func diskUsage(path string) (disk diskSpace, err error) {
	b, err := exec.Command("df", path).Output()
	if err != nil {
		return diskSpace{}, errors.WithStack(err)
	}

	// (Example)
	// Filesystem 512-blocks      Used Available Capacity iused      ifree %iused  Mounted on
	// /dev/disk1  234573824 219903288  14158536    94% 2469867 4292497412    0%   /
	str := string(b)
	v := regexp.MustCompile(" +").Split(strings.Split(str, "\n")[1], -1)

	fmt.Println(str)

	disk.All, _ = strconv.Atoi(v[1])
	disk.Free, _ = strconv.Atoi(v[5])
	disk.Used, _ = strconv.Atoi(v[2])
	return
}

func diskUsageWin(path string) (disk diskSpace, err error) {
	b, err := exec.Command("fsutil", "volume", "diskfree", path).Output()
	if err != nil {
		return diskSpace{}, errors.WithStack(err)
	}

	// (Example)
	// 空きバイト総数           : 1111111111111
	// バイト総数              : 3333333333333
	// 利用可能な空きバイト総数  : 2222222222222
	str := string(b)
	l := strings.Split(str, "\n")

	fmt.Println(str)

	disk.All, _ = strconv.Atoi(regexp.MustCompile(`[\s\t:]+`).Split(l[1], -1)[1])
	disk.Free, _ = strconv.Atoi(regexp.MustCompile(`[\s\t:]+`).Split(l[0], -1)[1])
	disk.Used = disk.All - disk.Free
	return
}
