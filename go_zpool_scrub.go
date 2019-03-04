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
  Status string
}
//TO DO: modularize code into separate functions

func main() {
  //arg0 := "zpool status"

  //use zpool list to get zpools and store in pool struct

  //cmd := exec.Command("bash", "-c", "zpool status")

  cmd := exec.Command("bash", "-c", "zpool list")

  stdout, err := cmd.Output()

  if err != nil {
    println(err.Error())
    return
  }

  ln := strings.Split(string(stdout), "\n") //split into lines
  //create array for storing pool structs that has length ln - 1
  pool_size := len(ln) - 2
  pools := make([]Pool, pool_size)
  for i := 1; i < len(ln) - 1; i++ { //iterate thru each line
    //fmt.Println(ln[i])
    s := strings.Split(ln[i], " ") //split each line
    pool := s[0] //name of the pool ->start of each line
    fmt.Println(pool)
    pools[i-1].Name = pool
  }
  fmt.Println(len(pools))
  //call zpool status on each pool and store status in pool struct

  for i := 0; i < len(pools); i++ {
    cmd := exec.Command("bash", "-c", "zpool status", pools[i].Name)

    stdout, err := cmd.Output()

    if err != nil {
      println(err.Error())
      return
    }
    //find the date of last srub and store in Status
    pools[i].Status = string(stdout)
    //fmt.Println(string(stdout))
  }

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