#!/bin/python3

# [{"Key":"metric/cpu","CreateIndex":105,"ModifyIndex":105,"LockIndex":0,"Flags":0,"Value":"MTA=","Session":""},{"Key":"metric/cpu1","CreateIndex":175,"ModifyIndex":175,"LockIndex":0,"Flags":0,"Value":"MjA=","Session":""}]

with open('/app/log.txt','+w') as f:
    inp = input()
    f.write(inp)