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
  State string
  Scan string
  Scan_Date time.Time
  Scanned bool
}

//creates Pool struct for each pool and stores in pools list
func Get_zpool_Names() []Pool {
  cmd := exec.Command("bash", "-c", "zpool list")
  //cmd := exec.Command("bash", "-c", "zpool list -H -o name,health")

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
    //health := s[1]
    //fmt.Println(pool)
    //fmt.Println(health)
    pools[i-1].Name = pool
    pools[i-1].Scanned = false
    //pools[i-1].State = health
  }
  
  return pools
}

func Get_zpools() []Pool {
  cmd := exec.Command("bash", "-c", "zpool list -H -o name,health")
  stdout, err := cmd.Output()

  p := make([]Pool, 5)

  if err != nil {
    fmt.Println("Not Working")
    println(err.Error())
    return p
  }

  ln := strings.Split(string(stdout), "\n") //split into lines
  fmt.Println(ln)
  pool_size := len(ln)
  for k := 0; k < len(string(ln); k++ {
    fmt.Println("k: " + stirng(ln[k]))
  }
  fmt.Println("pool list size: " + string(pool_size))
  pools := make([]Pool, pool_size)

  fmt.Println("creating pool list: ")
  for i := 0; i < len(ln); i++ {
    data := strings.Split(ln[i], " ")
    pool_name := data[0]
    pool_health := data[1]
    pools[i].Name = pool_name
    pools[i].Scanned = false
    pools[i].State = pool_health
    fmt.Println(pools[i])
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

    //find the date of last scrub and store in Status
    ln := strings.Split(string(stdout1), "\n")

    //find the line with the scan info:
    scan := 0
    for k := 0; k < len(ln); k++ { //go through each line
      line := string(ln[k])
      for j := 0; j < (len(line) - 6); j++ { //search for the line with scan info
        if string(line[j:j+4]) == "scan" {
          scan = j
          break
        }
      }
    }
    scan_output := ln[scan] //line containing the scrub info
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
          day = "0" + string(pools[k].Scan[i])
        }
        i = i + 15
        year := string(pools[k].Scan[i-4:i])
        //fmt.Println(month, day, year)
        date := year + "-" + month + "-" + day
        t, _ := time.Parse(shortForm, date)
        pools[k].Scan_Date = t
        //fmt.Println(t)
        //fmt.Println("adding scan date to " + pools[k].Name)
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
  //all pools have been scrubbed so time to compare dates
  j = 0
  for i := 1; i < len(pools); i++ { //if all have been scrubbed, find the oldest
    if pools[i].Scanned == true {
      if pools[j].Scan_Date.After(pools[i].Scan_Date) {
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
  pools := Get_zpools()
  Scrub_Least_Recent(pools)
}

