#!/bin/bash
#while true; do sleep 1; echo 1; done

[[ -x /usr/bin/taosadapter ]] && /usr/bin/taosadapter &
taosd
#./TDengineConnector
