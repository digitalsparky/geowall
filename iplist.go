package main

import (
  "net"
  "fmt"
)
type IPList struct {
  Versions []string
  Countries []string
}

type IPListResult struct {
  V4 []string
  V6 []string
}

type IPSubnet {
  Subnet string
  Version int
}


func (*i IPList) Get() *IPListResult {
  // Download the IP list for each individual version and country
  // Store in an array and return
  var v4list []IPSubnet
  var v6list []string

  var cidr = "192.168.10.244/24"

  // for each cidr returned from the API
  _, block, err := net.ParseCIDR(cidr)
  if err != nil {
    fmt.Errorf("Error! %s", err)

  }
  fmt.Println(block)




  return new IPListResult{
    V4: v4list
    V6: v6list
  }
}