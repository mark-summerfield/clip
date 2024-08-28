module github.com/mark-summerfield/clip

go 1.23

replace github.com/mark-summerfield/uterm => /home/mark/app/golib/uterm

replace github.com/mark-summerfield/set => /home/mark/app/golib/set

require (
	github.com/kopoli/go-terminal-size v0.0.0-20170219200355-5c97524c8b54
	github.com/mark-summerfield/set v1.0.0
	github.com/mark-summerfield/uterm v1.0.0
)

require (
	github.com/mattn/go-isatty v0.0.20 // indirect
	golang.org/x/sys v0.24.0 // indirect
)
