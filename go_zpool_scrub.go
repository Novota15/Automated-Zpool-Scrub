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

  ln := strings.Split(string(stdout), "\n")
  for i := 1; i < len(ln); i++ {
    fmt.Println(ln[i])
    s := strings.Split(ln[i], " ")
    pool := s[0]
    fmt.Println(pool)
  }
  
  print(string(stdout))
  


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