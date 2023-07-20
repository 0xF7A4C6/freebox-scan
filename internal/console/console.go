package console

import (
	"fmt"
	"strings"
	"time"
	"github.com/gookit/color"
)

func SetTitle(title string) {
	fmt.Printf("\033]0;%s\007", title)
}

func Log(Content string) {
	date := strings.ReplaceAll(time.Now().Format("15:04:05"), ":", "<fg=353a3b>:</>")
	content := fmt.Sprintf("[%s] %s.", date, Content)

	content = strings.ReplaceAll(content, "[", "<fg=6c7273>[</>")
	content = strings.ReplaceAll(content, "]", "<fg=6c7273>]</>")

	content = strings.ReplaceAll(content, "(", "<fg=6c7273;op=bold>(")
	content = strings.ReplaceAll(content, ")", ")</>")

	content = strings.ReplaceAll(content, "HIT", "<fg=7deb46>HIT</>")
	content = strings.ReplaceAll(content, "NOP", "<fg=f0543c>NOP</>")

	content = strings.ReplaceAll(content, "FOUND", "<fg=4af7ee>FOUND</>")
	content = strings.ReplaceAll(content, "INVALID", "<fg=f7a14a>INVALID</>")

	color.Println(content)
}
