package logs

import (
	"time"
	"fmt"
	"os"

	"github.com/rprobaina/lpfs"
)

func Logs() error {

	file, err := os.OpenFile("./logs.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Errorf("unable to open the file logs.txt")
		return err
	}

	t := time.Now().Format(time.DateTime)
	c, _ := lpfs.GetMemCached()
	f, _ := lpfs.GetSwapFilename()

	msg := fmt.Sprintf("Date: %v Memory: %v Swapfile: %v\n", t, c, f)

	_, err = file.WriteString(msg)
	if err != nil {
		fmt.Errorf("unable to write in the file logs.txt")
		return err
	}

	file.Close()

	return nil
}
