#!/bin/bash

/usr/bin/taosadapter &
taosd &
./TDengineConnector 
