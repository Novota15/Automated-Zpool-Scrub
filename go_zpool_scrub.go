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

  s := strings.Split(string(stdout), "/n")
  pool1, pool2 := s[0], s[1]
  print(string(stdout))
  fmt.Println(pool1, pool2)
  


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