#import <Cocoa/Cocoa.h>
#import <CoreGraphics/CoreGraphics.h>
#include "lookaway_cocoa_darwin.h"

void lookawayOnReady(void);
void lookawayOnTick(void);
void lookawayOnBreakNow(void);
void lookawayOnPause(void);
void lookawayOnSkip(void);
void lookawayOnQuit(void);
void lookawayOnUnlocked(void);
void lookawayOnWelcomeDone(void);

static NSStatusItem *gStatus;
static NSMenuItem *gNextItem;
static NSMenuItem *gStatsItem;
static NSMenuItem *gPauseItem;
static NSWindow *gWelcome;
static NSMutableArray<NSWindow *> *gWindows;
static NSMutableArray<NSTextField *> *gCounts;

@interface LookawayTarget : NSObject
@end
@implementation LookawayTarget
- (void)breakNow:(id)sender { lookawayOnBreakNow(); }
- (void)pause:(id)sender { lookawayOnPause(); }
- (void)skip:(id)sender { lookawayOnSkip(); }
- (void)quit:(id)sender { lookawayOnQuit(); }
- (void)gotIt:(id)sender { lookawayOnWelcomeDone(); }
- (void)screenUnlocked:(NSNotification *)note {
	(void)note;
	lookawayOnUnlocked();
}
- (void)didWake:(NSNotification *)note {
	(void)note;
	lookawayOnUnlocked();
}
@end

bool LookawayScreenIsLocked(void) {
	CFDictionaryRef dict = CGSessionCopyCurrentDictionary();
	if (!dict) {
		return false;
	}
	CFBooleanRef locked = CFDictionaryGetValue(dict, CFSTR("CGSSessionScreenIsLocked"));
	bool isLocked = locked != NULL && CFBooleanGetValue(locked);
	CFRelease(dict);
	return isLocked;
}

static LookawayTarget *gTarget;

static BOOL LookawayClaimSoleInstance(void) {
	NSString *bid = [[NSBundle mainBundle] bundleIdentifier];
	if (![bid isEqualToString:@"com.adesokanayo.lookaway"]) {
		return YES;
	}
	pid_t me = [[NSProcessInfo processInfo] processIdentifier];
	for (NSRunningApplication *app in [NSRunningApplication runningApplicationsWithBundleIdentifier:bid]) {
		if (app.processIdentifier != me) {
			return NO;
		}
	}
	return YES;
}

static NSTextField *label(NSString *text, NSColor *color, CGFloat size, BOOL bold, NSRect frame) {
	NSTextField *f = [NSTextField labelWithString:text];
	f.font = bold ? [NSFont boldSystemFontOfSize:size] : [NSFont systemFontOfSize:size];
	f.textColor = color;
	f.alignment = NSTextAlignmentCenter;
	f.drawsBackground = NO;
	f.frame = frame;
	f.maximumNumberOfLines = 0;
	f.lineBreakMode = NSLineBreakByWordWrapping;
	return f;
}

void LookawayHideOverlay(void) {
	for (NSWindow *w in gWindows) {
		[w orderOut:nil];
		[w close];
	}
	[gWindows removeAllObjects];
	[gCounts removeAllObjects];
}

