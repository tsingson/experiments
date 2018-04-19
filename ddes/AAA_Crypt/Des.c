#include "Des.h"
#include <stdio.h>
#include <stdlib.h>

#define DATA_SIZE 1024
static struct AVDES g_des; 

static uint8_t des_key[] ={ 0xa3, 0xbe, 0x93, 0xff, 0x10, 0x34, 0x5f, 0xde,
    0xc6, 0x2e, 0x57, 0x83, 0x29, 0x7c, 0x8e, 0xf6,
    0xa3, 0x58, 0x34, 0x27, 0x13, 0x2c, 0x4e, 0xd2 };     
 


int EncryptionFun(int len,const uint8_t  * inPut,uint8_t  * outPut)
{
  
   
    len = len / 8 +(len % 8 == 0 ? 0 : 1);
    av_des_init(&g_des, des_key, 192, 0);
    av_des_crypt(& g_des, outPut,inPut, len, NULL, 0);
    return len*8;
    
}


int DecryptionFun(int len,const uint8_t  * inPut,uint8_t  * outPut)
{
    len = len / 8 +(len % 8 == 0 ? 0 : 1);
    av_des_init(&g_des, des_key, 192, 1);    
    av_des_crypt(&g_des, outPut, inPut, len, NULL, 1); 
    return len*8;
}









