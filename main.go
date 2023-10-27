package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/go-redis/redis"
)

func main() {
	var (
		configFile string
	)

	defaul := "config.toml"
	//defaul := "D:\\go\\redisbatchcmd\\config.toml"
	flag.StringVar(&configFile, "c", defaul, "config file")
	flag.Parse()

	if !IsExist(configFile) { // 如果存在，加载配置文件，不存在则跳过配置文件，使用默认值
		err := errors.Errorf("not found config file:%s", configFile)
		fmt.Printf("Connect Failed! Err: %v\n", err)
		return
	}
	err := InitConfig(configFile)
	if err != nil {
		fmt.Printf("Connect Failed! Err: %v\n", err)
		return
	}

	redisAddr := fmt.Sprintf("%s:%d", _config.Host, _config.Port)
	db := redis.NewClient(&redis.Options{
		Addr:        redisAddr,        // redis服务ip:port
		Password:    _config.Password, // redis的认证密码
		DB:          _config.Database, // 连接的database库
		IdleTimeout: 300,              // 默认Idle超时时间
		PoolSize:    100,              // 连接池
	})
	_, err = db.Ping().Result()
	if err != nil {
		fmt.Printf("Connect Failed! Err: %v\n", err)
		return
	}
	b, err := os.ReadFile(_config.File)
	if err != nil {
		fmt.Printf("read file! Err: %v\n", err)
		return
	}
	lines := bytes.Split(b, []byte("\n"))
	for _, line := range lines {
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		str := string(line)
		strArr := strings.Split(str, " ")
		cmd := strings.TrimSpace(strArr[0]) // 第一个为类型命令
		args := make([]interface{}, 0)
		for _, str := range strArr {
			args = append(args, strings.TrimSpace(str))
		}
		cmdType := strings.ToLower(cmdTypeMap[cmd])
		switch cmdType {
		case "duration":
			cmd := redis.NewDurationCmd(time.Second, args...)
			err = db.Process(cmd)
			if err != nil {
				fmt.Printf("db.Process Err: %s\n", err.Error())
				return
			}
			_, err = cmd.Result()
			printCmd(cmd.String(), err)

		case "int":
			cmd := redis.NewIntCmd(args...)
			err = db.Process(cmd)
			if err != nil {
				fmt.Printf("db.Process Err: %v\n", err.Error())
				return
			}
			_, err = cmd.Result()
			printCmd(cmd.String(), err)
		case "float":
			cmd := redis.NewFloatCmd(args...)
			err = db.Process(cmd)
			if err != nil {
				fmt.Printf("db.Process Err: %v\n", err.Error())
				return
			}
			_, err = cmd.Result()
			printCmd(cmd.String(), err)
		case "string":
			cmd := redis.NewStringCmd(args...)
			err = db.Process(cmd)
			if err != nil {
				fmt.Printf("db.Process Err: %v\n", err.Error())
				return
			}
			_, err = cmd.Result()
			printCmd(cmd.String(), err)
		case "stringstructmap":
			cmd := redis.NewStringStructMapCmd(args...)
			err = db.Process(cmd)
			if err != nil {
				fmt.Printf("db.Process Err: %v\n", err.Error())
				return
			}
			_, err = cmd.Result()
			printCmd(cmd.String(), err)
		case "stringstringmap":
			cmd := redis.NewStringStringMapCmd(args...)
			err = db.Process(cmd)
			if err != nil {
				fmt.Printf("db.Process Err: %v\n", err.Error())
				return
			}
			_, err = cmd.Result()
			printCmd(cmd.String(), err)
		case "stringslice":
			cmd := redis.NewStringSliceCmd(args...)
			err = db.Process(cmd)
			if err != nil {
				fmt.Printf("db.Process Err: %v\n", err.Error())
				return
			}
			_, err = cmd.Result()
			printCmd(cmd.String(), err)
		case "stringintmap":
			cmd := redis.NewStringIntMapCmd(args...)
			err = db.Process(cmd)
			if err != nil {
				fmt.Printf("db.Process Err: %v\n", err.Error())
				return
			}
			_, err = cmd.Result()
			printCmd(cmd.String(), err)
		case "time":
			cmd := redis.NewTimeCmd(args...)
			err = db.Process(cmd)
			if err != nil {
				fmt.Printf("db.Process Err: %v\n", err.Error())
				return
			}
			_, err = cmd.Result()
			printCmd(cmd.String(), err)
		case "zslice":
			cmd := redis.NewZSliceCmd(args...)
			err = db.Process(cmd)
			if err != nil {
				fmt.Printf("db.Process Err: %v\n", err.Error())
				return
			}
			_, err = cmd.Result()
			printCmd(cmd.String(), err)
		case "zwithkey":
			cmd := redis.NewZWithKeyCmd(args...)
			err = db.Process(cmd)
			if err != nil {
				fmt.Printf("db.Process Err: %v\n", err.Error())
				return
			}
			_, err = cmd.Result()
			printCmd(cmd.String(), err)
		case "slice":
			cmd := redis.NewSliceCmd(args...)
			err = db.Process(cmd)
			if err != nil {
				fmt.Printf("db.Process Err: %v\n", err.Error())
				return
			}
			_, err = cmd.Result()
			printCmd(cmd.String(), err)
		case "scan":
			cmd := redis.NewScanCmd(func(cmd redis.Cmder) error {
				err := db.Process(cmd)
				return err
			})
			//db.Options().ReadTimeout = 10 * time.Hour
			// err = db.Process(cmd)
			// if err != nil {
			// 	fmt.Printf("db.Process Err: %v\n", err.Error())
			// 	return
			// }
			printCmd(cmd.String(), err)
		case "cmd":
			err = errors.Errorf("not sport%s", cmd)
			printCmd("", err)
		case "status":
			cmd := redis.NewStatusCmd(args...)
			err = db.Process(cmd)
			if err != nil {
				fmt.Printf("db.Process Err: %v\n", err.Error())
				return
			}
			_, err = cmd.Result()
			printCmd(cmd.String(), err)
		case "boolslice":
			cmd := redis.NewBoolSliceCmd(args...)
			err = db.Process(cmd)
			if err != nil {
				fmt.Printf("db.Process Err: %v\n", err.Error())
				return
			}
			_, err = cmd.Result()
			printCmd(cmd.String(), err)
		case "bool":
			cmd := redis.NewBoolCmd(args...)
			err = db.Process(cmd)
			if err != nil {
				fmt.Printf("db.Process Err: %v\n", err.Error())
				return
			}
			_, err = cmd.Result()
			printCmd(cmd.String(), err)
		default:
			err = errors.Errorf("not sport%s", cmd)
			printCmd("", err)
		}
	}
}

func printCmd(cmdStr string, err error) {
	fmt.Printf("run:%s,err:%s\n", cmdStr, printError(err))
}

func printError(err error) (s string) {
	if err != nil {
		return err.Error()
	}
	return ""
}
