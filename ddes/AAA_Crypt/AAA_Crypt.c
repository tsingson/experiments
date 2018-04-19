#include "AAA_Crypt.h"
#include <string.h>
#include <stdlib.h>
#include <stdio.h>


#define DATA_SIZE 1024

#define FFMAX(a,b) ((a) > (b) ? (a) : (b))




int  AAA_Encrypt(const uint8_t  * inPut,int inlen,uint8_t  * outPut,int outlen)
{
    int SIZE=FFMAX(DATA_SIZE,FFMAX(inlen,outlen));
    uint8_t* DES_data=(uint8_t *)malloc(SIZE);
    if(DES_data==NULL)
        return -1;

    memset(DES_data,0,SIZE);
    int len=EncryptionFun(inlen, inPut, DES_data);

    if (BASE64_encode(outPut, outlen, DES_data, len) == NULL)
    {
        //printf("BASE64_encode error!\n");
        return -1;
    }

    return 0;

}
int  AAA_DeCrypt(const uint8_t  * inPut,int inlen,uint8_t  * outPut,int outlen)
{

    int SIZE=FFMAX(DATA_SIZE,FFMAX(inlen,outlen));
    uint8_t* DES_data=(uint8_t *)malloc(SIZE);
    if(DES_data==NULL)
        return -1;

    memset(DES_data,0,DATA_SIZE);
    int decodelen =BASE64_decode(DES_data,outlen,inPut,inlen);
    if (decodelen <= 0){
        //printf("BASE64_decode error!\n");
        return -1;
    }
    DecryptionFun(decodelen, DES_data, outPut);

    return 0;
}

uint8_t* DesEnCrypt(char *in_data_string ) {

    uint8_t* in_data=(uint8_t *)malloc(DATA_SIZE);
    uint8_t* out_data=(uint8_t *)malloc(DATA_SIZE);
    //  if(in_data==NULL)
    //    return  ;
    memset(in_data,0,DATA_SIZE);
    snprintf((char *)in_data,DATA_SIZE,"%s", in_data_string);

    AAA_Encrypt(in_data,strlen((char *)in_data),out_data,DATA_SIZE);
    return out_data;
}






struct MyString MyEnCrypt(char *in_data_string ) {
    int len = strlen((char*) in_data_string)

   // uint8_t* in_data=(uint8_t *)malloc(DATA_SIZE);
    uint8_t* out_data=(uint8_t *)malloc(DATA_SIZE);

    uint8_t* in_data=(uint8_t *)malloc(len);

  //  memset(in_data,0,DATA_SIZE);

    snprintf((char *)in_data,DATA_SIZE,"%s", in_data_string );

    AAA_Encrypt(in_data,strlen((char *)in_data),out_data,DATA_SIZE);
    printf("AAA加密结果:%s.\n",out_data);

    struct MyString str;
    int outlen = strlen((char *) out_data)

    char* p = malloc(outlen);


    memcpy(p, out_data, outlen);

    struct MyString str;
    str.s = p;
    str.len = len;

    return str;
}


int Test()
{

    uint8_t* in_data=(uint8_t *)malloc(DATA_SIZE);
    uint8_t* out_data=(uint8_t *)malloc(DATA_SIZE);
    if(in_data==NULL)
        return -1;
    memset(in_data,0,DATA_SIZE);
    snprintf((char *)in_data,DATA_SIZE,"%s","{\"status\": \"200 SUCCESS\", \"portal\": \"396e9b16-8818-46d3-ac6a-34ec99b7bbe9\", \"token\": \"734c151c-7b30-11e4-a01c-00e081b1495b\", \"userid\": \"af7baaac-1a65-400f-b767-8d3ae9876d2f\"}");
    printf("######################加密测试 begin:>>>>>>>>>>>>>##########################\n");
    AAA_Encrypt(in_data,strlen((char *)in_data),out_data,DATA_SIZE);
    printf("AAA加密结果:%s.\n",out_data);
    printf("######################加密测试 <<<<<<<<<<<<<<<end!##########################\n");

    printf("######################解密测试 begin:>>>>>>>>>>>>>##########################\n");
    memset(in_data,0,DATA_SIZE);
    AAA_DeCrypt(out_data,strlen((char*)out_data),in_data,DATA_SIZE);
    printf("AAA解密结果:%s.\n",in_data);
    printf("######################解密测试 <<<<<<<<<<<<<<<end!##########################\n");


    return 0;
}




