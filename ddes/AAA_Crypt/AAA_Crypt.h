#ifndef AAA_CRYPT__H___
#include "Des.h"
#include "Base64.h"

struct MyString
{
    char* s;
    int len;
};
int  AAA_Encrypt(const uint8_t  * inPut,int inlen,uint8_t  * outPut,int outlen);
int  AAA_DeCrypt(const uint8_t  * inPut,int inlen,uint8_t  * outPut,int outlen);
uint8_t* DesEnCrypt(char *in_data_string );
struct MyString MyEnCrypt(char *in_data_string );
int Test();

#endif
