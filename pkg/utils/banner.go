package utils

import (
	"fmt"
)

const syncshASCII = `
 _____                      _____  _     
/  ___|                    /  ___|| |    
\ ` + "`" + `--. _   _  _ __    ___ \ ` + "`" + `--. | |__  
 ` + "`" + `--. \| | | || '_ \  / __| ` + "`" + `--. \| '_ \ 
/\__/ /| |_| || | | || (__ /\__/ /| | | |
\____/  \__, ||_| |_| \___|\____/ |_| |_|
         __/ |                           
        |___/                            
`

func ShowSyncshBanner() {
	fmt.Print(syncshASCII)
	fmt.Println()
}

func ShowInitBanner() {
	ShowSyncshBanner()
}
