#ifndef _HELPER_INCLUDED
#define _HELPER_INCLUDED

#include <stdio.h>
#include <iostream>
#include <Windows.h>
#include <algorithm>
#include <string.h>

#include "Tasking.h"
#include "Libraries/cJSON/cJSON.h"
#include "Libraries/TinyAES/aes.hpp"

void err(const char* errval);

wchar_t* convertCharArrayToLPCWSTR(const char* charArray);


std::string Encode(const std::string data);
std::string Decode(const std::string& input, std::string& out);

std::string cleaner(std::string encoded);

// Serialization + Deserialization

std::string SerializeData(Response *response);
Task DeserializeData(std::string data);

// registeration 
BOOL GetInfo(std::string& domain, std::string& strUser, std::string& strdomain);
BOOL GetLogonFromToken(HANDLE hToken, std::string& strUser, std::string& strdomain);
VOID regSerDeSer(State* state, std::string& data, BOOL in);


#endif