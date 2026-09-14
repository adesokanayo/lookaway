#pragma once

#include <stdbool.h>

void LookawayRunApp(void);
bool LookawayScreenIsLocked(void);
bool LookawayInputIsQuiet(void);
void LookawaySetEye(bool warning, bool flash);
void LookawayShowOverlay(const char *secs);
void LookawayUpdateOverlay(const char *secs);
void LookawayHideOverlay(void);
void LookawaySetMenu(const char *status, const char *next, const char *stats, const char *pause);
void LookawayShowWelcome(void);
void LookawayHideWelcome(void);
void LookawayDispatchTick(void);
void LookawayQuit(void);
