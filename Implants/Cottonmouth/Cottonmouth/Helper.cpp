#include <algorithm>
#include "Helper.h"
#include "iphlpapi.h"

#define MAX_NAME  256


void err(const char* errval) {

	printf("Error: %s\nError Code:%d", errval, GetLastError());
	//exit(1);

}

wchar_t* convertCharArrayToLPCWSTR(const char* charArray) {

	wchar_t* wString = new wchar_t[4096];
	ZeroMemory(wString, 4096 * sizeof(wchar_t));
	MultiByteToWideChar(CP_ACP, 0, charArray, -1, wString, 4096);
	return wString;

}

std::string Encode(const std::string data) {
    static constexpr char sEncodingTable[] = {
      'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H',
      'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P',
      'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X',
      'Y', 'Z', 'a', 'b', 'c', 'd', 'e', 'f',
      'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n',
      'o', 'p', 'q', 'r', 's', 't', 'u', 'v',
      'w', 'x', 'y', 'z', '0', '1', '2', '3',
      '4', '5', '6', '7', '8', '9', '+', '/'
    };

    size_t in_len = data.size();
    size_t out_len = 4 * ((in_len + 2) / 3);
    std::string ret(out_len, '\0');
    size_t i;
    char* p = const_cast<char*>(ret.c_str());

    for (i = 0; i < in_len - 2; i += 3) {
        *p++ = sEncodingTable[(data[i] >> 2) & 0x3F];
        *p++ = sEncodingTable[((data[i] & 0x3) << 4) | ((int)(data[i + 1] & 0xF0) >> 4)];
        *p++ = sEncodingTable[((data[i + 1] & 0xF) << 2) | ((int)(data[i + 2] & 0xC0) >> 6)];
        *p++ = sEncodingTable[data[i + 2] & 0x3F];
    }
    if (i < in_len) {
        *p++ = sEncodingTable[(data[i] >> 2) & 0x3F];
        if (i == (in_len - 1)) {
            *p++ = sEncodingTable[((data[i] & 0x3) << 4)];
            *p++ = '=';
        }
        else {
            *p++ = sEncodingTable[((data[i] & 0x3) << 4) | ((int)(data[i + 1] & 0xF0) >> 4)];
            *p++ = sEncodingTable[((data[i + 1] & 0xF) << 2)];
        }
        *p++ = '=';
    }

    // escape + 
    for (int i = 0; i < ret.length(); ++i) {
        if (ret[i] == '+')
            ret.replace(i, 1, "%2b");
    }
	
    return ret;
}

std::string Decode(const std::string& input, std::string& out) {
    static constexpr unsigned char kDecodingTable[] = {
      64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64,
      64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64,
      64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 62, 64, 64, 64, 63,
      52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 64, 64, 64, 64, 64, 64,
      64,  0,  1,  2,  3,  4,  5,  6,  7,  8,  9, 10, 11, 12, 13, 14,
      15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 64, 64, 64, 64, 64,
      64, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40,
      41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 64, 64, 64, 64, 64,
      64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64,
      64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64,
      64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64,
      64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64,
      64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64,
      64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64,
      64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64,
      64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64, 64
    };

    size_t in_len = input.size();
    if (in_len % 4 != 0) return "Input data size is not a multiple of 4";

    size_t out_len = in_len / 4 * 3;
    if (input[in_len - 1] == '=') out_len--;
    if (input[in_len - 2] == '=') out_len--;

    out.resize(out_len);

    for (size_t i = 0, j = 0; i < in_len;) {
        uint32_t a = input[i] == '=' ? 0 & i++ : kDecodingTable[static_cast<int>(input[i++])];
        uint32_t b = input[i] == '=' ? 0 & i++ : kDecodingTable[static_cast<int>(input[i++])];
        uint32_t c = input[i] == '=' ? 0 & i++ : kDecodingTable[static_cast<int>(input[i++])];
        uint32_t d = input[i] == '=' ? 0 & i++ : kDecodingTable[static_cast<int>(input[i++])];

        uint32_t triple = (a << 3 * 6) + (b << 2 * 6) + (c << 1 * 6) + (d << 0 * 6);

        if (j < out_len) out[j++] = (triple >> 2 * 8) & 0xFF;
        if (j < out_len) out[j++] = (triple >> 1 * 8) & 0xFF;
        if (j < out_len) out[j++] = (triple >> 0 * 8) & 0xFF;
    }

    return "";
}

std::string cleaner(std::string input) {
    
    
    std::string decoded;
    
    std::string encoded = input;
    encoded.erase(std::remove(encoded.begin(), encoded.end(), '\n'), encoded.end());  
    
    if(!encoded.empty())
        Decode(encoded, decoded);
   
    return decoded;
}

