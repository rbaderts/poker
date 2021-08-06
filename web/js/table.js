'use strict';

import Vue from "
import axios from "../libs/axios.min.js";
/// import { connect, StringCodec } from '../libs/nats.js'

import main from './main';



var tableInstance = null;
var OurSeatNumber = 0;
var NatsClient = null;

//import { connect, StringCodec } from '../nats.js'


function TableRoot(win, jwt, ourSeatNumber, seatsub) {
	OurSeatNumber = ourSeatNumber;
	var emptyTable = {
		id: 0,
		name: "",
		numSeats: 0,
		seats: [],
		commonCards: [],
		commonCardCount: 0,
		players: [],
		player: null,
		beaconId: 0,
		messages: []
	};


	var emptyGameResult = {
		winner: -1,
		winningHand: ""
	};

	tableInstance = new Vue({

		el: '#table_container',
		props: {

			potValue: {
				type: Number,
				default: 0
			},
			bettingRound: {
				type: Number,
				default: 0
			}
		},
		data: {
			//common_cards: [],
			//players: [],
			message: '',
			button: 0,
			callAmount: 0,
			currentBet: 0,
			requestId: 0,
			showBetPanel: false,
			seatsubject: seatsub,
			gameStage: 0,
			jwt: jwt,
			isFetching: true,
			players: [],
			player: null,
			ourCards: [],
			beaconId: 0,
			commonCardCount: 0,
			commonCards: [{index: 0, suit: 0},
				{index: 0, suit: 0},
				{index: 0, suit: 0}, {index: 0, suit: 0},
				{index: 0, suit: 0}],
			figures: [
				'1',   // spades
				'2',   // hearts
				'3',   // diamonds
				'4'    // clubs
			],
			values: [
				'A',
				'2',
				'3',
				'4',
				'5',
				'6',
				'7',
				'8',
				'9',
				'10',
				'J',
				'Q',
				'K'

			],
			messages: []
		},
		created: function () {
			try {
				this.fetchTableStatus(window.tableId);
//				if (tableData) {
//					this.players = tableData.players;
//				}
//				this.isFetching = false;
				Attach(this.jwt, window.tableId, OurSeatNumber, this.seatsubject)
			} catch (err) {
				console.log(err);
//				alert(err);
			}
		},
		destroyed: function () {
			console.log("table destroyed");
			var vm = this;
			axios.post(`/api/tables/` + window.tableId + `/unjoin/` + OurSeatNumber, "test", {withCredentials: true});
			//axios.get('/tables/' + vm.id + "/unjoin", {withCredentials: true} )
		},
		computed: {
			// a computed getter
			cards: function () {
				let all = []
				for (let figure of this.figures) {
					for (let value of this.values) {
						all.push({
							suit: figure,
							index: value
						})
					}
				}
				return all
			}
		},
		/*
		components: {
			'seat': seat
		},
		 */

		methods: {
			initialize: function () {
				this.callAmount = 0;
				this.player = null;
				this.currentBet = 0;
				this.requestId = 0;
				this.showBetPanel = false;
				this.gameStage = 0;
				this.potValue = 0;
				this.bettingRound = 0;
				this.message = 'Good luck';
				this.currentBet = 0;
				this.beaconId = 0;
				this.commonCardCount = 0;
				this.commonCards = [{index: 0, suit: 0},
					{index: 0, suit: 0},
					{index: 0, suit: 0},
					{index: 0, suit: 0},
					{index: 0, suit: 0}];

			},
			initializeGame: function (tabledata) {
				let vm = this;
				let i = 0;
				console.log("tabledata = " + JSON.stringify(tabledata));
				for (i = 0; i < tabledata.players.length; i++) {
					let p = tabledata.players[i];
					if (i == tabledata.button) {
						p.dealer = true;
					} else {
						p.dealer = false;
					}
					p["ontable"] = 0;
					p["message"] = "";
					p["folded"] = false;
					if (i == OurSeatNumber) {
						p.card1 = null;
						p.card2 = null;
						p.hand = "";
						vm.player = p;
						Vue.set(vm.players, i, p);
						vm.seatsubject = p["channel"]

					} else {
						Vue.set(vm.players, i, p);
					}
				}
				console.log("Cleared Hands");

			},

			updateOdds: function (odds) {
				let i = 0;
				let vm = this;
				for (const seatnum in odds) {
					let i = parseInt(seatnum);
					let wins = odds[i].wins * 100;
					let ties = odds[i].ties * 100;
					if (i == OurSeatNumber) {
						vm.player.percentWin = Math.round(wins * 1e2) / 1e2;
						vm.player.percentTie = Math.round(ties * 1e2) / 1e2;
					} else {
						vm.players[i].percentWin = Math.round(wins * 1e2) / 1e2;
						vm.players[i].percentTie = Math.round(ties * 1e2) / 1e2;
						//vm.players[i].percentWin = wins.toFixed(1);
						//vm.players[i].percentTie = ties.toFixed(1);
					}
				}
				/*
					console.log(`${property}: ${object[property]}`);

					for (i = 0; i < odds.length; i++) {
					if (i == OurSeatNumber) {
						vm.player.percentWin = odds[i].win;
						vm.player.percentTie = odds[i].tie;
					} else {
						vm.players[i].percentWin = odds[i].win;
						vm.players[i].percentTie = odds[i].tie;
					}
				}
				 */

			},

			fetchTableStatus: function (id) {
				var vm = this;
				axios.get('/api/tables/' + window.tableId, {withCredentials: true})
					.then(function (response) {
						const tb = response.data;
						let i = 0;
						console.log("tb = " + JSON.stringify(tb));
						for (i = 0; i < tb.players.length; i++) {
							let p = tb.players[i];
							if (p == null) {
								continue;
							}
							if (i == tb.button) {
								p.dealer = true;
							} else {
								p.dealer = false;
							}
							p["ontable"] = 0;
							p["message"] = "";
							p["folded"] = false;
//							console.log("seatnum = " + p.seatnum.toString() + ", OurSeatNumber = " + OurSeatNumber.toString());
							if (i == OurSeatNumber) {
								if (vm.players[i] == null) {
									vm.player = p;
								} else {
									vm.player = vm.players[i];
									vm.player["ontable"] = p["ontable"];
									vm.player["message"] = p["message"];
									vm.player["folded"] = p["folded"];
									vm.player.dealer = p["dealer"];
								}
								vm.player["stack"] = p["stack"];
								vm.player["percentwin"] = p["percentwin"];
								vm.player["percenttie"] = p["percenttie"];
								vm.player["folded"] = p["folded"];
								vm.player.hand = "";
								Vue.set(vm.players, i, vm.player);
//								vm.seatsubject = p["channel"];
//								p["card1"] = null;
//								p["card2"] = null;
							} else {
								Vue.set(vm.players, i, p);
							}
						}
						console.log("Cleared Hands");
						vm.isFetching = false;
						return tb
					})
					.catch(function (error) {
						// handle error
						alert(error)
					})
					.then(function () {
					});
			},

			hasPlayer: function (id) {
				var i;
				for (i = 0; i < this.players.length; i++) {
					if (this.players[i].id == id) {
						return i
					}
				}
				console.log("player: " + id + " not found");
				return -1
			},
			updatePlayerHand: function (desc) {
				this.player.hand = desc;
				this.players[OurSeatNumber].hand = desc;
			},
			update: function (action) {
				this.players[action.seat].ontable = action.ontable;
				this.players[action.seat].stack = action.stack;
				this.players[action.seat].hasCards = action.hasCards;

			},
			betUpdate: function (seat, amount) {
				this.players[seat].ontable += amount;
				this.players[seat].stack -= amount

			},
			holeCardsUpdate: function (seat, amount) {
				this.players[seat].ontable += amount;
				this.players[seat].stack -= amount
			},

			AddMessage: function (msg) {
				this.messages.push(msg)
			},

			turn: function (card) {
				this.commonCards.push(card);
			},
			removePlayer: function (seat) {
				let vm = this;
				console.log("removePlayer");
				if (seat == OurSeatNumber) {
					vm.player = null
				}
				Vue.delete(vm.players, seat)
			},
			updatePlayer: function (seat) {

				let vm = this;
				console.log("updatePlayer");
				axios.get('/api/tables/' + window.tableId, {withCredentials: true})
					.then(function (response) {
						let tb = response.data;
						console.log("update, tb = " + JSON.stringify(tb));
						Vue.set(vm.players, seat, tb.players[seat]);
					})
					.catch(function (error) {
						//alert(error)
						console.log(error);
					})
					.then(function () {
					});
			},
			AddToBet: function (amount) {
				console.log("AddToBet: " + amount);
				//this.player["stack"] = this.player["stack"]-amount;

				//this.addontable += amount
			},
			handleBet: function (callAmount, currentBet, requestId) {

				console.log("handleBet: requestId = " + requestId);
				this.callAmount = callAmount;
				this.currentBet = currentBet;
				this.requestId = requestId;
				this.stackAvailable = this.player["stack"];
				this.showBetPanel = true;

			},

			newgame: function () {
				axios.post(`/api/tables/` + window.tableId + `/startgame`, "test", {withCredentials: true});
			},
			leavetable: function () {
				axios.post(`/api/tables/` + window.tableId + `/unjoin/` + OurSeatNumber, "test", {withCredentials: true});
			}
		}
	});
	return tableInstance;
    window.CurrentCallAmount = 0;
}

