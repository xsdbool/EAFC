package api

// UtasErrorCode is a map that associates error codes with their corresponding error messages.
var UtasErrorCode = map[int]string{
	440:   "REGION_MISMATCH",
	442:   "ACCOUNT_MISMATCH",
	457:   "EVENT_EXPIRED",
	458:   "CAPTCHA_REQUIRED",
	460:   "UT_BAD_REQUEST",
	461:   "PERMISSION_DENIED",
	462:   "STATE_INVALID",
	463:   "NO_BID_TOKENS",
	464:   "NO_TITLE_ID",
	465:   "NO_USER",
	466:   "NAME_EXISTS",
	467:   "PROFANITY",
	468:   "LOGGED_IN_ON_CONSOLE_LEGACY",
	469:   "DELETING_LAST_SQUAD",
	470:   "NOT_ENOUGH_CREDIT",
	471:   "ITEM_EXISTS",
	472:   "DUPLICATE_ITEM_TYPE",
	473:   "DESTINATION_FULL",
	474:   "LOGGED_IN_ON_CONSOLE",
	475:   "NO_CARD_EXISTS",
	476:   "CARD_IN_TRADE",
	477:   "INVALID_DECK",
	478:   "NO_TRADE_EXISTS",
	479:   "INVALID_OWNER",
	480:   "SERVICE_IS_DISABLED",
	486:   "PLAYER_HAS_RED_CARD",
	487:   "REMOVE_WATCH_FAILURE",
	488:   "SWAP_ITEM_WITH_ITSELF",
	489:   "DID_CREATE_EXCEEDED",
	490:   "DID_LOGIN_EXCEEDED",
	491:   "DEVICE_SUSPENDED",
	492:   "SBC_EXPIRED",
	494:   "LOCKED_TRANSFER_MARKET",
	495:   "SOME_ITEMS_NOT_FREE",
	512:   "SOFTBAN_FUNCTION",
	521:   "SOFTBAN_FUNCTION",
	20000: "ACCOUNT_BANNED",
	20001: "UPDATE_REQUIRED",
	20002: "INVALID_CREDENTIALS",
	20003: "GEOIP_DENIED",
	20004: "UNRECOVERABLE",
}
