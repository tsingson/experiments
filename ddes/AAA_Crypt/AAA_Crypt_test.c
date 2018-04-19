
#ifdef TESTMAIN
int main()
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
#endif
