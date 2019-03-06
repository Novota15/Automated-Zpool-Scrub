# Go zpool Scrub

This is a go program that reads the dates that each zpool was scrubbed and executes the zpool scrub command on the zpool that was scrubbed the longest time ago (accurate up to the day).

## Installing Go

use ```sudo pkg install go```

make sure ```/usr/local/go/bin``` is added to PATH which can be done by using the command ```PATH=~/``` and adding to the current PATH. Do not replace the current PATH, just add to what is seen with ```echo $PATH```

## Running the program

The program must be run with root privileges in order for the scrub command to execute correctly
 
use ```sudo go run go_zpool_scrub.go```