void LookawayShowOverlay(const char *secs) {
	LookawayHideOverlay();
	if (!gWindows) {
		gWindows = [NSMutableArray new];
	}
	if (!gCounts) {
		gCounts = [NSMutableArray new];
	}

	NSColor *bg = [NSColor colorWithSRGBRed:0.04 green:0.05 blue:0.07 alpha:1];
	NSColor *fg = [NSColor colorWithSRGBRed:0.96 green:0.97 blue:0.98 alpha:1];
	NSColor *muted = [NSColor colorWithSRGBRed:0.70 green:0.74 blue:0.78 alpha:1];
	NSString *count = [NSString stringWithUTF8String:secs];
	NSScreen *main = [NSScreen mainScreen];

	for (NSScreen *screen in [NSScreen screens]) {
		NSRect frame = screen.frame;
		NSPanel *panel = [[NSPanel alloc] initWithContentRect:frame
			styleMask:(NSWindowStyleMaskBorderless | NSWindowStyleMaskNonactivatingPanel)
			backing:NSBackingStoreBuffered
			defer:NO];
		panel.level = NSScreenSaverWindowLevel;
		panel.collectionBehavior = NSWindowCollectionBehaviorCanJoinAllSpaces |
			NSWindowCollectionBehaviorFullScreenAuxiliary |
			NSWindowCollectionBehaviorStationary;
		panel.opaque = YES;
		panel.hasShadow = NO;
		panel.backgroundColor = bg;
		panel.ignoresMouseEvents = NO;
		panel.hidesOnDeactivate = NO;
		panel.releasedWhenClosed = NO;

		CGFloat w = frame.size.width;
		CGFloat h = frame.size.height;
		NSView *c = panel.contentView;
		[c addSubview:label(@"Look away", fg, 42, YES, NSMakeRect((w - 720) / 2, h * 0.58, 720, 56))];
		[c addSubview:label(@"Look about 20 feet away.", muted, 20, NO, NSMakeRect((w - 720) / 2, h * 0.48, 720, 40))];
		NSTextField *n = label(count, fg, 96, YES, NSMakeRect((w - 240) / 2, h * 0.30, 240, 110));
		n.font = [NSFont monospacedDigitSystemFontOfSize:96 weight:NSFontWeightSemibold];
		[c addSubview:n];
		[gCounts addObject:n];

		if (NSEqualRects(screen.frame, main.frame)) {
			NSButton *skip = [NSButton buttonWithTitle:@"Skip" target:gTarget action:@selector(skip:)];
			skip.bezelStyle = NSBezelStyleInline;
			skip.frame = NSMakeRect((w - 120) / 2, h * 0.16, 120, 32);
			[c addSubview:skip];
		}

		[panel orderFrontRegardless];
		[gWindows addObject:panel];
	}
	[[NSSound soundNamed:@"Tink"] play];
}

void LookawayUpdateOverlay(const char *secs) {
	NSString *t = [NSString stringWithUTF8String:secs];
	for (NSTextField *f in gCounts) {
		f.stringValue = t;
	}
}

void LookawaySetMenu(const char *status, const char *next, const char *stats, const char *pause) {
	gStatus.button.title = [NSString stringWithUTF8String:status];
	gNextItem.title = [NSString stringWithUTF8String:next];
	gStatsItem.title = [NSString stringWithUTF8String:stats];
	gPauseItem.title = [NSString stringWithUTF8String:pause];
}

void LookawayHideWelcome(void) {
	[gWelcome orderOut:nil];
	[gWelcome close];
	gWelcome = nil;
}

void LookawayShowWelcome(void) {
	LookawayHideWelcome();
	NSColor *bg = [NSColor colorWithSRGBRed:0.04 green:0.05 blue:0.07 alpha:1];
	NSColor *fg = [NSColor colorWithSRGBRed:0.96 green:0.97 blue:0.98 alpha:1];
	NSColor *muted = [NSColor colorWithSRGBRed:0.70 green:0.74 blue:0.78 alpha:1];
	NSRect frame = NSMakeRect(0, 0, 460, 260);
	NSPanel *panel = [[NSPanel alloc] initWithContentRect:frame
		styleMask:(NSWindowStyleMaskTitled | NSWindowStyleMaskNonactivatingPanel)
		backing:NSBackingStoreBuffered
		defer:NO];
	panel.title = @"Lookaway";
	panel.level = NSFloatingWindowLevel;
	panel.opaque = YES;
	panel.backgroundColor = bg;
	panel.hidesOnDeactivate = NO;
	panel.releasedWhenClosed = NO;
	NSView *c = panel.contentView;
	[c addSubview:label(@"Lookaway is in the menu bar.", fg, 18, YES, NSMakeRect(24, 190, 412, 28))];
	[c addSubview:label(@"Look for the eye. Every 20 minutes the screen takes over for 20 seconds. Look about 20 feet away. That counts as one lookaway.", muted, 13, NO, NSMakeRect(24, 88, 412, 96))];
	NSButton *ok = [NSButton buttonWithTitle:@"Got it" target:gTarget action:@selector(gotIt:)];
	ok.bezelStyle = NSBezelStyleRounded;
	ok.frame = NSMakeRect(170, 24, 120, 32);
	[c addSubview:ok];
	[panel center];
	[panel orderFrontRegardless];
	gWelcome = panel;
	dispatch_after(dispatch_time(DISPATCH_TIME_NOW, 8 * NSEC_PER_SEC), dispatch_get_main_queue(), ^{
		if (gWelcome != nil) {
			lookawayOnWelcomeDone();
		}
	});
}

