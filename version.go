package contrib

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

var (
	// 版本号识别
	regVersion   = regexp.MustCompile(`(\w+)?-?v(\d+)\.(\d+)\.(\d+)\(?(\w+)?\)?-?(\w{8})?`)
	regGitCommit = regexp.MustCompile(`v?((\d+)\.(\d+)\.(\d+))?(-\d+)?(-?g?([0-9a-f]+))?(-dirty)?`) // v0.0.1-1-gd4f800c-dirty
)

// 版本信息
type Version struct {
	Major int // 主版本号:向前不兼容	变化通常意味着模块的巨大的变化
	Minor int // 次版本号:对同级兼容	通常只反映了一些较大的更改, 比如模块的API增加等
	Patch int // 补丁版本:对同级兼容	通常情况下如果只是对模块的修改而不影响API接口
	// optional
	Model  string     // 模块名称
	Hash   string     // 可执行文件哈希
	Commit *string    // 代码提交哈希
	Build  *time.Time // 模块编译时间
}

// 解析版本号
func (v *Version) ParseStr(s string) (err error) {
	if v == nil {
		err = ErrUndefined
		return
	}
	var (
		vals = regVersion.FindStringSubmatch(s)
	)
	v.Model = vals[1]
	if _v, _err := strconv.ParseUint(vals[2], 10, 32); _err != nil {
		err = ErrParamInvalid
		return
	} else {
		v.Major = int(_v)
	}
	if _v, _err := strconv.ParseUint(vals[3], 10, 32); _err != nil {
		err = ErrParamInvalid
		return
	} else {
		v.Minor = int(_v)
	}
	if _v, _err := strconv.ParseUint(vals[4], 10, 32); _err != nil {
		err = ErrParamInvalid
		return
	} else {
		v.Patch = int(_v)
	}
	if len(vals[5]) > 0 {
		// 优先解析时间
		if _v, _err := time.Parse("01021504", vals[5]); _err == nil {
			v.Build = &_v
		} else {
			_s := vals[5]
			v.Commit = &_s
		}
	}
	if len(vals[6]) > 0 {
		v.Hash = vals[6]
	}
	return
}

// ParseCommit 从gitCommit解析 v0.0.1-1-gd4f800c-dirty
func (v *Version) ParseCommit(gitCommit string) (err error) {
	if len(gitCommit) == 0 {
		return
	}
	var (
		vals   = regGitCommit.FindStringSubmatch(gitCommit)
		ver    = vals[1]
		commit = vals[7]
	)
	if len(ver) == 0 && len(commit) == 0 {
		return
	}

	//
	if len(ver) > 0 {
		//
		if _v, _err := strconv.ParseUint(vals[2], 10, 32); _err == nil {
			v.Major = int(_v)
		} else {
			err = _err
			return
		}
		//
		if _v, _err := strconv.ParseUint(vals[3], 10, 32); _err == nil {
			v.Minor = int(_v)
		} else {
			err = _err
			return
		}
		//
		if _v, _err := strconv.ParseUint(vals[4], 10, 32); _err == nil {
			v.Patch = int(_v)
		} else {
			err = _err
			return
		}
	}

	if len(commit) > 0 {
		if vals[8] == "-dirty" {
			commit += "-dirty"
		}
		v.Commit = &commit
	} else {
		v.Commit = nil
	}
	return
}

//
func (v *Version) ParseInt(i_ int32) (err error) {
	if v == nil {
		err = ErrUndefined
		return
	}
	i := int(i_)
	v.Major = i / 1000000
	i = i - (v.Major * 1000000)
	v.Minor = i / 1000
	i -= v.Minor * 1000
	v.Patch = i
	return
}

// 取运行文件的哈希前32位
func (v *Version) GetRunFileHash() (ret string) {
	if len(os.Args) > 0 {
		if p, err := filepath.Abs(os.Args[0]); err == nil {
			if f, err := os.Open(p); err == nil {
				defer f.Close()
				h := sha1.New()
				if _, err = io.Copy(h, f); err == nil {
					ret = fmt.Sprintf("%x", h.Sum(nil))
				}
			}
		}
	}
	if len(ret) > 8 {
		ret = ret[0:8]
	}
	return
}

// 取运行文件的修改时间
func (v *Version) GetRunFileTime() (ret *time.Time) {
	if len(os.Args) > 0 {
		if p, err := filepath.Abs(os.Args[0]); err == nil {
			if info, err := os.Stat(p); err == nil {
				if t := info.ModTime(); t.Unix() > 0 {
					ret = &t
				}
			}
		}
	}
	return
}

//
func (v *Version) String() (ret string) {
	if v == nil {
		return "v0.0.0"
	}
	if len(v.Model) > 0 {
		ret = v.Model + "-"
	}
	if len(v.Hash) == 0 {
		// 尝试填补Hash
		v.Hash = v.GetRunFileHash()
	}
	if v.Commit == nil && v.Build == nil {
		// 尝试使用修改时间来表示编译时间
		v.Build = v.GetRunFileTime()
	}
	//
	ret = fmt.Sprintf(`%sv%d.%d.%d`, ret, v.Major, v.Minor, v.Patch)
	// 优先输出编译版本
	if v.Commit != nil {
		ret = fmt.Sprintf(`%s(%s)`, ret, *v.Commit)
	} else if v.Build != nil {
		ret = fmt.Sprintf(`%s(%s)`, ret, v.Build.Format("01021504"))
	}
	if len(v.Hash) > 0 {
		ret = fmt.Sprintf(`%s-%s`, ret, v.Hash)
	}
	return
}

//
func (v *Version) Int32() (ret int32) {
	if v == nil {
		return
	}
	ret = int32(v.Major)*1000000 + int32(v.Minor)*1000 + int32(v.Patch)
	return
}