//const init = async function () {

const Attach = async function (jwt, tableId, ourseatnumber, seatSub) {
//function Attach(jwt, tableId, ourseatnumber, seatSub) {

	var subject = "t_"+tableId.toString();
	console.log("subscribing to channel:" + subject);

	console.log("Attache: jwt: " + jwt + ", tableId = " + tableId);

	/*
	NatsClient = await natsws.connect(
		{
			servers: ["ws://localhost:9222", "wss://localhost:2229", "localhost:9111"],
		},
	);
	 */

	window.centrifuge = new Centrifuge('ws://localhost:8000/connection/websocket');

	window.centrifuge.setToken(jwt);
	window.centrifuge.connect();




	var tableChannelCallbacks = {
		"publish": function(message) {
				console.log("received msg on " + subject);
				//let msg = JSON.parse(message);
				console.log("msg = " + JSON.stringify(message));
				handleTableServerCommand(message);
		},
		"join": function(message) {
			// See below description of join message format
			console.log("tableChannel join: " + message);
		},
		"leave": function(message) {
			// See below description of leave message format
			console.log("tableChannel leave: " + message);
		},
		"subscribe": function(context) {
			// See below description of subscribe callback context format
			console.log("tableChannel subscribe: " + context);
		},
		"error": function(errContext) {
			// See below description of subscribe error callback context format
			console.log("tableChannel error: " + err);
		},
		"unsubscribe": function(context) {
			// See below description of unsubscribe event callback context format
			console.log("tableChannel unsubscribe: " + context);
		}
	};
	
	var seatChannelCallbacks = {
		"publish": function(message) {
			console.log("received msg on seat channel: " + sub);
			//let msg = JSON.parse(message);
			console.log("msg = " + JSON.stringify(message));
			handleSeatCommand(message);
		},
		"join": function(message) {
			// See below description of join message format
			console.log("seatChannel join: " + message);
		},
		"leave": function(message) {
			// See below description of leave message format
			console.log("seatChannel leave: " + message);
		},
		"subscribe": function(context) {
			// See below description of subscribe callback context format
			console.log("seatChannel subscribe: " + context);
		},
		"error": function(errContext) {
			// See below description of subscribe error callback context format
			console.log("seatChannel error: " + err);
		},
		"unsubscribe": function(context) {
			// See below description of unsubscribe event callback context format
			console.log("seatChannel unsubscribe: " + context);
		}
	};

	const sub = await nc.subscribe(subject);
	(async () => {
		for await (const m of subject) {
			console.log("m = " + m);
			m.respond(sc.encode(`I can help ${sc.decode(m.data)}`));
		}
	})().then();

//	NatsClient.subscribe(sub, )
//	window.centrifuge.subscribe(sub, tableChannelCallbacks);

	console.log("subscribing to channel:" + seatSub);
//	window.centrifuge.subscribe(seatSub, seatChannelCallbacks);

	/*
	window.centrifuge.on('connect', function(context) {
		console.log("centrifugo connectedd")
		// now client connected to Centrifugo and authorized
	});

	window.centrifuge.on('disconnect', function(context) {
		// do whatever you need in case of disconnect from server
		console.log("centrifugo disconnected")
	});
	*/

//	window.addEventListener("unload", LeaveTable);
	/*
	chrome.tabs.onAttached.addListener(function(tabId, props) {
		consolje.log("tab attached");
	});
	 */

	async function dispatch(s) {
		let subj = s.getSubject();
		console.log(`listening for ${subj}`);
		const c = (13 - subj.length);
		const pad = "".padEnd(c);
		for await (const m of s) {
			console.log(
				`[${subj}]${pad} #${s.getProcessed()} - ${m.subject} ${
					m.data ? " " + sc.decode(m.data) : ""
				}`,
			);
		}
	}

}

