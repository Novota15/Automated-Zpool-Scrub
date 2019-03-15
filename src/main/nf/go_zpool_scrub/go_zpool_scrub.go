package main

import (
   "fmt"
   "strings"
   "os/exec"
   "strconv"
   "time"
   "flag"
)

type Pool struct {
  Name string
  State string
  Scan string
  Scan_Date time.Time
  Scanned bool
}

func Get_All_zpools(exclusion_list []string) []Pool {
  cmd := exec.Command("bash", "-c", "zpool list -H -o name,health")
  stdout, err := cmd.Output()

  p := make([]Pool, 5)

  if err != nil {
    fmt.Println("Not Working")
    println(err.Error())
    return p
  }

  ln := strings.Split(string(stdout), "\n") //split into lines
  
  pools_size := len(ln) - 1
  pools := make([]Pool, pools_size)
  excluded_pools_count := 0
  //creating a list of ALL pools
  for i := 0; i < pools_size; i++ {
    data := strings.Split(ln[i], "\t")
    pool_name := data[0]
    pool_health := data[1]
    pools[i].Name = pool_name
    pools[i].Scanned = false
    pools[i].State = pool_health
    //check if pool is in exclusion list to get a count of pools needed to be taken out
    for j := 0; j < len(exclusion_list); j++ {
      if pools[i].Name == exclusion_list[j] {
        excluded_pools_count++
      }
    }
    //fmt.Println(pools[i])
  }
  //create new pool list
  pools2 := make([]Pool, (pools_size - excluded_pools_count))
  index_correction := 0
  included := false
  //create new list that does not include pools in exclusion list
  for j := 0; j < pools_size; j++ {
    for k := 0; k < len(exclusion_list); k++ {
      if pools[j].Name != exclusion_list[k] {
        pools2[j - index_correction] = pools[j]
        included = true
      }
    }
    if included == true {
      index_correction++
      included = false
    }
  }
  fmt.Println(pools2)
  return pools2
}

func Get_Online_zpools(pools []Pool) []Pool{
  length := 0
  for i := 0; i < len(pools); i++ {
    if pools[i].State == "ONLINE" {
      length++
    }
  }
  online_pools := make([]Pool, length)
  fmt.Println("Online Pools: ")
  for k := 0; k < length; k++ {
    if pools[k].State == "ONLINE" {
      online_pools[k] = pools[k]
      fmt.Println(online_pools[k])
    }
  }
  return online_pools
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
    scan := ""
    for k := 0; k < len(ln); k++ { //go through each line
      line := string(ln[k])
      for j := 0; j < (len(line) - 6); j++ { //search for the line with scan info
        if string(line[j:j+4]) == "scan" {
          scan = line
          break
        }
      }
    }
    scan_output := scan //line containing the scrub info
    //fmt.Println(scan_output)
    pools[i].Scan = scan_output
    //fmt.Println(pools[i].Name)
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
        i = i + 3
        day := string(pools[k].Scan[i-2:i])
        //add 0 to the day if day < 10
        x, _ := strconv.Atoi(day)
        //fmt.Println("day of scan seen: " + day)
        if x < 10 {
          day = "0" + string(pools[k].Scan[i-1])
        }
        i = i + 13
        year := string(pools[k].Scan[i-3:i+1])
        //fmt.Println("date seen: ", month, day, year)
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

  return
}

func Parse_Exclusion_List(exclude string) []string{
  exclusion_list := strings.Split(exclude, ",")
  
  return exclusion_list
}

func main() {
  //get zpool exclusion list from command line flag
  wordPtr := flag.String("exclude", "foo", "a string")
  flag.Parse()
  //fmt.Println("exclusion list: ", *wordPtr)
  //fmt.Println("tail: ", flag.Args())
  exclusion_list := Parse_Exclusion_List(*wordPtr)

  pools := Get_All_zpools(exclusion_list)
  online_pools := Get_Online_zpools(pools)
  Scrub_Least_Recent(online_pools)
}

