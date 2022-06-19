#!/usr/bin/python
# -*- coding: utf-8 -*-

import sys
import json
import requests
import time

cookies="PHPSESSID=sj67hqp37m3i1qqtpdc07qals4; hashtemp=e89aa09e622dc0c32502ed8164949e62; SERVERID=e4de745f5d5da3d9c53a7b3c040ad22b|1655605444|1655605387"

def updateTime(data):
    url = "https://www.zhifa315.com/index.php/site/study/update_time_new"
    headers = {
        "Content-Type": "application/x-www-form-urlencoded",
        "Cookie": cookies,
    }
    
    res = requests.post(url, data=data, headers=headers)
    if res.status_code != 200:
        print res
        return None
    print res.text
    return json.loads(res.text)

def getData(tm, mid, cid, ecid, start, source, status):
    return {
        "timestamp": tm,
        "material_id": mid,
        "c_id": cid,
        "playtime": 60,
        "ec_id": ecid,
        "time_start": start,
        "source": source,
        "status": status
    }

# curl -X POST -H "Content-Type: application/x-www-form-urlencoded"  -H "cookie: PHPSESSID=sj67hqp37m3i1qqtpdc07qals4; hashtemp=e89aa09e622dc0c32502ed8164949e62; SERVERID=e4de745f5d5da3d9c53a7b3c040ad22b|1655605444|1655605387" -d "timestamp=1792107_study_1655605444&material_id=2549&c_id=922&playtime=60&ec_id=2291&time_start=6501&source=1&status=2" https://www.zhifa315.com/index.php/site/study/update_time_new
if __name__ == "__main__":
    if len(sys.argv) < 7:
        print "usage: ./refresh.py timestamp mid cid ecid source status"
        exit(1)


    start = 60
    last_total_time = 0
    while True:
        data = getData(sys.argv[1], sys.argv[2], sys.argv[3], sys.argv[4], start, sys.argv[5], sys.argv[6])
        res = updateTime(data)
        if res == None:
            break

        print res
        if res["total_time"] == last_total_time:
            print "refresh succ, total=%d"%last_total_time
            break
        
        last_total_time = res["total_time"]
        start = last_total_time + 60
        time.sleep(1)