function handleSeatCommand(cmd) {
	let msg = cmd;
	let data = msg["data"];
	//let msg = JSON.parse(cmd);
	let typ = data["typ"];
	let requestId = data["id"];
	console.log("handleSeatCommand: " + GetCommandName(typ));

	if (typ == GameSeatUpdated) {
		let hand = data["data"]["hand"];
		tableInstance.updatePlayerHand(hand)
	} else if (typ == BetRequest)  {
		let callAmount = data["data"]["callAmount"];
		let currentBet = data["data"]["currentBet"];
	    tableInstance.handleBet(callAmount, currentBet, requestId);

///		tableInstance.$refs.player.handleBet(callAmount, currentBet, requestId)
//		tableInstance.handleBet(callAmount, currentBet, requestId);
	} else if (typ == GameHoleCardsDraw) {
		console.log("Received GameHoleCardsDraw");
		let card1 = data["data"]["card1"];
		let card2 = data["data"]["card2"];
		console.log("Received card1: " + JSON.stringify(card1));
		console.log("Received card2: " + JSON.stringify(card2));
		if ((tableInstance.players[OurSeatNumber]) && (card1 != null) ) {
			tableInstance.players[OurSeatNumber].card1 =
				{suit: card1["suit"], index: card1["index"]};
			tableInstance.players[OurSeatNumber].card2 =
				{suit: card2["suit"], index: card2["index"]};
		} else {
			//tableInstance.players[OurSeatNumber].card1 =
		//		{suit: 0, index: 0};
		//	tableInstance.players[OurSeatNumber].card2 =
		//		{suit: 0, indx: 0};

		}
		createjs.Sound.play(CardSlide1);
	}

}