void LookawayQuit(void) {
	LookawayHideWelcome();
	LookawayHideOverlay();
	[NSApp terminate:nil];
}

void LookawayDispatchTick(void) {
	dispatch_async(dispatch_get_main_queue(), ^{
		lookawayOnTick();
	});
}

void LookawayRunApp(void) {
	@autoreleasepool {
		[NSApplication sharedApplication];
		[NSApp setActivationPolicy:NSApplicationActivationPolicyAccessory];
		if (!LookawayClaimSoleInstance()) {
			return;
		}
		gTarget = [LookawayTarget new];
		gStatus = [[NSStatusBar systemStatusBar] statusItemWithLength:NSVariableStatusItemLength];
		NSImage *img = [NSImage imageWithSystemSymbolName:@"eye" accessibilityDescription:@"Lookaway"];
		img.template = YES;
		gStatus.button.image = img;
		gStatus.button.imagePosition = NSImageLeft;

		NSMenu *menu = [[NSMenu alloc] initWithTitle:@"Lookaway"];
		gNextItem = [[NSMenuItem alloc] initWithTitle:@"Next break" action:nil keyEquivalent:@""];
		gNextItem.enabled = NO;
		[menu addItem:gNextItem];
		gStatsItem = [[NSMenuItem alloc] initWithTitle:@"Today 0 lookaways" action:nil keyEquivalent:@""];
		gStatsItem.enabled = NO;
		[menu addItem:gStatsItem];
		[menu addItem:[NSMenuItem separatorItem]];
		[menu addItem:[[NSMenuItem alloc] initWithTitle:@"Break now" action:@selector(breakNow:) keyEquivalent:@""]];
		gPauseItem = [[NSMenuItem alloc] initWithTitle:@"Pause" action:@selector(pause:) keyEquivalent:@""];
		[menu addItem:gPauseItem];
		[menu addItem:[[NSMenuItem alloc] initWithTitle:@"Skip" action:@selector(skip:) keyEquivalent:@""]];
		[menu addItem:[NSMenuItem separatorItem]];
		[menu addItem:[[NSMenuItem alloc] initWithTitle:@"Quit" action:@selector(quit:) keyEquivalent:@"q"]];
		for (NSMenuItem *it in menu.itemArray) {
			it.target = gTarget;
		}
		gNextItem.target = nil;
		gStatsItem.target = nil;
		gStatus.menu = menu;

		NSDistributedNotificationCenter *dnc = [NSDistributedNotificationCenter defaultCenter];
		[dnc addObserver:gTarget
		        selector:@selector(screenUnlocked:)
		            name:@"com.apple.screenIsUnlocked"
		          object:nil
		suspensionBehavior:NSNotificationSuspensionBehaviorDeliverImmediately];
		[[[NSWorkspace sharedWorkspace] notificationCenter]
			addObserver:gTarget
			   selector:@selector(didWake:)
			       name:NSWorkspaceDidWakeNotification
			     object:nil];

		lookawayOnReady();
		[NSApp run];
	}
}