std::string SerializeData( Response* response) {

    cJSON* Json = cJSON_CreateObject();

    cJSON_AddItemToObject(Json, "UUID", cJSON_CreateString(response->UUID.c_str()));
    cJSON_AddItemToObject(Json, "TASKID", cJSON_CreateString(response->taskid.c_str()));
    cJSON_AddItemToObject(Json, "DATA", cJSON_CreateString(response->Payload.c_str()));


    std::string printable = cJSON_Print(Json);
    return printable;

}

Task DeserializeData(std::string data) {

    Task task;

    cJSON* taskData = cJSON_Parse(data.c_str());

    task.taskid = cJSON_GetObjectItem(taskData, "TASKID")->valuestring;
    task.task = cJSON_GetObjectItem(taskData, "TASK")->valuestring;
    task.method = cJSON_GetObjectItem(taskData, "METHOD")->valuestring;
    task.target = cJSON_GetObjectItem(taskData, "TARGET")->valuestring;
    task.data = cJSON_GetObjectItem(taskData, "DATA")->valuestring;
    
    return task;

}

VOID regSerDeSer(State* state, std::string& data, BOOL in) {

    cJSON* Json;


    if (in) {
        
        Json = cJSON_Parse(data.c_str());
        
        state->uuid = cJSON_GetObjectItem(Json, "Sessionhash")->valuestring;
        state->cryptkey = cJSON_GetObjectItem(Json, "Cryptkey")->valuestring;

    
    } else {
        // serialize 
        Json = cJSON_CreateObject();

        cJSON_AddItemToObject(Json, "Implanttype", cJSON_CreateString(Imp_Type));
        cJSON_AddItemToObject(Json, "Compname", cJSON_CreateString(state->cname.c_str()));
        cJSON_AddItemToObject(Json, "Username", cJSON_CreateString(state->uname.c_str()));
        cJSON_AddItemToObject(Json, "Domain", cJSON_CreateString(state->udomain.c_str()));
        std::string response = cJSON_Print(Json);

        data = response;

    }

}


BOOL GetInfo(std::string& domain, std::string& strUser, std::string& strdomain){
    
    DWORD bufSz = MAX_COMPUTERNAME_LENGTH;
    char buf[MAX_COMPUTERNAME_LENGTH];
    GetComputerNameExA(ComputerNameDnsFullyQualified, buf, &bufSz);
    domain = buf;

    HANDLE hProcess = OpenProcess(PROCESS_QUERY_INFORMATION, FALSE, GetCurrentProcessId());
    if (hProcess == NULL)
        return FALSE;

    HANDLE hToken = NULL;
    if (!OpenProcessToken(hProcess, TOKEN_QUERY, &hToken))
    {
        CloseHandle(hProcess);
        return FALSE;
    }

    BOOL bres = GetLogonFromToken(hToken, strUser, strdomain);

    CloseHandle(hToken);
    CloseHandle(hProcess);
    
    return TRUE;
    

    
}

BOOL GetLogonFromToken(HANDLE hToken, std::string &strUser, std::string &strdomain){

    DWORD dwSize = MAX_NAME;
    BOOL bSuccess = FALSE;
    DWORD dwLength = 0;
    strUser = "";
    strdomain = "";
    PTOKEN_USER ptu = NULL;
        
    if (NULL == hToken)
        goto Cleanup;

    if (!GetTokenInformation(hToken, TokenUser, (LPVOID)ptu, NULL, &dwLength)){
        
        if (GetLastError() != ERROR_INSUFFICIENT_BUFFER)
            goto Cleanup;

        ptu = (PTOKEN_USER)HeapAlloc(GetProcessHeap(),
            HEAP_ZERO_MEMORY, dwLength);

        if (ptu == NULL)
            goto Cleanup;

    }

    if (!GetTokenInformation( hToken, TokenUser, (LPVOID)ptu, dwLength, &dwLength))
        goto Cleanup;
    
    SID_NAME_USE SidType;
    char lpName[MAX_NAME];
    char lpDomain[MAX_NAME];

    if (!LookupAccountSidA(NULL, ptu->User.Sid, lpName, &dwSize, lpDomain, &dwSize, &SidType))
    {
        DWORD dwResult = GetLastError();
        if (dwResult == ERROR_NONE_MAPPED)
            strcpy_s(lpName, "NONE_MAPPED");
        
        else
            err("LookupAccountSid Error");
     
    }
    else
    {
        strUser = lpName;
        strdomain = lpDomain;
        bSuccess = TRUE;
    }

    Cleanup:

        if (ptu != NULL)
            HeapFree(GetProcessHeap(), 0, (LPVOID)ptu);
        return bSuccess;
}
