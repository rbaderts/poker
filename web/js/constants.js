'use strict';


const NODE_WIDTH=160;
const NODE_SPACE=26;
const NODE_HEIGHT=44;
const HEIGHT = 1000;
const WIDTH = 1000;
const MIN_HEIGHT = 1000;
const MIN_WIDTH = 1000;

const RIGHT_TO_LEFT = 1;
const CENTERED = 2;

const SIZE_BASE = 6;

const COLOR_DARK_BACKGROUND = '#3C3D40';
const COLOR_NAV_BACKGROUND = '#51585E';

const COLOR_1 = '#071425';
const COLOR_2 = '#444f5e';
const COLOR_3 = '#ba411f';
const COLOR_4 = '#fbfcfe';
const COLOR_5 = '#9197a1';
const COLOR_6 = '#f1cb1c';
const COLOR_7 = '#fcef95';



/*
const COLOR_1 = '#5e315b';
const COLOR_2 = '#8c3f5d';
const COLOR_3 = '#ba6156';
const COLOR_4 = '#f2a65a';
const COLOR_5 = '#ffe478';
const COLOR_6 = '#cfff70';
const COLOR_7 = '#8fde5d';
*/
const COLOR_8 = '#3ca370';
const COLOR_9 = '#1E4F50';
const COLOR_10 = '#323e4f';
const COLOR_11 = '#322947';
const COLOR_12 = '#4b6bcb';
const COLOR_13 = '#863845';
const COLOR_14 = '#9f8670';
const COLOR_15 = '#aaaaab';
const COLOR_16 = '#bfbfcb';
const COLOR_17 = '#D5D5D7';
const COLOR_18 = '#EBEBEF';
const COLOR_19 = '#606070';
const COLOR_20 = '#43434f';
const COLOR_21 = '#272736';
const COLOR_22 = '#FFC300';
const COLOR_23 = '#57294b';
const COLOR_24 = '#964253';
const COLOR_25 = '#e36956';
const COLOR_26 = '#384863';
const COLOR_27 = '#ff9166';
const COLOR_28 = '#eb564b';
const COLOR_29 = '#b0305c';

const POOL_TABLE_GREEN_1 = '#044c36';
const POOL_TABLE_GREEN_2 = '#4c1004';
const POOL_TABLE_GREEN_3 = '#4c3404';

const BetEvent                  = 1;
const BetRequest =                    2;

const GameHoleCardsDraw =             3;
const GamePlayerAction   =            4;

const PlayerPaidOut =          5;
const GamePlayerActionBet   =        7;
const GamePlayerActionFold   =        8;
const BettingRoundComplete   =        9;
const GameCardsDealt            =      10;
const RequestStatus             =      11;
const TablePlayerMessage        =      12;
const GamePlayerActionMuckCards  =    13;
const GamePlayerActionRequestTime  =  14;
const GamePlayerActionSitOut       =  15;
const GamePlayerActionReturned     =  16;
const GamePlayerActionRequested    =  17;
const GameFlop                     = 18;
const GameTurn                     = 19;
const GameRiver                    = 20;
const GameShowdown                 = 21;
const GamePotDistributed           = 22;
const TablePlayerJoined            = 23;
const TablePlayerLeft              = 24;
const TablePlayerSittingOut        = 25;
const TablePlayerReturned          = 26;
const TablePlayerStackChange       = 27;
const TableGameStarted             = 28;
const TableStatusUpdated           = 32;
const GameCompleted                = 29;
const TablePlayersChange           = 30;
const GameStageChanged             = 31;
const TableOddsUpdated             = 33;
const GameSeatUpdated              = 34;
const TableMessage                 = 35;

function GetCommandName(cmd) {
	switch (cmd) {
		case BetEvent:
			return 'BetEvent';
		case BetRequest:
			return 'BetRequest';
		case GameHoleCardsDraw:
			return "GameHoleCardsDraw";

		case GamePlayerAction: return "GamePlayerAction";

		case PlayerPaidOut: return "PlayerPaidOut";

		case GamePlayerActionFold: return "GamePlayerActionFold";
		case GamePlayerActionBet: return "GamePlayerActionBet";
		case RequestStatus: return "RequestStatus";

		case GameCardsDealt: return "GameCardsDealt";
		case GamePlayerActionMuckCards: return "GamePlayerActionMuckCards" ;
		case GamePlayerActionRequestTime: return "GamePlayerActionRequestTime";

		case GamePlayerActionSitOut: return "GamePlayerActionSitOut";
		case GamePlayerActionReturned: return "GamePlayerActionReturned";
		case GamePlayerActionRequested: return "GamePlayerActionRequested";
		case GameFlop: return "GameFlop";
		case GameTurn: return "GameTurn";
		case GameRiver: return "GameRiver";
		case GameShowdown: return "GameShowdown";
		case GamePotDistributed: return "GamePotDistributed";
		case BettingRoundComplete: return "BettingRoundComplete";
		case TablePlayerJoined: return "TablePlayerJoined";
		case TablePlayerMessage: return "TablePlayerMessage";
		case TablePlayerLeft: return "TablePlayerLeft";
		case TablePlayerSittingOut: return "TablePlayerSittingOut";
		case TablePlayerReturned : return "TablePlayerReturned";
		case TablePlayerStackChange : return "TablePlayerStackChange";
		case TableGameStarted : return "TableGameStarted";
		case TableStatusUpdated : return "TableStatusUpdated";
		case GameCompleted: return "GameCompleted";
		case TablePlayersChange: return "TablePlayersChange";
		case GameStageChanged: return "GameStageChanged";
		case TableOddsUpdated: return "TableOddsUpdated";
		case GameSeatUpdated: return "GameSeatUpdated";
	}

	return "cmd not found";
}

const Fold       = 1;
const Call       = 2;
const Check      = 3;
const Raise      = 4;
const AllIn      = 5;
const SmallBlind = 6;
const BigBlind   = 7;

