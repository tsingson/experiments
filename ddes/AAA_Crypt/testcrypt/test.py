#! /usr/bin/env python
#-*-coding:UTF-8-*-
import os;
import sys,json;
from ctypes import *
import ctypes
crypt_handle = cdll.LoadLibrary(os.getcwd() + '/libCrypt.so')


if __name__=='__main__':
    
    
    #url="register&info={\"username\":\"zhanghouming\",\"password\":\"beach\",\"sn\":\"12.00-09.10-10000000\"}"
    #2014-12-03加解密库出现问题，加密出来是4YCD1Mw3+l83xOPZCE16hLfPzWkU+T5tQZLazlo+kz73UvkjvVH59NVoGNqMok10bXXG4fInig5d3TQU3XkorxabJlgWy2JOgIqrOw7p6Y+FHC0/O5b8RKxf2HIJKXnqj/oC3eItF4M0owAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=
    
    #url="{\"status\": \"200 SUCCESS\", \"portal\": \"396e9b16-8818-46d3-ac6a-34ec99b7bbe9\", \"token\": \"734c151c-7b30-11e4-a01c-00e081b1495b\", \"userid\": \"af7baaac-1a65-400f-b767-8d3ae9876d2f\"}"
    url="play&info={\"sessionid\":989eea0f-3d70-498b-8172-b16f926a2eee,\"token\":4a836b6c-794b-11e4-8587-00e081b1495b,\"code\":\"123456789\"}"
    message1 = create_string_buffer(1024)
    print "################################加密#######################################"
    ret=crypt_handle.AAA_Encrypt(url,len(url),message1,1024)
    print "AAA_Encrypt 加密,ret=[%d],message:[%s]"%(ret,message1.value)
    
   
    print "################################解密#######################################"
    message2 = create_string_buffer(1024)
    ret=crypt_handle.AAA_DeCrypt(message1.value,len(message1.value),message2,1024)
    print "AAA_Encrypt 加密,ret=[%d],message:[%s]"%(ret,message2.value)
    
