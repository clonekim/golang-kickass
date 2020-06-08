# Go Bird

구글Drive를 이용한 컨텐츠매니저 프로그램입니다.  
gui는 없습니다. (만들어야 합니다)

# Build

```
$ go build gobird.go
```

# Run

반드시 구글클라우드 콘솔에서 서비스 계정을  
만든 후 credention파일(json)을 옵션으로 넣어야 합니다.

```
$ ./gobird --help

bonjour@bonjour:~/gobird$ ./gobird  --help
NAME:
   gobird 0.0.1 - google drive app implemented by Go.

USAGE:
   gobird [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --test              Test mode (default: false)
   --debug             debug option (default: false)
   --port value        port number (default: 8000)
   --credential value  credential file to google auth
   --help, -h          show help (default: false)
```
