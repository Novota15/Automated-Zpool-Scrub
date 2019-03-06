package main

//TO DO: Add sorting algorithmm


import (
   "fmt"
   "strings"
   //"os"
   "os/exec"
   //"bytes"
   //"io"
   "strconv"
   "time"
   //"parseany"
)

type Pool struct {
  Name string
  Scan string
  Scan_Date time.Time 
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
    fmt.Println(pool)
    pools[i-1].Name = pool
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
    fmt.Println(scan_output)
    pools[i].Scan = scan_output
  }
  return
}

//sort pools list by time of scrubs
func Get_zpool_Scrub_Date(pools []Pool) {
  //parse through scan info to get info about month, day, and year
  const shortForm = "2006-Jan-02"
  for k := 0; k < len(pools); k++ {
    fmt.Println(pools[k].Name)
    for i := 2; i < len(string(pools[k].Scan)); i++ {
      if string(pools[k].Scan[i-2:i]) == "on" { //date of scrub begins after "on"
        //fmt.Println("on")
        i = i + 8
        month := string(pools[k].Scan[i-3:i])
        i = i + 3
        day := string(pools[k].Scan[i-2:i])
        //add 0 to the day if day < 10
        x, _ := strconv.Atoi(day)
        if x < 10 {
          day = "0" + string(pools[k].Scan[i-1:i])
        }
        i = i + 14
        year := string(pools[k].Scan[i-4:i])
        //fmt.Println(month, day, year)
        date := year + "-" + month + "-" + day
        t, _ := time.Parse(shortForm, date)
        pools[k].Scan_Date = t
        //fmt.Println(t)
        //fmt.Println(pools[k].Scan_Date)
      }
      // fmt.Println(string(item))
    }
  }
  //fmt.Println("here:")
  //fmt.Println(pools[0].Scan_Date)
}

// func swap(a Pool, b Pool) {
//   temp = *a
//   *a = *b
//   *b = temp
// }

// func partition(pools []Pool, low int, high int) int{
//   pivot := pools[high]
//   i := low - 1

//   for j := low; j <= high - 1; j++ {
//     if time.Date(pools[j].Scan_Date).After(time.Date(pivot.Scan_Date)) {
//       i++
//       swap(pools[i], pools[j])
//     }
//   }
//   swap(pools[i+1], pools[high])
//   return (i+1)
// }

// func Sort_by_Date(pools []Pool, low int, high int) { //quicksort algorithm
//   t1 := time.Date(pools[low].Scan_Date)
//   t2 := time.Date(pools[high].Scan_Date)
//   if t1.After(t2) {
//     pi := partition(pools, low, high)
//     Sort_by_Date(pools, low, pi - 1)
//     Sort_by_Date(pools, pi + 1, high)
//   }
// }

//returns the index of the pool with the oldest scrub
func Find_Oldest_Scrub(pools []Pool) int{
  j := 0
  for i := 1; i < len(pools); i++ {
    if pools[j].Scan_Date.After(pools[i].Scan_Date) {
      j = i
    }
  }
  return j
}

func Perform_Scrub(pool Pool) {
  cmd := exec.Command("bash", "-c", "zpool scrub" + pool.Name)

  _, err := cmd.Output()

    if err != nil {
      println(err.Error())
      return
    }
  fmt.Println("ran scrub on " + pool.Name)

  return
}

func main() {

  pools := Get_zpool_Names()
  Get_zpool_scan(pools)
  Get_zpool_Scrub_Date(pools)
  j := Find_Oldest_Scrub(pools)
  fmt.Println("Oldest Scrub:")
  fmt.Println(pools[j])
  Perform_Scrub(pools[j])

}