function handleTableServerCommand(cmd) {
	//let msg = JSON.parse(cmd);

	let msg = cmd;
	let data = msg["data"];
	//let typ = data["actionType"];
	let typ = data["typ"];
	let RequestId = data["id"];

	console.log("command: " + GetCommandName(typ));

	let cmdData = data["data"];
	if (typ == TableMessage) {
		let msg = cmdData["message"];
		tableInstance.AddMessage(msg);


	} else if (typ == TableGameStarted) {
		let tableData = cmdData;
		tableInstance.initialize();
		tableInstance.initializeGame(tableData);
		tableInstance.bettingRound = 1;
		//tableInstance.fetchTableStatus(window.tableId);
//		let seat = cmdData["seat"];
	} else if (typ == TableOddsUpdated) {

		console.log("cmdData = " + JSON.stringify(cmdData));
		let oddsData = cmdData["odds"];
		console.log("OddsData = " + JSON.stringify(oddsData));
		tableInstance.updateOdds(oddsData)

	} else if (typ == TableStatusUpdated) {
		tableInstance.fetchTableStatus();
	} else if (typ == TablePlayersChange) {
		let seat = cmdData["seat"];
		tableInstance.updatePlayer(seat);
	} else if (typ == TablePlayerMessage) {
		let seat = cmdData["seat"];
		let message = cmdData["message"];
		if (seat == OurSeatNumber) {
			tableInstance.$refs.player.showPlayerMessage(message)
		}
	} else if (typ == GameShowdown) {
//		tableInstance.gameResult.winner = cmdData["winningSeat"];
//		tableInstance.gameResult.winningHand = cmdData["winningHand"];
		//let msg = "Seat " + cmdData["winningSeat"] + " wins with " + cmdData['winningHand'];
		let msg = cmdData["message"];
		let payouts = cmdData["payouts"];
//		let hands = cmdData["hands"];
//		tableInstance.message = msg;

		tableInstance.AddMessage(msg);


//		for (const )
		/*
		for (const [seatnumstr, cards] of Object.entries(hands)) {
			let seatnum = parseInt(seatnumstr);
			tableInstance.players[seatnum].card1 = cards[0];
			tableInstance.players[seatnum].card2 = cards[1];
		}
		 */

		//console.log("Show All hands = " + JSON.stringify(hands));

		//tableInstance.handHistory.push(msg);

		for (let i = 0; i < payouts.length; i++) {
			let payout = payouts[i];
			tableInstance.players[payout.seat].card1 = payout.card1
			tableInstance.players[payout.seat].card2 = payout.card2
			tableInstance.players[payout.seat].stack = payout.amount;
		}



	} else if (typ == GameStageChanged) {
		tableInstance.gameStage = cmdData["stage"]
	} else if (typ == GameCardsDealt) {

		let cards = cmdData;
		let i = 0;
		for (i = 0; i < cards.length; i++) {
         	Vue.set(tableInstance.commonCards, tableInstance.commonCardCount, cards[i]);
			tableInstance.commonCardCount += 1
		}

		createjs.Sound.play(CardSlide1);

	} else if (typ == TablePlayerJoined) {
		let seat = cmdData["seatNum"];
		console.log("TablePlayerJoin: seat = " + seat);
		tableInstance.updatePlayer(seat);

	} else if (typ == TablePlayerLeft) {
		let seat = cmdData["seatNum"];
		console.log("TablePlayerLeft: seat = " + seat);
		tableInstance.removePlayer(seat);

	} else if (typ == GamePlayerActionFold) {
		let seat = cmdData["seatNum"];
		console.log("seat: " + seat + " folded");
		tableInstance.message = "Player Folded";
		tableInstance.players[seat].folded = true;
		createjs.Sound.play(ChipSound1);

	} else if (typ == GamePlayerActionBet) {
		let betType = cmdData["betType"];
		let amount = cmdData["amount"];
		let seat = cmdData["seatNum"];

		tableInstance.potValue += amount;

		console.log("PlayerBet: seat = " + seat + ", type = " + betType);

		if (betType == SmallBlind) {
			console.log("   Small Blind amount = " + amount);
			tableInstance.players[seat]["ontable"] += amount;
			tableInstance.players[seat].stack -= amount;

		} else if (betType == BigBlind) {
			let amount = cmdData["amount"];
			console.log("   Big Blind amount = " + amount);
			tableInstance.players[seat]["ontable"] += amount;
			tableInstance.players[seat].stack -= amount;

		} else if (betType == Call) {
			let amount = cmdData["amount"];
			console.log("   Call amount = " + amount);
			tableInstance.players[seat]["ontable"] += amount;
			tableInstance.players[seat].stack -= amount;
			tableInstance.message = "Player Called";
			createjs.Sound.play(ChipSound1);

		} else if (betType == Raise) {
			let amount = cmdData["amount"];
			console.log("   Raise  amount = " + amount);
			tableInstance.players[seat]["ontable"] += amount;
			tableInstance.players[seat].stack -= amount;
			tableInstance.message = "Player Raised";
			createjs.Sound.play(ChipsStack6);

		} else if (betType == Check) {
			console.log("   Check");
			tableInstance.message = "Player Checked";
			createjs.Sound.play(ChipSound1);

		} else if (betType == Fold) {
			console.log("   Fold");
			tableInstance.message = "Player Folded";
			tableInstance.players[seat].folded = true;
			createjs.Sound.play(ChipSound1);


    	} else if (betType == AllIn) {
			let amount = cmdData["amount"];
			console.log("   AllIn");
			console.log("   player Stack current: " + tableInstance.players[seat].stack + "\n");
			console.log("   bet amount : " + amount + "\n");
			tableInstance.players[seat]["ontable"] += amount;
			tableInstance.players[seat].stack -= amount;
			tableInstance.message = "Player All In";
			createjs.Sound.play(AllInSound);

		}
	} else if (typ == PlayerPaidOut) {
		let seat = cmdData["seat"];
		let amount = cmdData["amount"];
		tableInstance.players[seat].stack += amount
    } else if (typ == BettingRoundComplete) {
		tableInstance.bettingRound = tableInstance.bettingRound+1;
		let amount = cmdData["total"];
		for (let i = 0; i < tableInstance.players.length; i++) {
			tableInstance.players[i].ontable = 0;
		}
		//tableInstance.potAmount += amount;
		createjs.Sound.play(ChipsStack6);
	}

}

$('#bet-submit-button').on('click', function (e) {

	let seatsub = tableInstance.seatsub;
	let betform = document.getElementById('form');
	let amountElem = document.getElementById('bet-amount');
	let amount = amountElem.value;
	console.log('amount = ' + amount);

	let reqId = betform.getAttribute("requestId");
	console.log("respond to be request with Id : + " + reqId);

	let action = {
		Id: ActionCounter,
		ResponseTo: Number(reqId),
		data: {amount:Number(amount)}
	};

	let id = window.tableId;
	axios.post(`/api/tables/`+id+`/bet`, action, {withCredentials: true} );

//	WS.send(JSON.stringify(action));

	betpanel.style.display = 'none';
	/*
		Id            int64
	ResponseTo    int64
	SourceSessionId  string
	ActionType    ActionType          `json:"actionType"`
	Data          interface{}         `json:"data"
	 */
});

function LeaveTable() {
	vm.destroy()
}

export default TableRoot;

