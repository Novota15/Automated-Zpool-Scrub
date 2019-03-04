package main

import (
   "fmt"
   "strings"
   //"os"
   "os/exec"
   //"bytes"
   //"io"
)

type Pool struct {
  Name string
  Scan string
  Scan_Date int 
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
func Sort_zpool_scrubs(pools []Pool) {
  //parse through scan info to get info about month, day, and year
  for _, pool := range pools {
    fmt.Println(pool.Name)
    for i := 2; i < len(string(pool.Scan)); i++ {
      if string(pool.Scan[i-2:i]) == "on" { //date of scrub begins after "on"
        fmt.Println("on")
        i = i + 8
        month := string(pool.Scan[i-3:i])
        i = i + 3
        day := string(pool.Scan[i-2:i])
        i = i + 14
        year := string(pool.Scan[i-4:i])
        fmt.Println(month, day, year)
        Convert_Date_to_Int(month, day, year)
      }
      // fmt.Println(string(item))
    }
  }
}

func Convert_Date_to_Int(month, day, year string) int {
  return 5
}

func main() {
  //arg0 := "zpool status"

  //use zpool list to get zpools and store in pool struct

  //cmd := exec.Command("bash", "-c", "zpool status")

  pools := Get_zpool_Names()
  Get_zpool_scan(pools)
  Sort_zpool_scrubs(pools)
  //call zpool status on each pool and store status in pool struct
  //print(string(stdout))
  

}



// func main() {
//     app := "echo"
//     //app := "buah"

//     arg0 := "-e"
//     arg1 := "Hello world"
//     arg2 := "\n\tfrom"
//     arg3 := "golang"

//     cmd := exec.Command(app, arg0, arg1, arg2, arg3)
//     stdout, err := cmd.Output()

//     if err != nil {
//         println(err.Error())
//         return
//     }

//     print(string(stdout))
// }