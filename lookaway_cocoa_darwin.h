#pragma once

void LookawayRunApp(void);
void LookawayShowOverlay(const char *secs);
void LookawayUpdateOverlay(const char *secs);
void LookawayHideOverlay(void);
void LookawaySetMenu(const char *status, const char *next, const char *pause);
void LookawayDispatchTick(void);
void LookawayQuit(void);
