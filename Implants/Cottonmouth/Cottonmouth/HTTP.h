#ifndef _HTTP_INCLUDED
#define _HTTP_INCLUDED

#include "Core.h"
#include "Tasking.h"



int prep(State* state);

std::string Get(State state);
VOID Post(State *state, std::string &outbound, BOOL reg);

BOOL regimp(State* state);
BOOL Callback(State state, std::string &inbound, std::string &outbound);


#endif