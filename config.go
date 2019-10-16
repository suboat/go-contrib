package contrib

import (
	"github.com/suboat/go-contrib/log"
	"github.com/tudyzhb/yaml"

	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"
)

// 配置文件
type Config struct {
	sync.RWMutex
	pathYaml   *string                 // yaml配置文件保存地址
	savePoint  interface{}             //
	hookChange func(interface{}) error //
	comments   map[string]string       // 配置文件的备注信息
	silent     bool                    // true: 不打印日志
}

// toJson
func (c *Config) ToJson() (ret string, err error) {
	if c.savePoint == nil {
		err = ErrUndefined
		return
	}
	var b []byte
	if b, err = json.Marshal(c.savePoint); err != nil {
		return
	} else {
		ret = string(b)
	}
	return
}

// fromJson
func (c *Config) FromJson(content string) (ret string, err error) {
	c.Lock()
	defer c.Unlock()
	if c.savePoint == nil {
		err = ErrUndefined
		return
	}
	if err = json.Unmarshal([]byte(content), c.savePoint); err != nil {
		return
	}
	if ret, err = c.ToJson(); err != nil {
		return
	}
	return
}

// 不打印日志
func (c *Config) SetSilent(v bool) {
	c.silent = v
}

// 保存配置
func (c *Config) Save(newConfig interface{}) (err error) {
	c.Lock()
	defer c.Unlock()
	var (
		savePath    string
		readContent []byte
		saveContent []byte
	)
	if savePath, err = c.GetSavePath(); err != nil {
		return
	}
	// 读旧记录
	readContent, _ = ioutil.ReadFile(savePath)

	if newConfig == nil {
		newConfig = c.savePoint
	}
	if saveContent, err = yaml.MarshalWithComments(newConfig, c.comments); err != nil {
		return
	} else if bytes.Equal(readContent, saveContent) == true {
		// 不重复保存
		return
	}

	// 写入记录
	if err = ioutil.WriteFile(savePath, saveContent, 0666); err != nil {
		return
	}
	if c.silent == false {
		log.Debug(fmt.Sprintf("[config] save config to %s bytes:%d->%d",
			savePath, len(readContent), len(saveContent)))
	}

	// hook
	if c.hookChange != nil {
		if _err := c.hookChange(newConfig); _err != nil {
			log.Warn(fmt.Sprintf(`[config] hookChange error: %v`, _err))
		}
	}
	return
}

// 取保存地址
func (c *Config) GetSavePath() (ret string, err error) {
	if c.pathYaml == nil {
		err = ErrUndefined
		return
	} else {
		ret = *c.pathYaml
	}
	return
}

// 设置保存地址
func (c *Config) SetSavePath(savePath string) (err error) {
	if len(savePath) == 0 {
		c.pathYaml = nil
		err = ErrUndefined
		return
	} else {
		c.pathYaml = &savePath
	}
	return
}

// 设置要保存的对象
func (c *Config) SetSavePoint(saveTarget interface{}) (err error) {
	c.savePoint = saveTarget
	return
}

// 设置要保存的对象
func (c *Config) SetHookChange(f func(interface{}) error) (err error) {
	c.hookChange = f
	return
}

// 设置备注信息
func (c *Config) SetComments(comments map[string]string) (err error) {
	comments_ := make(map[string]string)
	for key := range comments {
		comments_[strings.ToLower(key)] = comments[key]
	}
	c.comments = comments_
	return
}
