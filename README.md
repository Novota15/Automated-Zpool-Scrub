# Go zpool Scrub

This is a program written in GO that automates zpool scrubbing.

It reads the dates that each zpool was scrubbed and executes the ```zpool scrub``` command on the zpool that was scrubbed the longest time ago (accurate up to the day).

## When Installing Go

make sure ```/usr/local/go/bin``` is added to PATH which can be done by using the command ```PATH=~/``` and adding to the current PATH. Do not replace the current PATH, just add to what is seen with ```echo $PATH```

## Running the program

The program must be run with root privileges in order for the scrub command to execute correctly

You can use ```sudo go run go_zpool_scrub.go```

OR run the executable by typing ```~/go-zpool-scrub/bin/go_zpool_scrub```

The executable is created by executing ```go install``` inside the directory ```go-zpool-scrub/src/main/nf/go_zpool_scrub``` which will place the executable inside the ```/bin``` directory

###### Running as a Cron Job

write ```sudo crontab -e```
edit the crontab:
1) press ```esc```
2) press ```i``` to begin editing the file
3) paste cron command into the file.

Example for running at 1 am every day: 
```0 1 * * * /usr/local/go/bin /home/student/go-zpool-scrub/go_zpool_scrub.go > /dev/null 2>&1```.

4) press ```esc``` again to exit editing mode
5) type ```:wq``` to save

cron jobs can be viewed with ```sudo crontab -l```