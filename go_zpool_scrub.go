package main

import (
   "fmt"
   "strings"
   "os/exec"
   "strconv"
   "time"
)

type Pool struct {
  Name string
  Scan string
  Scan_Date time.Time
  Scanned bool
}

//creates Pool struct for each pool and stores in pools list
func Get_zpool_Names() []Pool {
  cmd := exec.Command("bash", "-c", "zpool list")

  stdout, err := cmd.Output()

  p := make([]Pool, 5)

  if err != nil {
    fmt.Println("Not Working")
    println(err.Error())
    return p
  }

  ln := strings.Split(string(stdout), "\n") //split into lines
  //create array for storing pool structs that has length ln - 2
  pool_size := len(ln) - 2
  pools := make([]Pool, pool_size)
  for i := 1; i < len(ln) - 1; i++ { //iterate thru each line
    //fmt.Println(ln[i])
    s := strings.Split(ln[i], " ") //split each line
    pool := s[0] //name of the pool ->start of each line
    //fmt.Println(pool)
    pools[i-1].Name = pool
    pools[i-1].Scanned = false
  }
  
  return pools
}

//gets the scan info for each pool in pools list
func Get_zpool_scan(pools []Pool) {
  for i := 0; i < len(pools); i++ {
    cmd1 := exec.Command("bash", "-c", "zpool status " + pools[i].Name)

    stdout1, err1 := cmd1.Output()

    if err1 != nil {
      println(err1.Error())
      return
    }

    //find the date of last srub and store in Status
    ln := strings.Split(string(stdout1), "\n")
    scan_output := ln[2] //line containing the scrub info
    //fmt.Println(scan_output)
    pools[i].Scan = scan_output
    //fmt.Println(string(pools[i].Scan))
  }
  return
}

//sort pools list by time of scrubs
func Get_zpool_Scrub_Date(pools []Pool) {
  //parse through scan info to get info about month, day, and year
  const shortForm = "2006-Jan-02"
  for k := 0; k < len(pools); k++ {
    fmt.Println(pools[k].Name)
    for i := 3; i < len(string(pools[k].Scan)); i++ {
      if string(pools[k].Scan[i-3:i]) == "on " { //date of scrub begins after "on"
        //fmt.Println("on")
        pools[k].Scanned = true
        i = i + 7
        month := string(pools[k].Scan[i-3:i])
        i = i + 2
        day := string(pools[k].Scan[i-2:i])
        //add 0 to the day if day < 10
        x, _ := strconv.Atoi(day)
        if x < 10 {
          day = "0" + string(pools[k].Scan[i-1:i])
        }
        i = i + 13
        year := string(pools[k].Scan[i-4:i])
        //fmt.Println(month, day, year)
        date := year + "-" + month + "-" + day
        t, _ := time.Parse(shortForm, date)
        pools[k].Scan_Date = t
        //fmt.Println(t)
        fmt.Println("adding scan date to " + pools[k].Name)
        fmt.Println(pools[k].Scan_Date)
      } else if i == (len(string(pools[k].Scan)) - 2) {
        //fmt.Println(pools[k].Name + " hasn't been scrubbed")
        pools[k].Scanned = false
        break
      }
    }
  }
}

//returns the index of the pool with the oldest scrub
func Find_Oldest_Scrub(pools []Pool) int{
  j := 0
  for j = 0; j < len(pools); j++ { //check if any pool has never been scrubbed
    if pools[j].Scanned == false {
      fmt.Println(pools[j].Name + " hasn't been scrubbed yet")
      return j
    }
  }
  j = 0
  for i := 1; i < len(pools); i++ { //if all have been scrubbed, find the oldest
    if pools[i].Scanned == true {
      fmt.Println(pools[j].Scan_Date)
      fmt.Println(pools[i].Scan_Date)
      if pools[j].Scan_Date.Before(pools[i].Scan_Date) {
      j = i
      }
    }
  }
  return j
}

func Perform_Scrub(pool Pool) {
  cmd := exec.Command("bash", "-c", "sudo zpool scrub " + pool.Name)

  _, err := cmd.Output()

    if err != nil {
      println(err.Error())
      return
    }

  fmt.Println("ran scrub on " + pool.Name)

  return
}

func Scrub_Least_Recent(pools []Pool) {
  Get_zpool_scan(pools)
  Get_zpool_Scrub_Date(pools)
  j := Find_Oldest_Scrub(pools)
  fmt.Println("pool to be scrubbed:")
  fmt.Println(pools[j])
  Perform_Scrub(pools[j])
}

func main() {
  pools := Get_zpool_Names()
  fmt.Println("length of pools list")
  fmt.Println(len(pools))
  Scrub_Least_Recent(pools)
}

