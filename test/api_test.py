# -*- coding: utf-8 -*-  

import requests
import json

# 关于：requests的传参
# params的时候之间接把参数加到url后面，只在get请求时使用：
# data用于post请求：不传递Content-Type，则默认使用form来post；如果设置Content-Type，则按Content-Type来


def testAuth():
    url = 'http://127.0.0.1:8088/2b/v1/user/user_info'
    headers = {'Useragent': json.dumps({'app_version': '1.01.01', 'mobile_platform': "1", 'mobile_system': 'ios', 'mobile_device_brand': 'apple', 'mobile_os_version': '8.0.1'}), 'Authorization': '5oFKwvqs2doI3czvLW/kyRUFJZeNfyh/ecLzVOBLkkSpxQOyRkZK6k4792ncesNk'}

    r = requests.get(url, headers=headers)
    print r.headers
    print r.content


def testLogin():
    url = 'http://127.0.0.1:8088/2b/v1/user/login'
    headers = {'Useragent': json.dumps({'app_version': '1.01.01', 'mobile_platform': "1", 'mobile_system': 'ios', 'mobile_device_brand': 'apple', 'mobile_os_version': '8.0.1'}), 'Content-Type': 'application/json'}
    data = {'MyAppId':0, 'WxId':'SDEFXCRYWFEG',  'city':'Chengdu', 'gender':1, 'nickName':'GHW💻', 'province':'Sichuan', 'avatarUrl':'https://wx.qlogo.cn/mmopen/vi_32/DYAIOgq83eolibHibbzicClicSTASgkgFpfhSrDKkT36DTXmJFk4dp3iaYZHHbRo8wuWdq8fsMDF3ib4DJL1Va1z6dNQ/132'}

    r = requests.post(url, headers=headers, data=json.dumps(data))
    print r.request.body
    print r.headers
    print r.text


def testGetUserInfo():
    url = 'http://127.0.0.1:8088/2b/v1/user/user_detail_info'
    headers = {'Useragent': json.dumps({'app_version': '1.01.01', 'mobile_platform': "1", 'mobile_system': 'ios', 'mobile_device_brand': 'apple', 'mobile_os_version': '8.0.1'}), 'Authorization': 'o5rDuSzSeOMhMCc+B0nnKqPB2I/XYZ0Of11dTDXFnEwWXKSCKXH+kX7fxagB5ugpiKDzlA/7l4JoxtlZb5nKUg=='}
    data = {}

    r = requests.get(url, headers=headers, params=data)
    print r.headers
    print r.content

    url = 'http://127.0.0.1:8088/2b/v1/user/user_simple_info'
    data = {'UserId':2}

    r = requests.get(url, headers=headers, params=data)
    print r.url
    print r.headers
    print r.content

def testCreateHouse():
    url = 'http://127.0.0.1:8088/2b/v1/house/create_house'
    headers = {'Useragent': json.dumps({'app_version': '1.01.01', 'mobile_platform': "1", 'mobile_system': 'ios', 'mobile_device_brand': 'apple', 'mobile_os_version': '8.0.1'}), 'Authorization': 'o5rDuSzSeOMhMCc+B0nnKqPB2I/XYZ0Of11dTDXFnEwWXKSCKXH+kX7fxagB5ugpiKDzlA/7l4JoxtlZb5nKUg==', 'Content-Type': 'application/json'}
    data = {'Name':'XXXXXX'}

    r = requests.post(url, headers=headers, data=json.dumps(data))
    print r.headers
    print r.content

def testGetHouseInfo():
    url = 'http://127.0.0.1:8088/2b/v1/house/house_info'
    headers = {'Useragent': json.dumps({'app_version': '1.01.01', 'mobile_platform': "1", 'mobile_system': 'ios', 'mobile_device_brand': 'apple', 'mobile_os_version': '8.0.1'}), 'Authorization': 'o5rDuSzSeOMhMCc+B0nnKqPB2I/XYZ0Of11dTDXFnEwWXKSCKXH+kX7fxagB5ugpiKDzlA/7l4JoxtlZb5nKUg==', 'Content-Type': 'application/json'}

    r = requests.get(url, headers=headers)
    print r.headers
    print r.content

def testSetHouse(postFix, data):
    url = 'http://127.0.0.1:8088/2b/v1/house/' + postFix
    headers = {'Useragent': json.dumps({'app_version': '1.01.01', 'mobile_platform': "1", 'mobile_system': 'ios', 'mobile_device_brand': 'apple', 'mobile_os_version': '8.0.1'}), 'Authorization': 'o5rDuSzSeOMhMCc+B0nnKqPB2I/XYZ0Of11dTDXFnEwWXKSCKXH+kX7fxagB5ugpiKDzlA/7l4JoxtlZb5nKUg==', 'Content-Type': 'application/json'}

    r = requests.post(url, headers=headers, data=json.dumps(data))
    print r.headers
    print r.content

def getUA():
    return json.dumps({'app_version': '0.0.1', 'mobile_platform': "devtools", 'mobile_system': 'iOS 10.0.1', 'mobile_device_brand': 'devtools'})

def getAuth():
    return "hoeUBQMwY4JBSsH9ER4wmvlT8qmfmvydOloklGLGrN7PqSR6M3WC7SyhdAl5vPon1JilpXW+iYPOkknxXE+OSg=="

def testGetOrder(offset, limit):
    url = 'http://127.0.0.1:7070/v1/order/list_order?offset=' + str(offset) + "&limit=" + str(limit)
    headers = {'Useragent': getUA(), 'Authorization': getAuth(), 'Content-Type': 'application/json'}

    r = requests.get(url, headers=headers)
    print r.headers
    print r.content

def testSetHouseInfo() :
    testSetHouse("set_name", {'Name': 'X1'})
    testSetHouse("set_pics", {'PicUrlList': [ "url1", "url2", "url3" ]})
    testSetHouse("set_desc", {'DescType': 1, 'Content': 'short desc'})
    testSetHouse("set_desc", {'DescType': 2, 'Content': 'detail desc'})
    testSetHouse("set_free_meal", {'FreeMealByBookRoom': False})
    testSetHouse("set_location", {'Latitude': 2.2222, 'Longitude': 3.333})
    testSetHouse("set_pickup_spot", {'PickupSpotLat': 2.2222, 'PickupSpotLng': 3.333})

def main():
    #testLogin()
    #testAuth()
    #testGetUserInfo()
    #testCreateHouse()
    #testGetHouseInfo()
    #testSetHouseInfo()

    testGetOrder(0, 10)

main()
